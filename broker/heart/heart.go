package heart;
// pxnMetrics Broker - heartbeat

import(
	Log      "log"
	Time     "time"
	Sync     "sync"
	PxnUtil  "github.com/PoiXson/pxnGoCommon/utils"
	PxnServ  "github.com/PoiXson/pxnGoCommon/service"
	Configs  "github.com/PoiXson/pxnMetrics/broker/configs"
	SecretDB "github.com/PoiXson/pxnMetrics/broker/databases/secret"
);



const HeartBeatInterval  = "200ms";
const HeartBeatStopLoops = 3;
const TaskQueueSize      = 1024;



type HeartBeat struct {
	Service   *PxnServ.Service
	Config    *Configs.CfgBroker
//TODO: nothing adds to the queue yet
	TaskQueue chan Task
	LastSync  Time.Time
	LastBatch Time.Time
	LastTask  Time.Time
	// shards
	MutShards Sync.Mutex
	Shards    []*ShardState
	NextShard uint8
	// IP/UID request counts
	SecretDB  *SecretDB.SecretUpdater
//TODO: data from batches?
	TotalRequests uint64 // all time
}

type ShardState struct {
	IsOnline   bool
	NeedsBatch bool
	LastSeen   Time.Time
	LastSync   Time.Time
	LastBatch  Time.Time
}



func New(service *PxnServ.Service, config *Configs.CfgBroker) *HeartBeat {
	num_shards := config.NumShards;
	shards := make([]*ShardState, num_shards);
	for index:=uint8(0); index<num_shards; index++ {
		shards[index] = &ShardState{}; }
	return &HeartBeat{
		Service:   service,
		Config:    config,
		TaskQueue: make(chan Task, TaskQueueSize),
		Shards:    shards,
		SecretDB:  SecretDB.New(config),
	};
}



func (heart *HeartBeat) Start() error {
	go heart.Serve();
	PxnUtil.SleepC();
	return nil;
}

func (heart *HeartBeat) Serve() {
	heart.Service.WaitGroup.Add(1);
	defer heart.Service.WaitGroup.Done();
	heart.Service.AddStopHook(func() {
		heart.MutShards.Lock();
		defer heart.MutShards.Unlock();
		for i:=uint8(0); i<heart.Config.NumShards; i++ {
			heart.Shards[i].NeedsBatch = false; }
	});
	num_shards := heart.Config.NumShards;
	Log.Printf("Starting HeartBeat..");
	loop_interval, err := Time.ParseDuration(HeartBeatInterval);
		if err != nil { panic(err); }
	batch_interval, err := Time.ParseDuration(heart.Config.BatchInterval);
		if err != nil { panic(err); }
	var batch_intr_per_shard Time.Duration;
	if num_shards > 0 {
		batch_intr_per_shard = batch_interval / Time.Duration(num_shards);
		Log.Printf("  Total Shards: %d", num_shards);
		Log.Printf("  Batch every: %s", batch_intr_per_shard);
	}
	timer := Time.NewTicker(loop_interval);
	last_batch := Time.Now();
	var stopping uint8 = 0;
	LOOP_SERVE:
	for { select {
		case task := <-heart.TaskQueue: {
//TODO: timing
			task.Run();
			heart.LastTask = Time.Now();
			LOOP_DRAIN:
			for { select {
				case <-timer.C: continue LOOP_DRAIN;
				default:        break    LOOP_DRAIN;
			}}
			stopping = 0;
			continue LOOP_SERVE;
		}
		case <-timer.C: {
			if heart.Service.IsStopping() {
				stopping++;
				if stopping >= HeartBeatStopLoops {
					break LOOP_SERVE; }
			} else
			if num_shards > 0 {
				now := Time.Now();
				// request next batch
				since_batch := now.Sub(last_batch);
				if since_batch >= batch_intr_per_shard {
					last_batch = now;
					next_index := heart.NextShard;
					heart.NextShard++;
					if heart.NextShard >= num_shards {
						heart.NextShard = 0; }
					heart.MutShards.Lock();
					shard := heart.Shards[next_index];
					if shard.IsOnline { shard.NeedsBatch = true;
					} else {            shard.NeedsBatch = false; }
					heart.MutShards.Unlock();
				}
			}
			continue LOOP_SERVE;
		}
	}}
	Log.Printf("Heart Beat stopped.");
}



func (heart *HeartBeat) QueueTask(task Task) {
	heart.TaskQueue <- task;
}



func (heart *HeartBeat) MarkSeen(index uint8) {
	heart.MutShards.Lock();
	defer heart.MutShards.Unlock();
	shard := heart.Shards[index];
	shard.IsOnline = true;
	shard.LastSeen = Time.Now();
}

func (heart *HeartBeat) MarkSynced(index uint8) bool {
	heart.MutShards.Lock();
	defer heart.MutShards.Unlock();
	shard := heart.Shards[index];
	shard.IsOnline = true;
	now := Time.Now();
	heart.LastSync = now;
	shard.LastSeen = now;
	shard.LastSync = now;
	if shard.NeedsBatch {
		shard.NeedsBatch = false;
		return true;
	}
	return false;
}

func (heart *HeartBeat) MarkBatched(index uint8) {
	heart.MutShards.Lock();
	defer heart.MutShards.Unlock();
	shard := heart.Shards[index];
	shard.IsOnline = true;
	now := Time.Now();
	heart.LastBatch = now;
	shard.LastSeen  = now;
	shard.LastBatch = now;
	shard.NeedsBatch = false;
}

func (heart *HeartBeat) MarkOffline(index uint8) {
	heart.MutShards.Lock();
	defer heart.MutShards.Unlock();
	shard := heart.Shards[index];
	shard.IsOnline = false;
	shard.LastSeen = Time.Now();
}
