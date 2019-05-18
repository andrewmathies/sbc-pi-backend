#!/bin/bash

wget https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz
sudo tar -xvf go1.12.5.linux-amd64.tar.gz
sudo mv go /usr/local
sudo rm go1.12.5.linux-amd64.tar.gz
sudo rm -rf go

export GOROOT=/usr/local/go
export GOPATH=/home/ubuntu/server
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
