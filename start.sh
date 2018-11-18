#!/bin/bash
go build -o ./main src/server/main.go
nohup ./main > server.log 2>&1 &
