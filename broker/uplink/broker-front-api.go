package uplink;
// pxnMetrics Broker - Frontend API

import(
	Log       "log"
	Fmt       "fmt"
	Time      "time"
	Context   "context"
	JSON      "encoding/json"
	Errors    "errors"
	GRPC      "google.golang.org/grpc"
	GEmpty    "google.golang.org/protobuf/types/known/emptypb"
	UtilsRPC  "github.com/PoiXson/pxnGoCommon/rpc"
	Configs   "github.com/PoiXson/pxnMetrics/broker/configs"
	Heart     "github.com/PoiXson/pxnMetrics/broker/heart"
	UserMan   "github.com/PoiXson/pxnMetrics/broker/userman"
	API_Front "github.com/PoiXson/pxnMetrics/api/front"
	API_Web   "github.com/PoiXson/pxnMetrics/api/web"
);



type BrokerFrontAPI struct {
	API_Front.UnimplementedServiceFrontendAPIServer
	Config *Configs.CfgBroker
	Heart  *Heart.HeartBeat
}



func NewFrontAPI(rpc *GRPC.Server, config *Configs.CfgBroker,
		heart *Heart.HeartBeat) *BrokerFrontAPI {
	api := BrokerFrontAPI{
		Config: config,
		Heart:  heart,
	};
	API_Front.RegisterServiceFrontendAPIServer(rpc, &api);
	return &api;
}



func (api *BrokerFrontAPI) FetchStatusJSON(ctx Context.Context,
		_ *GEmpty.Empty) (*API_Front.StatusJSON, error) {
	username := ctx.Value(UtilsRPC.KeyUsername).(string);
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
	if !user.PermitWeb {
		msg := Fmt.Sprintf("User lacks permissions: %s", username);
		Log.Printf(msg);
		return nil, Errors.New(msg);
	}
	now := Time.Now();
	// broker status
	var last_sync  int64;
	var last_batch int64;
	var last_task  int64;
	if api.Heart.LastSync.IsZero() { last_sync = 0;
	} else { last_sync = api.Heart.LastSync.Unix(); }
	if api.Heart.LastBatch.IsZero() { last_batch = 0;
	} else { last_batch = api.Heart.LastBatch.Unix(); }
	if api.Heart.LastTask.IsZero() { last_task = 0;
	} else { last_task = api.Heart.LastTask.Unix(); }
	json_broker := API_Web.BrokerStatus{
		Status: API_Web.ServerStatus{
			Name:      "Broker",
			Desc:      "",
			LastSeen:  uint32 (now.Unix()),
			LastSync:  uint32 (last_sync ),
			LastBatch: uint32 (last_batch),
//			LastTask:  uint32 (last_task ),
		},
		LastTask: uint32(last_task),
	};
//TODO
json_broker.Status.Status = "Online";
//	if shard.IsOnline { json_shards[index].Status = "Online";
//	} else {            json_shards[index].Status = "Offline" }
	// shard status
	num_shards := api.Config.NumShards;
	json_shards := make([]API_Web.ServerStatus, num_shards);
	for index:=uint8(0); index<num_shards; index++ {
		shard := api.Heart.Shards[index];
		var last_seen  int64;
		var last_sync  int64;
		var last_batch int64;
		if shard.LastSeen.IsZero() { last_seen = 0;
		} else { last_seen = shard.LastSeen.Unix(); }
		if shard.LastSync.IsZero() { last_sync = 0;
		} else { last_sync = shard.LastSync.Unix(); }
		if shard.LastBatch.IsZero() { last_batch = 0;
		} else { last_batch = shard.LastBatch.Unix(); }
		json_shards[index] = API_Web.ServerStatus{
			Name:      Fmt.Sprintf("Shard-%d", index+1),
			Desc:      "",
			LastSeen:  uint32 (last_seen ),
			LastSync:  uint32 (last_sync ),
			LastBatch: uint32 (last_batch),
//TODO
//			BatchWaiting: uint32 (shard.BatchWaiting),
//			QueueWaiting: uint32 (shard.QueueWaiting),
//			ReqPerSec:    float32(shard.ReqPerSec   ),
//			ReqPerMin:    float32(shard.ReqPerMin   ),
//			ReqPerHour:   float32(shard.ReqPerHour  ),
//			ReqPerDay:    float32(shard.ReqPerDay   ),
//			ReqTotal:     uint64 (shard.ReqTotal    ),
		};
		if shard.IsOnline { json_shards[index].Status = "Online";
		} else {            json_shards[index].Status = "Offline" }
	}
	json, err := JSON.Marshal(
		API_Web.Status{
			Broker: json_broker,
			Shards: json_shards,
		},
	);
	if err != nil { return nil, err; }
	return &API_Front.StatusJSON{ Data: json }, nil;
}
