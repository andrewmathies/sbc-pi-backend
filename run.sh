#!/bin/bash

cd sbc-pi-backend

echo "killing server process"
pkill server

echo "building new server binary"
go build

echo "getting rid of old log"
/bin/rm -f nohup.out

ip_addr="$(ifconfig eth0 | grep "inet " | awk '{print $2}')"

echo "starting server"
nohup ./sbc-pi-backend $ip_addr &
