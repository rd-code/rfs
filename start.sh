#!/bin/bash
rm -f main server.log
go build -o ./main src/server/main.go
killall main
nohup ./main > server.log 2>&1 &