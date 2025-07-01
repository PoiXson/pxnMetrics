package configs;
// pxnMetrics Shard - config

import(
	_ "gopkg.in/yaml.v2"
);



type CfgShard struct {
	BrokerAddr     string `yaml:"Broker-Addr"`
	BindPublic     string `yaml:"Bind-Public"`
	// loaded from broker rpc
	NumShards      uint8
	ShardIndex     uint8
	ChecksumBase   uint16
	ListenInterval string
	SyncInterval   string
	RateLimit      CfgRateLimit
}

type CfgRateLimit struct {
	TokenInterval string
	TokensPerHit  uint16
	TokensThresh  uint16
	TokensCap     uint16
}
