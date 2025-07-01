package uplink;
// pxnMetrics Broker - shard api

import(
	Log       "log"
	Fmt       "fmt"
	Math      "math"
	Context   "context"
	Errors    "errors"
	GRPC      "google.golang.org/grpc"
	GStatus   "google.golang.org/grpc/status"
	GCodes    "google.golang.org/grpc/codes"
	PxnRPC    "github.com/PoiXson/pxnGoCommon/rpc"
	Configs   "github.com/PoiXson/pxnMetrics/broker/configs"
	Heart     "github.com/PoiXson/pxnMetrics/broker/heart"
	UserMan   "github.com/PoiXson/pxnMetrics/broker/userman"
	API_Shard "github.com/PoiXson/pxnMetrics/api/shard"
	GEmpty    "google.golang.org/protobuf/types/known/emptypb"
);



type BrokerShardAPI struct {
	API_Shard.UnimplementedServiceShardAPIServer
	Config *Configs.CfgBroker
	heart  *Heart.HeartBeat
}



func NewShardAPI(rpc *GRPC.Server, config *Configs.CfgBroker,
		heart *Heart.HeartBeat) *BrokerShardAPI {
	api := BrokerShardAPI{
		Config: config,
		heart:  heart,
	};
	API_Shard.RegisterServiceShardAPIServer(rpc, &api);
	return &api;
}



func (api *BrokerShardAPI) Greet(ctx Context.Context,
		hello *API_Shard.Hello) (*API_Shard.Hey, error) {
	username := ctx.Value(PxnRPC.KeyUsername).(string);
	if username == "" {
		Log.Printf("Invalid RPC user");
		return nil, Errors.New("Invalid RPC user");
	}
	user, ok := ctx.Value(UserMan.KeyUserRPC).(*Configs.CfgUser);
	if !ok {
		Log.Printf("Invalid RPC User type");
		return nil, Errors.New("Invalid RPC User type");
	}
	if user == nil {
		msg := Fmt.Sprintf("Invalid RPC user: %s", username);
		Log.Printf(msg);
		return nil, Errors.New(msg);
	}
	shard_index := uint8(hello.ShardIndex);
	// check user to shard index
	found := false;
	for _, permit_index := range user.PermitShards {
		if permit_index == shard_index {
			found = true;
			break;
		}
	}
	if !found {
		msg := Fmt.Sprintf("User %s cannot serve shard %d",
			username, shard_index);
		Log.Printf(msg);
		return nil, GStatus.Error(GCodes.PermissionDenied, msg);
	}
	Log.Printf("[Shard-%d] Shard Connected! Index %d of %d",
		shard_index, shard_index, api.Config.NumShards);
	return &API_Shard.Hey{
		NumShards:      uint32(api.Config.NumShards),
		ChecksumBase:   uint32(api.Config.ChecksumBase),
		ListenInterval: api.Config.ListenInterval,
		SyncInterval:   api.Config.SyncInterval,
		// token bucket
		TokenInterval:  api.Config.RateLimit.TokenInterval,
		TokensPerHit:   uint32(api.Config.RateLimit.TokensPerHit),
		TokensThresh:   uint32(api.Config.RateLimit.TokensThresh),
		TokensCap:      uint32(api.Config.RateLimit.TokensCap),
	}, nil;
}



func (api *BrokerShardAPI) SyncDBs(ctx Context.Context, send *API_Shard.SyncSend,
		) (*API_Shard.SyncReply, error) {
	if send.ShardIndex == 0 || send.ShardIndex > Math.MaxUint8 {
		Log.Panicf("Invalid shard index: %d", send.ShardIndex); }
	index := uint8(send.ShardIndex);
//TODO: remove this
//Fmt.Printf("SYNC %d\n", index);
	secretdb := api.heart.SecretDB;
	ips, ids := secretdb.PushPull(
		index,
		send.TokenBuckets,
		send.ServerUIDs,
	);
	needs_batch := api.heart.MarkSynced(index);
	return &API_Shard.SyncReply{
		NeedsBatch:   needs_batch,
		TokenBuckets: ips,
		ServerUIDs:   ids,
	}, nil;
}



func (api *BrokerShardAPI) BatchOut(ctx Context.Context, batch_data *API_Shard.BatchData,
		) (*GEmpty.Empty, error) {
	if batch_data.ShardIndex == 0 || batch_data.ShardIndex > Math.MaxUint8 {
		Log.Panicf("Invalid shard index: %d", batch_data.ShardIndex); }
	index := uint8(batch_data.ShardIndex);
	api.heart.Shards[index].NeedsBatch = false;
	task := Heart.NewTask_BatchOut(
//TODO
	);
	api.heart.QueueTask(task);
	// shard going offline
	if batch_data.IsLast { api.heart.MarkOffline(index);
Fmt.Printf("Going Offline: %d\n", index);
	}
	return &GEmpty.Empty{}, nil;
}
