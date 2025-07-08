{{define "page-content"}}
pxnMetrics Wiki



== How It Works ==

1. The Minecraft server connects to https://api.minecraftmetrics.com/
  * If a server UID is needed, it is generated and provided here.
  * A list of available backend servers is also provided here.
2. The minecraft server picks a backend to submit metrics.
  * The data is sent as an encrypted UDP packet to one of the backend servers.





=== Shard Batching ===

1. Each shard synchronizes with the broker at a set interval.
  * Encoded IP addresses and server UID's are sent to the broker.
  * Reply also contains encoded IP addresses and server UID's.


1. Each shard batches out at a set interval, in sequence.
{{end}}
