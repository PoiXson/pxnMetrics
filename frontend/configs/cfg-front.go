package configs;
// pxnMetrics Frontend - config

import(
	_ "gopkg.in/yaml.v2"
);



type CfgFront struct {
	BrokerAddr string `yaml:"Broker-Addr"`
	BindWeb    string `yaml:"Bind-Web"`
	Proxied    bool   `yaml:"Proxied"`
	// loaded from broker rpc
	NumShards  uint8
}
