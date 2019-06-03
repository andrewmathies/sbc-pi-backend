#!/bin/bash

sudo pkill sbc-pi-backend
rm ./nohup.out
nohup sudo ./sbc-pi-backend &
