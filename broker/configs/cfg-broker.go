package configs;
// pxnMetrics Broker - config

import(
	_ "gopkg.in/yaml.v2"
);



type CfgBroker struct {
	BindRPC         string   `yaml:"Bind-RPC"`
	ChecksumBase    uint16   `yaml:"Checksum-Base"`
	ListenInterval  string   `yaml:"Listen-Interval"`
	SyncInterval    string   `yaml:"Sync-Interval"`
	BatchInterval   string   `yaml:"Batch-Interval"`
	RateLimit CfgRateLimit   `yaml:"Rate-Limit"`
	Users map[string]CfgUser `yaml:"Users"`
	NumShards       uint8
}

type CfgRateLimit struct {
	TokenInterval string `yaml:"Token-Interval"`
	TokensPerHit  uint16 `yaml:"Tokens-Per-Hit"`
	TokensThresh  uint16 `yaml:"Tokens-Thresh"`
	TokensCap     uint16 `yaml:"Tokens-Cap"`
}

type CfgUser struct {
	Desc         string   `yaml:"Desc"`
	PermitIPs    []string `yaml:"Permit-IPs"`
	PermitWeb    bool     `yaml:"Permit-Web"`
	PermitShards []uint8  `yaml:"Permit-Shards"`
}
