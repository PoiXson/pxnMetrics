package worker;

import(
	SecretDB "github.com/PoiXson/pxnMetrics/shard/worker/secretdb"
);



func (worker *Worker) ProcessV1(db_secret *SecretDB.DB,
		payload []byte) ([]byte, error) {
return []byte("{}"), nil;
}
