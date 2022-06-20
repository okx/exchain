#!/bin/sh
set -e
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
	*darwin* ) true ;;
	*Darwin* ) true ;;
	* ) false;;
	esac
}

is_linux() {
  case "$(uname -s)" in
  	*darwin* ) true ;;
  	*Darwin* ) true ;;
  	* ) false;;
  	esac
}

function check_go_verison() {
    # check go,awk is install
    hasgo=$(which go)
    if [ -z "$hasgo" ] ;then
    echo "command go is not found,please install go${GO_VERSION}"
    exit 1
    fi

    # checkout go version
    go_version=$(go version | awk '{print$3}' | awk '{ gsub(/go/,""); print $0 }' )
    find=$(echo $go_version $GO_VERSION| awk '{print index($1,$2)}')
    if [ "$find" != "1" ] ;then
      echo "exchain need go${GO_VERSION},please install"
      exit 1
    fi
}

function check_rocksdb_version() {
    prefix="librocksdb."
    suffix=".dylib"
    file_path="/usr/local/lib/"

    lsb_dist=$( get_distribution )
  	lsb_dist="$(echo "$lsb_dist" | tr '[:upper:]' '[:lower:]')"

  	case "$lsb_dist" in
    		ubuntu)
    		  file_path="/usr/lib"
    			suffix=".so"
    			;;
    		centos)
    		  file_path="/usr/lib"
    		  suffix=".so"
    			;;
    	  alpine)
    	    file_path="/usr/lib"
    	    suffix=".a"
    	    ;;
    		*)
    			if [ -z "$lsb_dist" ]; then
    				if is_darwin; then
    				  file_path="/usr/local/lib/"
              suffix=".dylib"
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
    rocksdb_version=$(ls -al ${file_path}librocksdb${suffix} | awk '{print$11}' | awk '{ gsub(/'$prefix'/,""); gsub(/'$suffix'/,"");print $0 }')

    if [ "$rocksdb_version" != "$ROCKSDB_VERSION" ] ;then
      echo "exchain need rocksdb-v${ROCKSDB_VERSION},please install"
      exit 1
    fi
    echo "check version success:"
    echo "      go check: go-"$GO_VERSION
    echo "      rocksdb check: "${file_path}librocksdb${suffix}"-"$rocksdb_version
    echo "------------------------------------------------------------------------"
}

hasawk=$(which awk)

if [ -z "$hasawk" ] ;then
    echo "please install awk"
    exit 1
fi

if [ "$GO_VERSION" != "0" ] ;then
    check_go_verison
else
    echo "ignore go version check"
fi


check_rocksdb_version

exit 0



