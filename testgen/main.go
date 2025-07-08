package main;

import(
//	OS      "os"
	Fmt     "fmt"
	Net     "net"
	Time    "time"
	Rand    "math/rand"
	Atomic  "sync/atomic"
	Adler   "hash/adler32"
	StrConv "strconv"
	JSON    "encoding/json"
	PxnUID  "github.com/PoiXson/pxnGoCommon/utils/uid"
	APIv1   "github.com/PoiXson/pxnMetrics/api/submit_v1"
);



const Address = "127.0.0.1:9001";
const ThreadCount = 1;
const ChecksumBase = 9000;
const RandSleep = 600;



func main() {
	print("\nStarting..\n\n");
	tool := NewTool();
	for i:=0; i<ThreadCount; i++ {
		go tool.Run();
	}
	tool.Ticker();
//	slp, _ := Time.ParseDuration("5.1s");
//	for { Time.Sleep(slp); }
//	print("<end>\n\n");
//	OS.Exit(0);
}



type Tool struct {
	Count Atomic.Uint64
}

func NewTool() *Tool {
	return &Tool{};
}



func (tool *Tool) Ticker() {
	tick_interval, _ := Time.ParseDuration("1s");
	ticker := Time.NewTicker(tick_interval);
	defer ticker.Stop();
	var last uint64 = 0;
	for { select {
		case <-ticker.C: {
			cnt := tool.Count.Load();
			Fmt.Printf(" %s per sec\n", Format(int64(cnt-last)));
			last = cnt;
		}
	}}
}

func (tool *Tool) Run() {
	uid := PxnUID.New(0);
	addr, err := Net.ResolveUDPAddr("udp", Address)
	if err != nil { panic(err); }
	conn, err := Net.DialUDP("udp", nil, addr);
	if err != nil { panic(err); }
	for {
index, _ := uid.Next();
//Fmt.Printf("UID: %s\n", index.ToString());
sleep := Time.Duration(Rand.Intn(RandSleep)) * Time.Millisecond; Time.Sleep(sleep);
		tool.Count.Add(1);
		// build submit packet
		timestamp := Time.Now().UnixMilli();
		num_players := int16(Rand.Intn(123));
		var platform string;
		switch Rand.Intn(7) {
			case 0, 1, 2, 3: platform = "PaperMC"; break;
			case 4:          platform = "Folia";   break;
			case 5, 6, 7:    platform = "Fabric";  break;
		}
		// packet
		json, err := JSON.Marshal(APIv1.Submit{
			Timestamp:  timestamp,
			ServerUID:  index.ToString(),
			Platform:   platform,
			NumPlayers: num_players,
		});
		if err != nil { panic(err); }
		data := []byte(json);
		size := len(data);
		hash32 := Adler.Checksum(data);
		hash16 := uint16(((hash32 >> 16) & 0xFFFF) ^ (hash32 & 0xFFFF)) ^ ChecksumBase;
		// send
		out := make([]byte, size+7);
		out[0] = 0x07; // header size
		out[1] = 0x00; out[2] = byte(size); // size
		out[3] = byte((hash16 >> 8) & 0xFF); // checksum
		out[4] = byte( hash16       & 0xFF); // checksum
		out[5] = 0x00; // encryption
		out[6] = 0x01; // protocol
		copy(out[7:], data);
		_, err = conn.Write(out);
		if err != nil { panic(err); }
	}
}



func Format(n int64) string {
	in := StrConv.FormatInt(n, 10);
	numOfDigits := len(in);
	// First character is the - sign (not a digit)
	if n < 0 { numOfDigits--; }
	numOfCommas := (numOfDigits - 1) / 3;
	out := make([]byte, len(in)+numOfCommas);
	if n < 0 { in, out[0] = in[1:], '-'; }
	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 { return string(out); }
		if k++; k == 3 {
			j, k = j-1, 0;
			out[j] = ',';
		}
	}
}
