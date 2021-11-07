#!/bin/sh

apt-get update
apt-get -y upgrade
apt-get -y install make screen

GOFILE=go1.12.7.linux-amd64.tar.gz

wget https://storage.googleapis.com/golang/${GOFILE}
tar -C /usr/local -xzf ${GOFILE}
