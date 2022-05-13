#!/bin/sh
#set -e
#set -x
VERSION_NUM=2.9.1
VERSION=gperftools-$VERSION_NUM
while [ $# -gt 0 ]; do
	case "$1" in
		--version)
			VERSION="$2"
			shift
			;;
		--*)
			echo "Illegal option $1"
			;;
	esac
	shift $(( $# > 0 ? 1 : 0 ))
done

echo $VERSION

command_exists() {
	command -v "$@" > /dev/null 2>&1
}

exit

is_wsl() {
	case "$(uname -r)" in
	*microsoft* ) true ;; # WSL 2
	*Microsoft* ) true ;; # WSL 1
	* ) false;;
	esac
}

is_darwin() {
	case "$(uname -s)" in
	*darwin* ) true ;;
	*Darwin* ) true ;;
	* ) false;;
	esac
}

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

install_linux() {
  $sh_c "git clone https://github.com/gperftools/gperftools.git"
  $sh_c "cd gperftools && git checkout ${VERSION}"
  $sh_c "cd gperftools && sh autogen.sh"
  $sh_c "cd gperftools && ./configure --disable-cpu-profiler --disable-heap-profiler --disable-heap-checker --disable-debugalloc --enable-minimal"
  $sh_c "cd gperftools && make install"
  $sh_c "ldconfig"
}

install_macos(){
}

do_install() {
	echo "# Executing tcmalloc install script, version: $VERSION"

	user="$(id -un 2>/dev/null || true)"

	sh_c='sh -c'
	if [ "$user" != 'root' ]; then
		if command_exists sudo; then
			sh_c='sudo -E sh -c'
		elif command_exists su; then
			sh_c='su -c'
		else
			cat >&2 <<-'EOF'
			Error: this installer needs the ability to run commands as root.
			We are unable to find either "sudo" or "su" available to make this happen.
			EOF
			exit 1
		fi
	fi

	# perform some very rudimentary platform detection
	lsb_dist=$( get_distribution )
	lsb_dist="$(echo "$lsb_dist" | tr '[:upper:]' '[:lower:]')"

	# Run setup for each distro accordingly
	case "$lsb_dist" in
		ubuntu)
			pre_reqs="git make autoconf automake libtool libunwind"
      $sh_c 'apt-get update -qq >/dev/null'
      $sh_c "apt-get install -y -qq $pre_reqs >/dev/null"
      install_linux
			exit 0
			;;
		centos)
		  pre_reqs="git make autoconf automake libtool libunwind"
      $sh_c "yum install -y -q $pre_reqs"
      install_linux
			exit 0
			;;
		*)
			if [ -z "$lsb_dist" ]; then
				if is_darwin; then
				  pre_reqs="git make"
          $sh_c "xcode-select --install"
          $sh_c "brew install $pre_reqs"
          install_macos
          exit 0
				fi
				if is_wsl; then
          echo
          echo "ERROR: Unsupported OS 'Windows'"
          echo "Please install tcmalloc from https://github.com/gperftools/gperftools/tree/gperftools-2.9.1"
          echo
          exit 1
        fi
			fi
      echo
      echo "ERROR: Unsupported distribution '$lsb_dist'"
      echo
      exit 1
			;;
	esac
	exit 1
}

# wrapped up in a function so that we have some protection against only getting
# half the file during "curl | sh"
do_install
