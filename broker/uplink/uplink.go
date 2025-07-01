package uplink;
// pxnMetrics Broker - rpc server

import(
	GRPC     "google.golang.org/grpc"
	UtilsRPC "github.com/PoiXson/pxnGoCommon/rpc"
	Service  "github.com/PoiXson/pxnGoCommon/service"
	Configs  "github.com/PoiXson/pxnMetrics/broker/configs"
	Heart    "github.com/PoiXson/pxnMetrics/broker/heart"
	UserMan  "github.com/PoiXson/pxnMetrics/broker/userman"
);



type UpLink struct {
	service *Service.Service
	config  *Configs.CfgBroker
	heart   *Heart.HeartBeat
	rpc     *UtilsRPC.ServerRPC
	// api's
	API_Shard *BrokerShardAPI
	API_Front *BrokerFrontAPI
}



func New(service *Service.Service, config *Configs.CfgBroker,
		heart *Heart.HeartBeat) *UpLink {
	rpc := UtilsRPC.NewServerRPC(service, config.BindRPC);
	return &UpLink{
		service: service,
		config:  config,
		heart:   heart,
		rpc:     rpc,
	};
}



func (uplink *UpLink) Start() error {
	allow_ips := make(map[string]string);
	for username, user := range uplink.config.Users {
		for _, ip := range user.PermitIPs {
			allow_ips[ip] = username;
		}
	}
	uplink.rpc.SetServerGRPC(GRPC.NewServer(
		GRPC.ChainUnaryInterceptor(
			UtilsRPC.NewAuthByIP(allow_ips),
			UserMan.NewUserManagerInterceptor(uplink.config),
		),
	));
	// api's
	uplink.API_Shard = NewShardAPI(
		uplink.rpc.GetServerGRPC(),
		uplink.config,
		uplink.heart,
	);
	uplink.API_Front = NewFrontAPI(
		uplink.rpc.GetServerGRPC(),
		uplink.config,
		uplink.heart,
	);
	return uplink.rpc.Start();
}
