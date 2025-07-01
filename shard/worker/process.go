package worker;
// pxnMetrics Shard - packet processor

import(
	Log       "log"
	Fmt       "fmt"
	Net       "net"
	Time      "time"
	Bytes     "bytes"
	Binary    "encoding/binary"
	Adler     "hash/adler32"
	Context   "context"
	Errors    "errors"
	GRPC      "google.golang.org/grpc"
	API_Shard "github.com/PoiXson/pxnMetrics/api/shard"
);



const TimeoutSync  = "10s";
const TimeoutBatch = "60s";



func (worker *Worker) Process(data []byte, src *Net.Addr) ([]byte, error) {
	reader := Bytes.NewReader(data);
	var first       byte;
	var size        uint16;
	var checksum    uint16;
	var index_crypt byte;
	var index_proto byte;
	var payload     []byte;
	// check first byte
	if err := Binary.Read(reader, Binary.BigEndian, &first); err != nil {
		return nil, err; }
	// header version
	switch first {
	case 0x07: // 7 byte header v1
		if err := Binary.Read(reader, Binary.BigEndian, &size       );
			err != nil { return nil, err; }
		if err := Binary.Read(reader, Binary.BigEndian, &checksum   );
			err != nil { return nil, err; }
		if err := Binary.Read(reader, Binary.BigEndian, &index_crypt);
			err != nil { return nil, err; }
		if err := Binary.Read(reader, Binary.BigEndian, &index_proto);
			err != nil { return nil, err; }
		var buffer Bytes.Buffer;
		if _, err := buffer.ReadFrom(reader); err != nil { return nil, err; }
		payload = buffer.Bytes();
	// invalid packet (first byte)
	default: return nil, Fmt.Errorf("Invalid first byte: %X (%d)", first, first);
	}
	// packet size
	if len(payload) != int(size) {
		return nil, Errors.New("Invalid packet length"      ); }
	if size <    10 { return nil, Errors.New("Invalid short packet length"); }
	if size > 10000 { return nil, Errors.New("Invalid long packet length" ); }
	// checksum
	hash32 := Adler.Checksum(payload);
	hash16 := uint16(((hash32 >> 16) & 0xFFFF) ^
		(hash32 & 0xFFFF)) ^ worker.config.ChecksumBase;
	if checksum != hash16 { return nil, Errors.New("Invalid packet checksum"); }
	// encryption
	switch index_crypt {
		// plain text
//TODO: disable this by default
		case 0x00: break;
		// AES-GCM 128
		case 0x01: panic(Errors.New("UNFINISHED"));
		// AES-GCM 192
		case 0x02: panic(Errors.New("UNFINISHED"));
		// AES-GCM 256
		case 0x03: panic(Errors.New("UNFINISHED"));
		// invalid encryption
		default: return nil, Fmt.Errorf("Invalid encryption value: %X (%d)",
			index_crypt, index_crypt);
	}
	// protocol version
	switch index_proto {
	// Submit V1
	case 0x01: {
		reply, err := worker.ProcessV1(worker.secret_db, payload);
		if err != nil { return nil, err; }
		// send reply
		hash32 := Adler.Checksum(reply);
		hash16 := uint16(((hash32 >> 16) & 0xFFFF) ^
			(hash32 & 0xFFFF)) ^ worker.config.ChecksumBase;
		size := len(reply);
		out := make([]byte, size+7);
		out[0] = 0x07;                          // header size
		out[1] = byte((size   >> 16) & 0xFFFF); // payload size high
		out[2] = byte( size          & 0xFFFF); // payload size low
		out[3] = byte((hash16 >> 16) & 0xFFFF); // checksum high
		out[4] = byte( hash16        & 0xFFFF); // checksum low
		out[5] = index_crypt;                   // encryption
		out[6] = index_proto;                   // protocol version
		copy(out[7:], reply);
		return out, nil;
	}
	// invalid version
	default: break;
	}
	return nil, Fmt.Errorf("Invalid protocol version: %X (%d)", index_proto, index_proto);
}



func (worker *Worker) DoSync(islast bool, shard_index uint8) {
	worker.service.WaitGroup.Add(1);
	defer worker.service.WaitGroup.Done();
	if islast {
		worker.mut_update.Lock();
	} else
	if !worker.mut_update.TryLock() {
		Log.Printf("[Shard-%d] Still syncing?!", shard_index);
		return;
	}
	defer worker.mut_update.Unlock();
	buckets, uids := worker.secret_db.UpdatePush();
	if islast { worker.dosync(shard_index, buckets, uids);
	} else { go worker.dosync(shard_index, buckets, uids); }
}

func (worker *Worker) DoBatch(islast bool, shard_index uint8, chip *Chip) {
	worker.service.WaitGroup.Add(1);
	defer worker.service.WaitGroup.Done();
	if islast {
		worker.mut_update.Lock();
	} else
	if !worker.mut_update.TryLock() {
		Log.Printf("[Shard-%d] Still batching?!", shard_index);
		return;
	}
	defer worker.mut_update.Unlock();
	if islast { worker.dobatch(islast, shard_index, chip);
	} else { go worker.dobatch(islast, shard_index, chip); }
}

func (worker *Worker) dosync(shard_index uint8, buckets map[string]int32, uids []uint64) {
//TODO: remove
print("DO SYNC!  ");
Fmt.Printf("Buckets: %d  uids: %d\n", len(buckets), len(uids));
for k, b := range buckets { Fmt.Printf("   %s   %d\n", k, b); }
	sync_send := API_Shard.SyncSend{
		ShardIndex:   uint32(shard_index),
		TokenBuckets: buckets,
		ServerUIDs:   uids,
	};
	timeout, _ := Time.ParseDuration(TimeoutSync);
	ctx, cancel := Context.WithTimeout(Context.Background(), timeout);
	defer cancel();
	sync_reply, err := worker.link.API.SyncDBs(
		ctx,
		&sync_send,
		GRPC.WaitForReady(true),
	);
	if err != nil {
		Log.Printf("SyncDBs error: %v", err);
		return;
	}
	worker.secret_db.UpdatePull(
		sync_reply.TokenBuckets,
		sync_reply.ServerUIDs,
	);
	if sync_reply.NeedsBatch { worker.needs_batch.Store(true); }
}

func (worker *Worker) dobatch(islast bool, shard_index uint8, chip *Chip) {
//TODO: remove
print("DO BATCH OUT!\n");
	batch_data := API_Shard.BatchData{
		ShardIndex: uint32(shard_index),
		IsLast:     islast,
//TODO: chip data
	};
	timeout, _ := Time.ParseDuration(TimeoutBatch);
	ctx, cancel := Context.WithTimeout(Context.Background(), timeout);
	defer cancel();
	if _, err := worker.link.API.BatchOut(
		ctx,
		&batch_data,
		GRPC.WaitForReady(true),
	); err != nil {
		if islast { print("\n"); Log.Printf("Final batchout failed!!!"); }
		Log.Printf("BatchOut error: %v", err);
		return;
	}
}
