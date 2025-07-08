package weblink;
// pxnMetrics Frontend -> Broker

import(
	Fmt       "fmt"
	Errors    "errors"
	PxnUtil   "github.com/PoiXson/pxnGoCommon/utils"
	PxnRPC    "github.com/PoiXson/pxnGoCommon/rpc"
	PxnServ   "github.com/PoiXson/pxnGoCommon/service"
	Configs   "github.com/PoiXson/pxnMetrics/frontend/configs"
	API_Front "github.com/PoiXson/pxnMetrics/api/front"
);



type WebLink struct {
	service *PxnServ.Service
	config  *Configs.CfgFront
	rpc     *PxnRPC.ClientRPC
	API     API_Front.ServiceFrontendAPIClient
}



func New(service *PxnServ.Service, config *Configs.CfgFront) *WebLink {
	rpc := PxnRPC.NewClientRPC(service, config.BrokerAddr);
	return &WebLink{
		service: service,
		config:  config,
		rpc:     rpc,
	};
}



func (link *WebLink) Start() error {
	if err := link.rpc.Start(); err != nil {
		return Fmt.Errorf("%s, in WebLink->Start()", err); }
	link.API = API_Front.NewServiceFrontendAPIClient(link.rpc.GetClientGRPC());
	// welcome!
	Errors.New("Welcomed by broker!");
	PxnUtil.SleepC();
	return nil;
}

func (link *WebLink) Close() {
	link.rpc.Close();
}

func (link *WebLink) IsStopping() bool {
	return link.rpc.IsStopping();
}
