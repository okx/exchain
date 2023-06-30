#!/bin/sh
#set -e

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

get_system_version() {
  lsb_dist=$(get_distribution)
  lsb_dist="$(echo "$lsb_dist" | tr '[:upper:]' '[:lower:]')"
  system_version=
  case "$lsb_dist" in
  ubuntu)
    system_version="ubuntu"
    ;;
  centos)
    system_version="centos"
    ;;
  alpine)
    system_version="alpine"
    ;;
  *)
    if [ -z "$lsb_dist" ]; then
      if is_darwin; then
        system_version="macos"
      fi
    else
      echo
      echo "ERROR: Unsupported system '$lsb_dist', only support centos,ubuntu,alpine,macos"
      echo
      exit 1
    fi
    ;;
  esac

  # checkout go version

  echo "$system_version"
}

echo $(get_system_version)
