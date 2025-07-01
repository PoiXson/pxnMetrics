clear
go mod tidy                             || exit $?
go run . --broker tcp://192.168.3.3:9901  || exit $?
