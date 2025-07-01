clear
go mod tidy  || exit $?
go run . --bind udp://127.0.0.1:9001 --broker tcp://192.168.3.3:9901 --shard-index=2     || exit $?
