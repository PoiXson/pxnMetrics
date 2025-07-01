package worker;
// pxnMetrics Shard - worker

import(
	Log      "log"
	Fmt      "fmt"
	Net      "net"
	Time     "time"
	Sync     "sync"
	Atomic   "sync/atomic"
	Utils    "github.com/PoiXson/pxnGoCommon/utils"
	PxnNet   "github.com/PoiXson/pxnGoCommon/utils/net"
	Service  "github.com/PoiXson/pxnGoCommon/service"
	Configs  "github.com/PoiXson/pxnMetrics/shard/configs"
	BackLink "github.com/PoiXson/pxnMetrics/shard/backlink"
	SecretDB "github.com/PoiXson/pxnMetrics/shard/worker/secretdb"
);



type Worker struct {
	mut_state    Sync.Mutex
	mut_update   Sync.Mutex
	service      *Service.Service
	config       *Configs.CfgShard
	link         *BackLink.BackLink
	listen       *Net.UDPConn
	// stats
//TODO: move to WorkerStats struct?
//TODO: add to status page
	PacketsGood  Atomic.Uint64
	PacketsBlock Atomic.Uint64
	PacketsError Atomic.Uint64
	// IP/UID request counts
	secret_db    *SecretDB.DB
	needs_batch  Atomic.Bool
}



func New(service *Service.Service, config *Configs.CfgShard,
		backlink *BackLink.BackLink) *Worker {
	return &Worker{
		service: service,
		config:  config,
		link:    backlink,
	};
}



func (worker *Worker) Start() error {
	worker.mut_state.Lock();
	defer worker.mut_state.Unlock();
	shard_index := worker.config.ShardIndex;
	num_shards  := worker.config.NumShards;
	if shard_index == 0 { return Fmt.Errorf("Invalid shard index: %d", shard_index); }
	if num_shards  == 0 { return Fmt.Errorf("Invalid total shards: %d", num_shards); }
	if shard_index > num_shards { return Fmt.Errorf(
		"Invalid shard index, out of range: %d > max %d", shard_index, num_shards); }
	worker.secret_db = SecretDB.New(worker.config);
	Log.Printf("[Shard-%d] Starting public listener.. %s",
		shard_index, worker.config.BindPublic);
	listen, err := PxnNet.NewServerUDP(worker.config.BindPublic);
	if err != nil { return err; }
	worker.listen = listen;
	go worker.Serve();
	Utils.SleepC();
	return nil;
}

func (worker *Worker) Close() {
	worker.service.WaitGroup.Add(1);
	defer worker.service.WaitGroup.Done();
	worker.mut_state.Lock();
	defer worker.mut_state.Unlock();
	if worker.listen != nil {
		Log.Printf("[Shard-%d] Closing public listener..", worker.config.ShardIndex);
		if err := worker.listen.Close(); err != nil { Log.Printf("%v"); }
		worker.listen = nil;
	}
}



func (worker *Worker) Serve() {
	worker.service.WaitGroup.Add(1);
	defer func() {
		worker.Close();
		worker.link.Close();
		worker.service.WaitGroup.Done();
	}();
	packet_timeout, err := Time.ParseDuration(worker.config.ListenInterval);
		if err != nil { Log.Panic(Fmt.Errorf("%v for Listen-Interval", err)); }
	sync_interval,  err := Time.ParseDuration(worker.config.SyncInterval);
		if err != nil { Log.Panic(Fmt.Errorf("%v for Sync-Interval", err)); }
	shard_index := worker.config.ShardIndex;
	last_sync := Time.Now();
	chip_current := NewChip();
	LOOP_WORKER:
	for {
		if worker.link   .IsStopping() { break LOOP_WORKER; }
		if worker.service.IsStopping() { break LOOP_WORKER; }
		now := Time.Now();
		since := now.Sub(last_sync);
		// call sync
		if since >= sync_interval {
			last_sync = now;
			worker.DoSync(false, shard_index);
		}
		// call batch
		if worker.needs_batch.Load() {
			worker.needs_batch.Store(false);
			chip_batch := chip_current;
			chip_current = NewChip();
			worker.DoBatch(false, shard_index, chip_batch);
		}
		// listen for packets
		if err := worker.DoListen(shard_index, chip_current, packet_timeout); err != nil {
			// timeout
			if neterr, ok := err.(Net.Error); ok && neterr.Timeout() {
				continue LOOP_WORKER; }
			if neterr, ok := err.(*Net.OpError); ok &&
			neterr.Err.Error() == "use of closed network connection" {
				Log.Printf("[Shard-%d] Socket closed!", shard_index);
				break LOOP_WORKER; }
			Fmt.Printf("Unknown UDP listen error: %v", err);
			worker.PacketsError.Add(1);
			Utils.SleepX();
			continue LOOP_WORKER;
		}
	}
	if !worker.service.IsStopping() {
		Log.Printf("Something went wrong! Worker is exiting..");
		worker.service.Stop();
	}
	worker.Close();
	print("\n"); Log.Printf("[Shard-%d] Final sync..", shard_index);
	worker.DoSync(true, shard_index);
	print("\n"); Log.Printf("[Shard-%d] Batching last chip..", shard_index);
	worker.DoBatch(true, shard_index, chip_current);
}



func (worker *Worker) DoListen(shard_index uint8, chip *Chip,
		packet_timeout Time.Duration) error {
	// set listen timeout
	timeout := Time.Now().Add(packet_timeout);
	if err := worker.listen.SetReadDeadline(timeout); err != nil { return err; }
	buffer := make([]byte, 1500);
	n, addr, err := worker.listen.ReadFrom(buffer);
	if err != nil { return err; }
	// token bucket per ip
	host, _, err := Net.SplitHostPort(addr.String());
	if err != nil { return err; }
	ip_tup := PxnNet.ParseAddrStr(host);
	if ip_tup == nil { return Fmt.Errorf("Invalid IP: %s", host); }
	if worker.secret_db.CheckTupleIP(ip_tup) {
		worker.PacketsBlock.Add(1);
//TODO: send rate limited reply
//print("Rate Limited!\n");
		return nil;
	}
Fmt.Printf("Got %d bytes\n", n);
	// process packet
	reply, err := worker.Process(buffer[:n], &addr);
	if err != nil {
		worker.PacketsError.Add(1);
		Log.Printf("[Shard-%d] Packet Error: %v", shard_index, err);
		return nil;
	}
	// send reply
	if _, err := worker.listen.WriteTo(reply, addr); err != nil { return err; }
//TODO: remove this?
//		worker.PacketsError.Add(1);
//		Log.Printf("[Shard-%d] Socket Write Error: %v", shard_index, err);
//		return;
//	}
	worker.PacketsGood.Add(1);
	return nil;
}
