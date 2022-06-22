#!/bin/sh
#set -e
GO_VERSION=$1
ROCKSDB_VERSION=$2

get_distribution() {
  lsb_dist=""
  # Every system that we officially support has /etc/os-release
  if [ -r /etc/os-release ]; then
    lsb_dist="$(. /etc/os-release && echo "$ID")"
  fi
  # Returning an empty string here should be alright since the
  # case statements don't act unless you provide an actual value
  echo "$lsb_dist"
}

is_darwin() {
  case "$(uname -s)" in
  *darwin*) true ;;
  *Darwin*) true ;;
  *) false ;;
  esac
}

version_lt() { test "$(echo "$@" | tr " " "\n" | sort -rV | head -n 1)" != "$1"; }


check_go_verison() {
  # check go,awk is install
  hasgo=$(which go)
  if [ -z "$hasgo" ]; then
    echo "ERROR: command go is not found,please install go${GO_VERSION}"
    exit 1
  fi

  # checkout go version
  go_version=$(go version | awk '{print$3}' | awk '{ gsub(/go/,""); print $0 }')
  if version_lt $go_version $GO_VERSION; then
     echo "ERROR: exchain requires go${GO_VERSION}+,please install"
     exit 1
  fi

  echo "go check success: "$go_version
}

check_rocksdb_version() {

  lsb_dist=$(get_distribution)
  lsb_dist="$(echo "$lsb_dist" | tr '[:upper:]' '[:lower:]')"
  rocksdb_version=
  case "$lsb_dist" in
  ubuntu)
    rocksdb_version=$(cat /usr/lib/pkgconfig/rocksdb.pc | grep Version: | awk '{print $2}')
    ;;
  centos)
    rocksdb_version=$(cat /usr/lib/pkgconfig/rocksdb.pc | grep Version: | awk '{print $2}')
    ;;
  alpine)
    rocksdb_version=$(cat /usr/local/lib/pkgconfig/rocksdb.pc | grep Version: | awk '{print $2}')
    ;;
  *)
    if [ -z "$lsb_dist" ]; then
      if is_darwin; then
        rocksdb_version=$(cat /usr/local/lib/pkgconfig/rocksdb.pc | grep Version: | awk '{print $2}')
      fi
    else
      echo
      echo "ERROR: Unsupported distribution '$lsb_dist'"
      echo
      exit 1
    fi
    ;;
  esac

  # checkout go version

  if [ "$rocksdb_version" != "$ROCKSDB_VERSION" ]; then
    echo "ERROR: exchain requires rocksdb-v${ROCKSDB_VERSION},current: v$rocksdb_version , please install with command (make rocksdb)"
    exit 1
  fi
  echo "rocksdb check success: "$rocksdb_version

}

echo "check go and rocksdb version: "

hasawk=$(which awk)

if [ -z "$hasawk" ]; then
  echo "please install awk"
  exit 1
fi

if [ "$GO_VERSION" != "0" ]; then
  check_go_verison
else
  echo "go version check: ignore "
fi

if [ "$ROCKSDB_VERSION" != "0" ]; then
  check_rocksdb_version
else
  echo "rocksdb version check: ignore "
fi

echo "------------------------------------------------------------------------"
exit 0
