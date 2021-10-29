#!/bin/sh

# This runs benchmarks, by default from master branch of
# github.com/tendermint/iavl
# You can customize this by optional command line args
#
# INSTALL_USER.sh [branch] [repouser]
#
# set repouser as your username to time your fork

BRANCH=${1:-master}
REPOUSER=${2:-tendermint}

export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$HOME/go/bin
export GOPATH=$HOME/go

go get -u github.com/${REPOUSER}/iavl
cd ~/go/src/github.com/${REPOUSER}/iavl
git checkout ${BRANCH}

make bench > results.txt

