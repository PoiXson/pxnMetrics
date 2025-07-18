package app;
// badge.minecraftmetrics.com

import(
	Log     "log"
	Flag    "flag"
	Flagz   "github.com/PoiXson/pxnGoCommon/utils/flagz"
	PxnFS   "github.com/PoiXson/pxnGoCommon/utils/fs"
	PxnUtil "github.com/PoiXson/pxnGoCommon/utils"
	PxnWeb  "github.com/PoiXson/pxnGoCommon/net/web"
	PxnServ "github.com/PoiXson/pxnGoCommon/service"
	BagLink "github.com/PoiXson/pxnMetrics/badge/badgelink"
	Configs "github.com/PoiXson/pxnMetrics/badge/configs"
	Pages   "github.com/PoiXson/pxnMetrics/badge/pages"
);



type AppBadge struct {
	service *PxnServ.Service
	websvr  *PxnWeb.WebServer
	pages   *Pages.Pages
	link    *BagLink.BadgeLink
	config  *Configs.CfgBadge
}



func New() PxnServ.AppFace {
	return &AppBadge{};
}

func (app *AppBadge) Main() {
	app.service = PxnServ.New();
	app.service.Start();
	app.flags_and_configs(DefaultConfigFile);
	// rpc
	app.link = BagLink.New(app.service, app.config.BrokerAddr);
	// web server
	app.websvr = PxnWeb.NewWebServer(
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
		PxnUtil.SleepCn(5);
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
	cfg, err := PxnFS.LoadConfig[Configs.CfgBadge](file);
	if err != nil { Log.Panicf("%s, when loading config %s", err, file); }
	// remote rpc
	if flag_broker    != "" { cfg.BrokerAddr = flag_broker;          }
	if cfg.BrokerAddr == "" { cfg.BrokerAddr = DefaultBrokerAddress; }
	// bind web
	if flag_bindweb != "" { cfg.BindWeb = flag_bindweb;          }
	if cfg.BindWeb  == "" { cfg.BindWeb = PxnWeb.DefaultBindWeb; }
	app.config = cfg;
}
