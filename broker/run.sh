
while true; do
	clear
	go mod tidy  || break
	go run .     || break
	sleep 4.5 ; echo ; sleep 1.5
done
