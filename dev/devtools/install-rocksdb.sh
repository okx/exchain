#!/bin/bash

set -euo pipefail

if [ ! -d rocksdb ]; then
  git clone https://github.com/facebook/rocksdb.git
fi

cd rocksdb
git checkout v6.26.1
make shared_lib
make install-shared

if [[ `uname` == 'Linux' ]]; then
    cp librocksdb.so* /usr/lib/
fi

rm -rf ../rocksdb
