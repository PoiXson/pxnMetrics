package secret;
// pxnMetrics Broker - secret database

import(
	Log     "log"
//	Fmt     "fmt"
	Time    "time"
	PxnNet  "github.com/PoiXson/pxnGoCommon/net"
	RateLim "github.com/PoiXson/pxnGoCommon/net/ratelimit"
	Configs "github.com/PoiXson/pxnMetrics/broker/configs"
);



type SecretUpdater struct {
	config   *Configs.CfgBroker
	IPs      *RateLim.TokBuckLim
	ShardUps []*ShardUpdate
}

type ShardUpdate struct {
	UpIPs  map[string]int32
	UpUIDs []uint64
}



func New(config *Configs.CfgBroker) *SecretUpdater {
	num_shards := config.NumShards;
	shardups := make([]*ShardUpdate, num_shards);
	for index:=uint8(0); index<num_shards; index++ {
		shardups[index] = &ShardUpdate{
			UpIPs: make(map[string]int32),
		};
	}
	if config.RateLimit.TokenInterval == "" {
		Log.Panic("Token interval parameter is required"); }
	token_interval, err := Time.ParseDuration(config.RateLimit.TokenInterval);
	if err != nil { Log.Panicf("Invalid token interval: %s",
		config.RateLimit.TokenInterval); }
	rate_limit := RateLim.NewTokenBucket();
	rate_limit.Interval     = token_interval;
	rate_limit.HitCost      = int32(config.RateLimit.TokensPerHit);
	rate_limit.TokensThresh = int32(config.RateLimit.TokensThresh);
	rate_limit.TokensCap    = int32(config.RateLimit.TokensCap);
	rate_limit.Start();
	return &SecretUpdater{
		config:   config,
		IPs:      rate_limit,
		ShardUps: shardups,
	};
}



func (db *SecretUpdater) PushPull(index uint8,
		ips map[string]int32, uids []uint64,
		) ( map[string]int32,      []uint64) {
	db.IPs.MutBuckets.Lock();
	defer db.IPs.MutBuckets.Unlock();
	// share tokens
	{
		num_shards := db.config.NumShards;
		// ip token buckets
		for ip_str, tokens := range ips {
			ip_tup := PxnNet.ParseTupStr(ip_str);
			if ip_tup == nil { Log.Panicf(
				"Invalid IP in secret db: %s", ip_str); }
			// local database
			bucket := db.IPs.GetBucket(ip_tup);
			if bucket == nil {
				Log.Panic("Failed to get token bucket"); }
			bucket.Tokens += tokens;
			// other shards
			for idx:=uint8(0); idx<num_shards; idx++ {
				if idx != index {
					val, ok := db.ShardUps[idx].UpIPs[ip_str];
					if !ok { val = 0; }
//TODO: remove this
//Fmt.Printf("  pushing %d to %d   >%d + %d     %s\n", index, idx, val, tokens, ip_str);
					db.ShardUps[idx].UpIPs[ip_str] = val + tokens;
				}
			}
		}
		// server uid's
		for idx:=uint8(0); idx<num_shards; idx++ {
			if idx != index { db.ShardUps[idx].UpUIDs =
				append(db.ShardUps[idx].UpUIDs, uids...); }}
	}
	// pull for this shard
	{
		num_uids   := len(db.ShardUps[index].UpUIDs);
		reply_ips  := make(map[string]int32);
		reply_uids := make([]uint64, num_uids);
		// ip token buckets
		for ip_str, tokens := range db.ShardUps[index].UpIPs {
			reply_ips[ip_str] = tokens; }
		// server uid's
		for idx:=0; idx<num_uids; idx++ {
			reply_uids[idx] = db.ShardUps[index].UpUIDs[idx]; }
		// clear
		clear(db.ShardUps[index].UpIPs);
		db.ShardUps[index].UpUIDs = make([]uint64, 0);
		return reply_ips, reply_uids;
	}

}
