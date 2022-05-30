#!/bin/sh
#set -e
set -x

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

install_linux(){
  install_unix
	$sh_c "cp ./ed25519_okc/target/release/libed25519_okc.dylib /usr/local/lib"
	$sh_c "rm -rf ed25519_okc"
}

install_macos(){
  install_unix
	$sh_c "cp ./ed25519_okc/target/release/libed25519_okc.so /usr/lib"
	$sh_c "rm -rf ed25519_okc"
}

install_unix() {
	$sh_c "rm -rf ed25519_okc"
	$sh_c "git clone https://github.com/giskook/ed25519_okc.git"
	$sh_c "cd ed25519_okc && cargo build --release"
}

install_rust(){
  rustc_version="$(rustc --version)"
  if [ $(rustc_version) != "rustc*" ]; then
    echo "installing rust..."
    $sh_c "curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
  else
    echo "rust already installed."
  fi
}

do_install() {
	echo "# Executing ed25519_okc install script"

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
			pre_reqs="git"
			$sh_c 'apt-get update -qq >/dev/null'
			$sh_c "apt-get install -y -qq $pre_reqs >/dev/null"
			install_rust
			install_linux
			exit 0
			;;
		centos)
			pre_reqs="git make autoconf automake libtool gcc-c++"
			$sh_c "yum install -y -q $pre_reqs"
			install_rust
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
					echo "Please install ed25519_okc from https://www.rust-lang.org/tools/install"
					echo
					exit 1
				fi
			fi
			echo
			echo "ERROR: Unsupported distribution '$lsb_dist'"
			echo "Please install ed25519_okc from https://www.rust-lang.org/tools/install"
			echo
			exit 1
			;;
	esac
	exit 1
}

# wrapped up in a function so that we have some protection against only getting
# half the file during "curl | sh"
do_install
