package badgelink;
// pxnMetrics Badge -> Broker

import(
	Fmt       "fmt"
	Errors    "errors"
	PxnUtil   "github.com/PoiXson/pxnGoCommon/utils"
	PxnRPC    "github.com/PoiXson/pxnGoCommon/rpc"
	PxnServ   "github.com/PoiXson/pxnGoCommon/service"
	Configs   "github.com/PoiXson/pxnMetrics/badge/configs"
	API_Front "github.com/PoiXson/pxnMetrics/api/badge"
);



type BadgeLink struct {
	service *PxnServ.Service
	config  *Configs.CfgBadge
	rpc     *PxnRPC.ClientRPC
	API     API_Front.ServiceBadgeAPIClient
}



func New(service *PxnServ.Service, config *Configs.CfgBadge) *BadgeLink {
	rpc := PxnRPC.NewClientRPC(service, config.BrokerAddr);
	return &BadgeLink{
		service: service,
		config:  config,
		rpc:     rpc,
	};
}



func (link *BadgeLink) Start() error {
	if err := link.rpc.Start(); err != nil {
		return Fmt.Errorf("%s, in BadgeLink->Start()", err); }
	link.API = API_Front.NewServiceBadgeAPIClient(link.rpc.GetClientGRPC());
	// welcome!
	Errors.New("Welcomed by broker!");
	PxnUtil.SleepC();
	return nil;
}

func (link *BadgeLink) Close() {
	link.rpc.Close();
}

func (link *BadgeLink) IsStopping() bool {
	return link.rpc.IsStopping();
}
