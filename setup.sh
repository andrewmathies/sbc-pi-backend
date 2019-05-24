#!/bin/bash

wget https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz
sudo tar -xvf go1.12.5.linux-amd64.tar.gz
sudo mv go /usr/bin

export GOROOT=/usr/bin/go
export GOPATH=/
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH


sudo rm go1.12.5.linux-amd64.tar.gz
sudo rm -rf go

