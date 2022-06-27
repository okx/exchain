#!/bin/sh
#set -e
#set -x
VERSION_NUM=5.2.1
VERSION=jemalloc-$VERSION_NUM

command_exists() {
	command -v "$@" > /dev/null 2>&1
}

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
	$sh_c "rm -rf jemalloc"
	$sh_c "git clone https://github.com/jemalloc/jemalloc.git"
	$sh_c "cd jemalloc && git checkout ${VERSION_NUM}"
	$sh_c "cd jemalloc && ./autogen.sh"
	$sh_c "cd jemalloc && ./configure --prefix=/usr --libdir=/usr/lib"
	$sh_c "cd jemalloc && make uninstall"
	$sh_c "cd jemalloc && make install"
	$sh_c "ldconfig"
	$sh_c "rm -rf jemalloc"
}

install_macos(){
	JEMALLOC=jemalloc-$VERSION_NUM
	$sh_c "wget -c https://github.com/jemalloc/jemalloc/releases/download/$VERSION_NUM/$JEMALLOC.tar.bz2"
	$sh_c "tar -xvf $JEMALLOC.tar.bz2"
	$sh_c "cd $JEMALLOC&& ./configure --disable-cpu-profiler --disable-heap-profiler --disable-heap-checker --disable-debugalloc --enable-minimal"
	$sh_c "cd $JEMALLOC&& make uninstall"
	$sh_c "cd $JEMALLOC&& make install"
	$sh_c "rm $JEMALLOC.tar.bz2"
	$sh_c "rm -r $JEMALLOC"
}

do_install() {
	echo "# Executing jemalloc install script, version: $VERSION"

	user="$(id -un 2>/dev/null || true)"

	sh_c='sh -c'
	if [ "$user" != 'root' ]; then
		if command_exists sudo; then
			sh_c='sudo -E sh -c'
		elif command_exists su; then
			sh_c='su -c'
		fi
	fi

	# perform some very rudimentary platform detection
	lsb_dist=$( get_distribution )
	lsb_dist="$(echo "$lsb_dist" | tr '[:upper:]' '[:lower:]')"

	# Run setup for each distro accordingly
	case "$lsb_dist" in
		ubuntu)
			pre_reqs="git make autoconf automake libtool gcc-c++"
			$sh_c 'apt-get update -qq >/dev/null'
			$sh_c "apt-get install -y -qq $pre_reqs >/dev/null"
			install_linux
			exit 0
			;;
		centos)
			pre_reqs="git make autoconf automake libtool gcc-c++"
			$sh_c "yum install -y -q $pre_reqs"
			install_linux
			exit 0
			;;
		*)
			if [ -z "$lsb_dist" ]; then
				if is_darwin; then
					pre_reqs="wget"
					$sh_c "xcode-select --install"
					$sh_c "brew install $pre_reqs"
					install_macos
					exit 0
				fi
				if is_wsl; then
					echo
					echo "ERROR: Unsupported OS 'Windows'"
					echo "Please install jemalloc from https://github.com/jemalloc/jemalloc"
					echo
					exit 1
				fi
			fi
			echo
			echo "ERROR: Unsupported distribution '$lsb_dist'"
			echo "Please install jemalloc from https://github.com/jemalloc/jemalloc"
			echo
			exit 1
			;;
	esac
	exit 1
}

# wrapped up in a function so that we have some protection against only getting
# half the file during "curl | sh"
do_install
