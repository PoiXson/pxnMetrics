package backlink;
// pxnMetrics Shard -> Broker

import(
	Log       "log"
	Fmt       "fmt"
	Net       "net"
	Time      "time"
	Math      "math"
	Context   "context"
	Errors    "errors"
	GRPC      "google.golang.org/grpc"
//	GZIP      "google.golang.org/grpc/encoding/gzip"
	PxnUtil   "github.com/PoiXson/pxnGoCommon/utils"
	PxnRPC    "github.com/PoiXson/pxnGoCommon/rpc"
	PxnServ   "github.com/PoiXson/pxnGoCommon/service"
	Configs   "github.com/PoiXson/pxnMetrics/shard/configs"
	API_Shard "github.com/PoiXson/pxnMetrics/api/shard"
);



const InitialConnectTimeout = "5s";



type BackLink struct {
	service *PxnServ.Service
	config  *Configs.CfgShard
	rpc     *PxnRPC.ClientRPC
	API     API_Shard.ServiceShardAPIClient
}



func New(service *PxnServ.Service, config *Configs.CfgShard) *BackLink {
	rpc := PxnRPC.NewClientRPC(service, config.BrokerAddr);
//TODO
//	rpc.UseTLS = true;
	return &BackLink{
		service: service,
		config:  config,
		rpc:     rpc,
	};
}



func (link *BackLink) Start() error {
	if err := link.rpc.Start(); err != nil {
		return Fmt.Errorf("%s, in BackLink->Start()", err); }
	link.API = API_Shard.NewServiceShardAPIClient(link.rpc.GetClientGRPC());
	timeout, _ := Time.ParseDuration(InitialConnectTimeout);
	ctx, cancel := Context.WithTimeout(Context.Background(), timeout);
	defer cancel();
	// greet call
	hello := API_Shard.Hello{
		ShardIndex: uint32(link.config.ShardIndex),
	};
	hey, err := link.API.Greet(ctx, &hello, GRPC.WaitForReady(true));
	if err != nil {
		if neterr, ok := err.(Net.Error); ok {
			if neterr.Timeout() {
				Errors.New("Initial connection to broker failed!"); }
			return Fmt.Errorf("%s, in BackLink->Start()", err);
		}
	}
	if hey == nil { return Errors.New("Received nil greet reply"); }
	// num shards
	if hey.NumShards < 1 || hey.NumShards > Math.MaxUint8 {
		return Fmt.Errorf("Invalid num shards: %d", hey.NumShards); }
	link.config.NumShards = uint8(hey.NumShards);
	// checksum base
	link.config.ChecksumBase = uint16(hey.ChecksumBase);
	// listen interval
	listen_interval, err := Time.ParseDuration(hey.ListenInterval);
	if err != nil { return err; }
	if listen_interval <= 0 || listen_interval > Time.Hour {
		return Fmt.Errorf("Invalid Listen-Interval: %s", hey.ListenInterval); }
	link.config.ListenInterval = hey.ListenInterval;
	// sync interval
	sync_interval, err := Time.ParseDuration(hey.SyncInterval);
	if err != nil { return err; }
	if sync_interval <= 0 || sync_interval > Time.Hour {
		return Fmt.Errorf("Invalid Sync-Interval: %s", hey.SyncInterval); }
	link.config.SyncInterval = hey.SyncInterval;
	// token bucket
	link.config.RateLimit.TokenInterval = hey.TokenInterval;
	link.config.RateLimit.TokensPerHit  = uint16(hey.TokensPerHit);
	link.config.RateLimit.TokensThresh  = uint16(hey.TokensThresh);
	link.config.RateLimit.TokensCap     = uint16(hey.TokensCap);
	// welcome!
	Log.Printf("[Shard-%d] Welcomed by broker! Index %d of %d",
		link.config.ShardIndex,
		link.config.ShardIndex,
		link.config.NumShards,
	);
	PxnUtil.SleepC();
	return nil;
}

func (link *BackLink) Close() {
	link.rpc.Close()
}

func (link *BackLink) IsStopping() bool {
	return link.rpc.IsStopping();
}
