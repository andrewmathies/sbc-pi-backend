#!/bin/bash

cd ~/server

echo "killing server process"
sudo pkill server

echo "building new server binary"
go build

echo "getting rid of old log"
/bin/rm -f nohup.out

echo "starting server"
sudo nohup sudo ./server &
