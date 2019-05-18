#!/bin/bash

cd ~/server

echo "killing server process"
sudo pkill server

echo "building new server binary"
go build

echo "getting rid of old log"
/bin/rm -f nohup.out

ip_addr="$(ifconfig eth0 | grep "inet " | awk '{print $2}')"

echo "starting server"
nohup sudo ./server $ip_addr &
