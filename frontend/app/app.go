package app;
// minecraftmetrics.com

import(
	Log     "log"
	Flag    "flag"
	Flagz   "github.com/PoiXson/pxnGoCommon/utils/flagz"
	Utils   "github.com/PoiXson/pxnGoCommon/utils"
	UtilsFS "github.com/PoiXson/pxnGoCommon/utils/fs"
	PxnNet  "github.com/PoiXson/pxnGoCommon/utils/net"
	Service "github.com/PoiXson/pxnGoCommon/service"
	WebLink "github.com/PoiXson/pxnMetrics/frontend/weblink"
	Configs "github.com/PoiXson/pxnMetrics/frontend/configs"
	PagesMetricsCom "github.com/PoiXson/pxnMetrics/frontend/pages/metrics/com"
	PagesMetricsTop "github.com/PoiXson/pxnMetrics/frontend/pages/metrics/top"
);



const DomainMetricsCom = "minecraftmetrics.com";
const DomainMetricsTop = "minecraftmetrics.top";



type AppFront struct {
	service   *Service.Service
	websvr    *PxnNet.WebServer
	pages_com *PagesMetricsCom.Pages
	pages_top *PagesMetricsTop.Pages
	link      *WebLink.WebLink
	config    *Configs.CfgFront
}



func New() Service.AppFace {
	return &AppFront{};
}

func (app *AppFront) Main() {
	app.service = Service.New();
	app.service.Start();
	app.flags_and_configs(DefaultConfigFile);
	// rpc
	app.link = WebLink.New(app.service, app.config.BrokerAddr);
	// web server
	app.websvr = PxnNet.NewWebServer(
		app.service,
		app.config.BindWeb,
		app.config.Proxied,
	);
	router := PxnNet.NewDomainsRouter();
	app.websvr.Router = router;
	router_metrics_com := router.DefDomain(DomainMetricsCom, true);
	router_metrics_top := router.AddDomain(DomainMetricsTop, true);
	app.pages_com = PagesMetricsCom.New(app.link, router_metrics_com);
	app.pages_top = PagesMetricsTop.New(app.link, router_metrics_top);
	// start things
	if err := app.link  .Start(); err != nil { Log.Panic(err); }
	if err := app.websvr.Start(); err != nil { Log.Panic(err); }
	// delay rpc close
	app.service.AddStopHook(func() { go func() {
		Utils.SleepCn(5);
		app.link.Close();
	}(); });
	app.service.WaitUntilEnd();
}



func (app *AppFront) flags_and_configs(file string) {
	var flag_broker  string;
	var flag_bindweb string;
	var flag_proxied bool;
	Flagz.String(&flag_broker,  "broker", "");
	Flagz.String(&flag_bindweb, "bind",   "");
	Flagz.Bool  (&flag_proxied, "proxied"   );
	Flag.Parse();
	// load config
	cfg, err := UtilsFS.LoadConfig[Configs.CfgFront](file);
	if err != nil { Log.Panicf("%s, when loading config %s", err, file); }
	// remote rpc
	if flag_broker    != "" { cfg.BrokerAddr = flag_broker;          }
	if cfg.BrokerAddr == "" { cfg.BrokerAddr = DefaultBrokerAddress; }
	// bind web
	if flag_bindweb != "" { cfg.BindWeb = flag_bindweb;          }
	if cfg.BindWeb  == "" { cfg.BindWeb = PxnNet.DefaultBindWeb; }
	if flag_proxied       { app.config.Proxied = true;           }
	app.config = cfg;
}
