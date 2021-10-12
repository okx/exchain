#!/bin/bash

set -euo pipefail

if [ ! -d rocksdb ]; then
  git clone https://github.com/facebook/rocksdb.git
fi

cd rocksdb
git checkout v6.15.5
make shared_lib
make install-shared
rm -r ../rocksdb
