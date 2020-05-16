#!/bin/bash

wget https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz
tar -xvf go1.14.3.linux-amd64.tar.gz
sudo mv go /usr/local/

export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

go get

go build -o bareos_exporter


