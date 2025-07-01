package web;
// pxnMetrics web api


type Status struct {
	Broker BrokerStatus   `json:"Broker"`
	Shards []ServerStatus `json:"Shards"`
}

type BrokerStatus struct {
	Status ServerStatus `json:"Status"`
	LastTask uint32     `json:"LastTask"`
}

type ServerStatus struct {
	Name         string  `json:"Name"`
	Desc         string  `json:"Desc"`
	Status       string  `json:"Status"`
	LastSeen     uint32  `json:"LastSeen"`
	LastSync     uint32  `json:"LastSync"`
	LastBatch    uint32  `json:"LastBatch"`
//	LastTask     uint32  `json:"LastTask"`
	TasksQueued  uint32  `json:"TasksQueued"`
//TODO
//	QueueWaiting uint32  `json:"QueueWaiting"`
//	ReqPerSec    float32 `json:"ReqPerSec"`
//	ReqPerMin    float32 `json:"ReqPerMin"`
//	ReqPerHour   float32 `json:"ReqPerHour"`
//	ReqPerDay    float32 `json:"ReqPerDay"`
//	ReqTotal     uint64  `json:"ReqTotal"`
}
