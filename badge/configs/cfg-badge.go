package configs;
// badge.minecraftmetrics.com config

import(
	_ "gopkg.in/yaml.v2"
);



type CfgBadge struct {
	BindWeb    string `yaml:"Bind-Web"`
	BrokerAddr string `yaml:"Broker-Addr"`
	Proxied    bool   `yaml:"Proxied"`
}
