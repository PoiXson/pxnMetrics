package secretdb;
// pxnMetrics Shard - secret database

import(
	Log     "log"
	Time    "time"
	Sync    "sync"
	PxnNet  "github.com/PoiXson/pxnGoCommon/utils/net"
	RateLim "github.com/PoiXson/pxnGoCommon/utils/net/ratelimit"
	Configs "github.com/PoiXson/pxnMetrics/shard/configs"
);



type DB struct {
	mut_db  Sync.Mutex
	config  *Configs.CfgShard
	IPs     *RateLim.TokBuckLim
	IPsNew  map[PxnNet.TupIP]int32
	UIDs    map[uint64]bool
	UIDsNew []uint64
}



func New(config *Configs.CfgShard) *DB {
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
	return &DB{
		config: config,
		IPs:    rate_limit,
		IPsNew: make(map[PxnNet.TupIP]int32),
	};
}



func (db *DB) CheckTupleIP(ip_tup *PxnNet.TupIP) bool {
	db.mut_db.Lock();
	defer db.mut_db.Unlock();
	blocked  := db.IPs.CheckTupleIP(ip_tup);
	hit_cost := db.IPs.HitCost;
	val, ok  := db.IPsNew[*ip_tup];
	if !ok { val = 0; }
	db.IPsNew[*ip_tup] = val + hit_cost;
	return blocked;
}



// rotates and returns
func (db *DB) UpdatePush() (map[string]int32, []uint64) {
	db.mut_db.Lock();
	defer db.mut_db.Unlock();
	// ip rate limit
	buckets_new := make(map[string]int32);
	for ip_tup, tokens := range db.IPsNew {
		buckets_new[ip_tup.String()] = tokens; }
	clear(db.IPsNew);
	// server uids
	num_new_uids := len(db.UIDsNew);
	uids_new := make([]uint64, num_new_uids);
	for i:=0; i<num_new_uids; i++ {
		uids_new[i] = db.UIDsNew[i]; }
	clear(db.UIDsNew);
	return buckets_new, uids_new;
}

// merges data to local
func (db *DB) UpdatePull(buckets map[string]int32, uids []uint64) {
	db.mut_db.Lock();
	defer db.mut_db.Unlock();
	// merge ip tokens
	for ip_str, tokens := range buckets {
		ip_tup := PxnNet.ParseTupStr(ip_str);
		if ip_tup == nil { Log.Panicf(
			"Invalid IP in secret db: %s", ip_str); }
		bucket := db.IPs.GetBucket(ip_tup);
		if bucket == nil { Log.Panic("Failed to get token bucket"); }
		bucket.Tokens += tokens;
	}
	// merge server uids
	for _, uid := range uids {
		db.UIDs[uid] = true; }
}
