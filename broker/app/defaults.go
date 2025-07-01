package app;
// pxnMetrics Broker - defaults



const LogPrefix = "[Broker] ";
const DefaultConfigFile = "config.yml"



const DefaultNumShards      = 1;
const DefaultChecksumBase   = 0;
const DefaultListenInterval = "500ms";
const DefaultSyncInterval   = "5s";
const DefaultBatchInterval  = "30s";
const DefaultBindRPC        = "tcp://127.0.0.1:9901";
//const DefaultBindRPC = "unix:///run/pxnMetrics/broker.sock";

// rate limiter
const DefaultTokenInterval = "10s";
const DefaultTokensPerHit  = 3;
const DefaultTokensThresh  = 35;
const DefaultTokensCap     = 50;
