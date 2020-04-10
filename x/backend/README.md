go get go-sqlite3

## Unexpected

1. github.com/okex/okchain/vendor/github.com/mattn/go-sqlite3
../../vendor/github.com/mattn/go-sqlite3/backup.go:14:10: fatal error: 'stdlib.h' file not found
#include <stdlib.h>

solution: https://github.com/mattn/go-sqlite3/issues/481

* try in ubuntu: sudo apt-get install g++
* try in mac(@linsheng.yu):  cd /Library/Developer/CommandLineTools/Packages/; open macOS_SDK_headers_for_macOS_10.14.pkg