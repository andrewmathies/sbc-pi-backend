#!/bin/bash

echo "Killing server"
sudo pkill sbc-pi-backend

echo "Removing old log file, building server"
rm ./nohup.out
go build

echo "Restarting server"
nohup sudo ./sbc-pi-backend &