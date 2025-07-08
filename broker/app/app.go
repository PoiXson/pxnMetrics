package app;
// pxnMetrics Broker

import(
	Log     "log"
	Fmt     "fmt"
	Time    "time"
	Math    "math"
	Flag    "flag"
	Flagz   "github.com/PoiXson/pxnGoCommon/utils/flagz"
	PxnFS   "github.com/PoiXson/pxnGoCommon/utils/fs"
	PxnServ "github.com/PoiXson/pxnGoCommon/service"
	Configs "github.com/PoiXson/pxnMetrics/broker/configs"
	Heart   "github.com/PoiXson/pxnMetrics/broker/heart"
	UpLink  "github.com/PoiXson/pxnMetrics/broker/uplink"
);



type AppBroker struct {
	service *PxnServ.Service
	heart   *Heart.HeartBeat
	link    *UpLink.UpLink
	config  *Configs.CfgBroker
}



func New() PxnServ.AppFace {
	return &AppBroker{};
}

func (app *AppBroker) Main() {
	app.service = PxnServ.New();
	app.service.Start();
	app.flags_and_configs(DefaultConfigFile);
	// databases
//TODO
	// heartbeat
	app.heart = Heart.New(app.service, app.config);
	if app.config.NumShards == 0 {
		Log.Printf("%sShard brokering is disabled", LogPrefix);
	} else {
		s := ""; if app.config.NumShards != 1 { s = "s"; }
		Log.Printf("%sListening for %d shard%s..",
			LogPrefix, app.config.NumShards, s);
	}
	// rpc
	app.link = UpLink.New(app.service, app.config, app.heart);
	// start things
//TODO: db
	if err := app.heart.Start(); err != nil {
		Log.Panicf("%s, when starting heartbeat", err); }
	if err := app.link.Start(); err != nil {
		Log.Panicf("%s, when starting uplink",    err); }
	app.service.WaitUntilEnd();
}



func (app *AppBroker) flags_and_configs(file string) {
	var flag_bind           string;
	var flag_checksum       int;
	var flag_interval_sync  string;
	var flag_interval_batch string;
	Flagz.String(&flag_bind,           "bind",           "");
	Flagz.Int   (&flag_checksum,       "checksum",       -1);
	Flagz.String(&flag_interval_sync,  "sync-interval",  "");
	Flagz.String(&flag_interval_batch, "batch-interval", "");
	Flag.Parse();
	// load config
	cfg, err := PxnFS.LoadConfig[Configs.CfgBroker](file);
	if err != nil { Log.Panicf("%s, when loading config %s", err, file); }
	// bind rpc
	if flag_bind   != "" { cfg.BindRPC = flag_bind;      }
	if cfg.BindRPC == "" { cfg.BindRPC = DefaultBindRPC; }
	// checksum base
	if flag_checksum > Math.MaxUint16 {
		Log.Panicf("Invalid checksum base value: ", flag_checksum); }
	if flag_checksum   >= 0 { cfg.ChecksumBase = uint16(flag_checksum);       }
	if cfg.ChecksumBase < 0 { cfg.ChecksumBase = uint16(DefaultChecksumBase); }
	// listen interval
	if cfg.ListenInterval == "" { cfg.ListenInterval = DefaultListenInterval; }
	listen_interval, err := Time.ParseDuration(cfg.ListenInterval);
	if err != nil { Log.Panic(err); }
	if listen_interval <= 0 || listen_interval > Time.Minute {
		Log.Panicf("Invalid listen-interval: %s", cfg.ListenInterval); }
	// sync interval
	if flag_interval_sync != "" { cfg.SyncInterval = flag_interval_sync;  }
	if cfg.SyncInterval   == "" { cfg.SyncInterval = DefaultSyncInterval; }
	sync_interval, err := Time.ParseDuration(cfg.SyncInterval);
	if err != nil { Log.Panic(err); }
	if sync_interval <= 0 || sync_interval > Time.Hour {
		Log.Panicf("Invalid sync-interval: %s", cfg.SyncInterval); }
	// batch interval
	if flag_interval_batch != "" { cfg.BatchInterval = flag_interval_batch;  }
	if cfg.BatchInterval   == "" { cfg.BatchInterval = DefaultBatchInterval; }
	batch_interval, err := Time.ParseDuration(cfg.BatchInterval);
	if err != nil { Log.Panic(err); }
	if batch_interval <= 0 || batch_interval > Time.Hour {
		Log.Panicf("Invalid batch-interval: %s", cfg.BatchInterval); }
	// rate limiter
	if cfg.RateLimit.TokenInterval =="" { cfg.RateLimit.TokenInterval = DefaultTokenInterval; }
	if cfg.RateLimit.TokensPerHit  == 0 { cfg.RateLimit.TokensPerHit  = DefaultTokensPerHit;  }
	if cfg.RateLimit.TokensThresh  == 0 { cfg.RateLimit.TokensThresh  = DefaultTokensThresh;  }
	if cfg.RateLimit.TokensCap     == 0 { cfg.RateLimit.TokensCap     = DefaultTokensCap;     }
	// users
	Fmt.Print("Users:\n");
	max_index := 0;
	for username, user := range cfg.Users {
		line := "";
		if user.PermitWeb { line = "Web"; }
		for _, index := range user.PermitShards {
			if line != "" { line += ", "; }
			line += Fmt.Sprintf("%d", index);
			if max_index < int(index) { max_index = int(index); }
		}
		Fmt.Printf("  <%s>  %s\n", username, line);
		for _, ip := range user.PermitIPs {
			Fmt.Printf("     %s\n", ip); }
	}
	if len(cfg.Users) == 0 { panic("No users configured!"); }
	if max_index > Math.MaxUint8 { Log.Panic("Invalid number of shards: ", max_index); }
	cfg.NumShards = uint8(max_index);
	app.config = cfg;
}
