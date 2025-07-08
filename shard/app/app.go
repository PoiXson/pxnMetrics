package app;
// pxnMetrics Shard

import(
	OS       "os"
	Log      "log"
	Fmt      "fmt"
	Math     "math"
	Flag     "flag"
	Flagz    "github.com/PoiXson/pxnGoCommon/utils/flagz"
	PxnFS    "github.com/PoiXson/pxnGoCommon/utils/fs"
	PxnServ  "github.com/PoiXson/pxnGoCommon/service"
	BackLink "github.com/PoiXson/pxnMetrics/shard/backlink"
	Worker   "github.com/PoiXson/pxnMetrics/shard/worker"
	Configs  "github.com/PoiXson/pxnMetrics/shard/configs"
);



type AppShard struct {
	service *PxnServ.Service
	link    *BackLink.BackLink
	worker  *Worker.Worker
	config  *Configs.CfgShard
}



func New() PxnServ.AppFace {
	return &AppShard{};
}

func (app *AppShard) Main() {
	app.service = PxnServ.New();
	app.service.Start();
	app.flags_and_configs(DefaultConfigFile);
	// rpc
	app.link = BackLink.New(app.service, app.config);
	// public listener
	app.worker = Worker.New(app.service, app.config, app.link);
	// start things
	if err := app.link  .Start(); err != nil { Log.Panic(err); }
	if err := app.worker.Start(); err != nil { Log.Panic(err); }
	app.service.WaitUntilEnd();
}



func (app *AppShard) flags_and_configs(file string) {
	var flag_broker     string;
	var flag_bindsubmit string;
	var flag_shardindex int;
	Flagz.String(&flag_broker,     "broker",      "");
	Flagz.String(&flag_bindsubmit, "bind",        "");
	Flagz.Int   (&flag_shardindex, "shard-index", -1);
	Flag.Parse();
	// load config
	cfg, err := PxnFS.LoadConfig[Configs.CfgShard](file);
	if err != nil { Log.Panicf("%s, when loading config %s", err, file); }
	// remote rpc
	if flag_broker    != "" { cfg.BrokerAddr = flag_broker;          }
	if cfg.BrokerAddr == "" { cfg.BrokerAddr = DefaultBrokerAddress; }
	// bind submit service
	if flag_bindsubmit != "" { cfg.BindPublic = flag_bindsubmit;   }
	if cfg.BindPublic  == "" { cfg.BindPublic = DefaultBindPublic; }
	// shard index
	if flag_shardindex == -1 {
		Fmt.Printf("--shard-index flag is required\n");
		Flag.Usage();
		OS.Exit(1);
	}
	if flag_shardindex < 1 || flag_shardindex > Math.MaxUint8 {
		Log.Panicf("Invalid shard index; max %d", Math.MaxUint8); }
	cfg.ShardIndex = uint8(flag_shardindex);
	app.config = cfg;
}
