package app;
// badge.minecraftmetrics.com

import(
	Log     "log"
	Flag    "flag"
	Flagz   "github.com/PoiXson/pxnGoCommon/utils/flagz"
	Utils   "github.com/PoiXson/pxnGoCommon/utils"
	UtilsFS "github.com/PoiXson/pxnGoCommon/utils/fs"
	PxnNet  "github.com/PoiXson/pxnGoCommon/utils/net"
	Service "github.com/PoiXson/pxnGoCommon/service"
	WebLink "github.com/PoiXson/pxnMetrics/frontend/weblink"
	Configs "github.com/PoiXson/pxnMetrics/badge/configs"
	Pages   "github.com/PoiXson/pxnMetrics/badge/pages"
);



type AppBadge struct {
	service *Service.Service
	websvr  *PxnNet.WebServer
	pages   *Pages.Pages
	link    *WebLink.WebLink
	config  *Configs.CfgBadge
}



func New() Service.AppFace {
	return &AppBadge{};
}

func (app *AppBadge) Main() {
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
	router := app.websvr.WithGorilla();
	app.pages = Pages.New(app.link, router);
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



func (app *AppBadge) flags_and_configs(file string) {
	var flag_broker  string;
	var flag_bindweb string;
	Flagz.String(&flag_broker,  "broker", "");
	Flagz.String(&flag_bindweb, "bind",   "");
	Flag.Parse();
	// load config
	cfg, err := UtilsFS.LoadConfig[Configs.CfgBadge](file);
	if err != nil { Log.Panicf("%s, when loading config %s", err, file); }
	// remote rpc
	if flag_broker    != "" { cfg.BrokerAddr = flag_broker;          }
	if cfg.BrokerAddr == "" { cfg.BrokerAddr = DefaultBrokerAddress; }
	// bind web
	if flag_bindweb != "" { cfg.BindWeb = flag_bindweb;          }
	if cfg.BindWeb  == "" { cfg.BindWeb = PxnNet.DefaultBindWeb; }
	app.config = cfg;
}
