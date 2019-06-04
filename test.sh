#!/bin/bash

sudo pkill sbc-pi-backend
rm ./nohup.out
go build
nohup sudo ./sbc-pi-backend &
