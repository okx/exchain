#!/bin/bash

printLogo() {
  echo '
     OOOOOOOOO     KKKKKKKKK    KKKKKKKEEEEEEEEEEEEEEEEEEEEEE                                        hhhhhhh                                 iiii
   OO:::::::::OO   K:::::::K    K:::::KE::::::::::::::::::::E                                        h:::::h                                i::::i
 OO:::::::::::::OO K:::::::K    K:::::KE::::::::::::::::::::E                                        h:::::h                                 iiii
O:::::::OOO:::::::OK:::::::K   K::::::KEE::::::EEEEEEEEE::::E                                        h:::::h
O::::::O   O::::::OKK::::::K  K:::::KKK  E:::::E       EEEEEExxxxxxx      xxxxxxx    cccccccccccccccc h::::h hhhhh         aaaaaaaaaaaaa   iiiiiii nnnn  nnnnnnnn
O:::::O     O:::::O  K:::::K K:::::K     E:::::E              x:::::x    x:::::x   cc:::::::::::::::c h::::hh:::::hhh      a::::::::::::a  i:::::i n:::nn::::::::nn
O:::::O     O:::::O  K::::::K:::::K      E::::::EEEEEEEEEE     x:::::x  x:::::x   c:::::::::::::::::c h::::::::::::::hh    aaaaaaaaa:::::a  i::::i n::::::::::::::nn
O:::::O     O:::::O  K:::::::::::K       E:::::::::::::::E      x:::::xx:::::x   c:::::::cccccc:::::c h:::::::hhh::::::h            a::::a  i::::i nn:::::::::::::::n
O:::::O     O:::::O  K:::::::::::K       E:::::::::::::::E       x::::::::::x    c::::::c     ccccccc h::::::h   h::::::h    aaaaaaa:::::a  i::::i   n:::::nnnn:::::n
O:::::O     O:::::O  K::::::K:::::K      E::::::EEEEEEEEEE        x::::::::x     c:::::c              h:::::h     h:::::h  aa::::::::::::a  i::::i   n::::n    n::::n
O:::::O     O:::::O  K:::::K K:::::K     E:::::E                  x::::::::x     c:::::c              h:::::h     h:::::h a::::aaaa::::::a  i::::i   n::::n    n::::n
O::::::O   O::::::OKK::::::K  K:::::KKK  E:::::E       EEEEEE    x::::::::::x    c::::::c     ccccccc h:::::h     h:::::ha::::a    a:::::a  i::::i   n::::n    n::::n
O:::::::OOO:::::::OK:::::::K   K::::::KEE::::::EEEEEEEE:::::E   x:::::xx:::::x   c:::::::cccccc:::::c h:::::h     h:::::ha::::a    a:::::a i::::::i  n::::n    n::::n
 OO:::::::::::::OO K:::::::K    K:::::KE::::::::::::::::::::E  x:::::x  x:::::x   c:::::::::::::::::c h:::::h     h:::::ha:::::aaaa::::::a i::::::i  n::::n    n::::n
   OO:::::::::OO   K:::::::K    K:::::KE::::::::::::::::::::E x:::::x    x:::::x   cc:::::::::::::::c h:::::h     h:::::h a::::::::::aa:::ai::::::i  n::::n    n::::n
     OOOOOOOOO     KKKKKKKKK    KKKKKKKEEEEEEEEEEEEEEEEEEEEEExxxxxxx      xxxxxxx    cccccccccccccccc hhhhhhh     hhhhhhh  aaaaaaaaaa  aaaaiiiiiiii  nnnnnn    nnnnnn
  '
}

getOSName() {
    if [ ! -f "/etc/os-release" ]; then
      if [ -f "/etc/centos-release" ]; then
        os=$(cat /etc/centos-release 2>/dev/null |awk -F' ' '{print $1}')
      fi
    else
      # shellcheck disable=SC2046
      os=$(cat /etc/os-release 2>/dev/null | grep ^ID= | awk -F= '{print $2}')

      if [ "$os" = "" ]; then
          # shellcheck disable=SC2046
          os=$(trim $(lsb_release -i 2>/dev/null | awk -F: '{print $2}'))
      fi
      if [ ! "$os" = "" ]; then
          # shellcheck disable=SC2021
          os=$(echo "$os" | tr '[A-Z]' '[a-z]')
      fi
    fi

    echo "$os"
}

GetArchitecture() {
  _cputype="$(uname -m)"
  _ostype="$(uname -s)"
  if [[ "$_ostype" == "Linux" ]]; then
          set -e
          if [[ $(whoami) == "root" ]]; then
                  MAKE_ME_ROOT=
          else
                  MAKE_ME_ROOT=sudo
          fi
          echo "echo Arch Linux detected."
          os="linux"
          deptool="ldd"
          goAchive="go1.17.linux-amd64.tar.gz"
          libArray=("/usr/local/lib/librocksdb.so.6.27.3" "/usr/local/lib/librocksdb.so.6.27" "/usr/local/lib/librocksdb.so.6" "/usr/local/lib/librocksdb.so" "/usr/lib/librocksdb.so.6.27.3" "/usr/lib/librocksdb.so.6.27" "/usr/lib/librocksdb.so.6" "/usr/lib/librocksdb.so")
          rocksdbdep=("/usr/local/lib/pkconfig/rocksdb.pc" "/usr/local/include/rocksdb" "/usr/local/lib/librocksdb.so*" "/usr//lib/pkconfig/rocksdb.pc" "/usr/include/rocksdb" "/usr/lib/librocksdb.so*")
          case "$(getOSName)" in
              *ubuntu*)
                  echo "detected ubuntu ..."
                  dynamicLink="TRUE"
                  $MAKE_ME_ROOT apt-get install -y wget g++ cmake make git gnutls-bin clang
                  ;;
              *centos*)
                  echo "detected centos ..."
                  dynamicLink="TRUE"
                  # binutils for link text error
                  $MAKE_ME_ROOT yum install -y wget tar gcc gcc-c++ automake autoconf libtool make which git perl-Digest-SHA glibc-static.x86_64 libstdc++-static clang binutils
                  if ! type cmake > /dev/null 2>&1; then
                    export CXXFLAGS="-stdlib=libstdc++" CC=/usr/bin/gcc CXX=/usr/bin/g++
                    $MAKE_ME_ROOT wget --no-check-certificate https://github.com/Kitware/CMake/releases/download/v3.15.5/cmake-3.15.5.tar.gz
                    tar -zxvf cmake-3.15.5.tar.gz && rm cmake-3.15.5.tar.gz
                    cd cmake-3.15.5
                    ./bootstrap && make -j4 && $MAKE_ME_ROOT make install
                    $MAKE_ME_ROOT cp /cmake-3.15.5/bin/cmake /usr/bin/
                    unset CXXFLAGS CC CXX
                  fi

                  ;;
              *CentOS*)
                  n=`cat /etc/centos-release`
                  echo "detected $n  not supported"
                  exit 1
                  ;;
              *alpine*)
                  echo "detected alpine ..."
                  dynamicLink="FALSE"
                  apk add make git libc-dev bash gcc cmake linux-headers eudev-dev g++ snappy-dev perl
                  cd /bin && rm -f sh && ln -s /bin/bash sh
                  ;;
              *)
                  echo unknow os $OS, exit!
                  exit 1
                  ;;
          esac
  elif [[ "$_ostype" == "Darwin"* ]]; then
          set -e
          echo "Mac OS (Darwin) detected."
          os="darwin"
          deptool="otool -L"
          libArray=("/usr/local/lib/librocksdb.6.27.dylib" "/usr/local/lib/librocksdb.6.dylib" "/usr/local/lib/librocksdb.dylib")
          rocksdbdep=("/usr/local/lib/pkconfig/rocksdb.pc" "/usr/local/include/rocksdb/" "/usr/local/lib/librocksdb.*")
          dynamicLink="TRUE"
          echo "$_cputype"
          if [ "$_cputype" == "arm64" ]
          then
            goAchive="go1.17.darwin-arm64.tar.gz"
          else
            goAchive="go1.17.darwin-amd64.tar.gz"
          fi
          brew install wget
  else
          echo "Unknown operating system.${_ostype}"
          echo "This OS is not supported with this script at present. Sorry."
          exit 1
  fi
}

download() {
  rm -rf "$HOME"/.exchain/src
  mkdir -p "$HOME"/.exchain/src
  tag=`wget -qO- -t1 -T2 --no-check-certificate "https://api.github.com/repos/okx/exchain/releases/latest" | grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g'`
  wget --no-check-certificate "https://github.com/okx/exchain/archive/refs/tags/${tag}.tar.gz" -O "$HOME"/.exchain/src/exchain.tar.gz
  ver=$(echo $tag| sed 's/v//g')
  cd "$HOME"/.exchain/src && tar zxvf exchain.tar.gz &&  cd exchain-"$ver"
}

function checkgoversion { echo "$@" | awk -F. '{ printf("%d%03d%03d%03d\n", $1,$2,$3,$4); }'; }

installRocksdb() {
  echo "install rocksdb...."
  if [ "$dynamicLink" == "TRUE" ]; then
    make rocksdb
  else
    installRocksdbStatic
  fi
  # shellcheck disable=SC2181
  if [ $? -gt 0 ]; then
    echo "install rocksdb error"
    exit 1
  fi
  echo "installRocksdb success"
}

installRocksdbStatic() {
  wget "https://github.com/facebook/rocksdb/archive/refs/tags/v6.27.3.tar.gz" --no-check-certificate -O /tmp/rocksdb.tar.gz && \
  cd /tmp/ && tar zxvf rocksdb.tar.gz && \
  cd rocksdb-6.27.3 && \
  sed -i 's/install -C /install -c /g' Makefile && \
  make libsnappy.a && $MAKE_ME_ROOT cp libsnappy.a /usr/lib && \
  make liblz4.a && $MAKE_ME_ROOT cp liblz4.a /usr/lib && \
  make static_lib PREFIX=/usr LIBDIR=/usr/lib && \
  $MAKE_ME_ROOT make install-static PREFIX=/usr LIBDIR=/usr/lib && \
  rm -rf /tmp/rocksdb
  echo "rocksdb install completed"
}

uninstallRocksdb() {
  # shellcheck disable=SC2068
  for lib in ${rocksdbdep[@]}
  do
    echo "rm lib ${lib}"
    $MAKE_ME_ROOT rm -rf $lib
  done
  echo "uninstallRocksdb ..."
}

installgo() {
  echo "install go ..."
  if [[ -d "/usr/local/go" ]]; then
    rm -rf "/usr/local/go"
  fi
  if [[ -f "${goAchive}" ]]; then
      rm ${goAchive}
  fi
  wget --no-check-certificate "https://golang.google.cn/dl/${goAchive}"
  $MAKE_ME_ROOT tar -zxvf ${goAchive} -C /usr/local/
  rm ${goAchive}
  cd ~
  if [[ -f ".bashrc" ]]; then
      echo "PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
      source ~/.bashrc
  fi
  if [[ -f ".zshrc" ]]; then
      echo "PATH=\$PATH:/usr/local/go/bin" >> ~/.zshrc
      source ~/.zshrc
  fi

  #/usr/local/go/bin/go env -w GOPROXY=https://goproxy.cn,direct
  /usr/local/go/bin/go env -w GOPROXY="https://goproxy.io"
  /usr/local/go/bin/go env -w GO111MODULE="on"
  export PATH=/usr/local/go/bin:$PATH
  echo "install go completed ..."
}

installWasmLib() {
  $MAKE_ME_ROOT wget --no-check-certificate "https://github.com/CosmWasm/wasmvm/releases/download/v1.0.0/libwasmvm_muslc.x86_64.a" -O /lib/libwasmvm_muslc.x86_64.a
  num=`md5sum /lib/libwasmvm_muslc.x86_64.a |grep f6282df732a13dec836cda1f399dd874b1e3163504dbd9607c6af915b2740479`
  if [[ $num -ne 0 ]]; then
    echo "installWasmLib error md5 not fit"
    exit 1
  fi
  $MAKE_ME_ROOT cp /lib/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.a
}

Prepare() {
  #for curl 56 GnuTLS recv error (-9) gnutls-bin and git config
  git config --global http.version HTTP/1.1
  git config --global http.sslVerify false
  git config --global http.postBuffer 1048576000
  echo "check go version"
  if ! type /usr/local/go/bin/go > /dev/null 2>&1; then
    installgo
  fi
  if ! type go > /dev/null 2>&1; then
    export PATH=$PATH:/usr/local/go/bin
  fi

  v=$(/usr/local/go/bin/go version | { read _ _ v _; echo ${v#go}; })
  # shellcheck disable=SC2046
  if [ $(checkgoversion "$v") -ge $(checkgoversion "1.17") ]
  then
    echo "$v"
    echo "should not install go"
  else
    echo "should install go version above 1.17"
    installgo
  fi

  echo "Prepare completed ...."
}

checkjcmalloc() {
  echo "check jcmalloc ..."
    # shellcheck disable=SC2068
    for lib in ${libArray[@]}
    do
    echo "check lib ${lib}"
    if [ -f "${lib}" ]; then
      # shellcheck disable=SC2126
      has=$($deptool "${lib}" |grep -E '[t|j][c|e]malloc' |wc -l)
      if [ "${has}" -gt 0 ]; then
            uninstallRocksdb
            installRocksdb
            return
      fi
    fi
    done

    echo "check rocksdb lib version"
    # shellcheck disable=SC2068
    for lib in ${libArray[@]}
    do
      if [ ! -f "${lib}" ]; then
        uninstallRocksdb
        installRocksdb
        return
      fi
    done
}

InstallExchain() {
  echo "InstallExchain...."

  download
  cd "$HOME"/.exchain/src/exchain-${ver}
  checkjcmalloc
  #if alpine add LINK_STATICALLY=true
  echo "compile exchain...."
  rm -rf ~/.cache/go-build
  if [ "$dynamicLink" == "TRUE" ]; then
    make mainnet WITH_ROCKSDB=true
  else
    installWasmLib
    make mainnet WITH_ROCKSDB=true LINK_STATICALLY=true
  fi

  if ! type exchaind > /dev/null 2>&1; then
    export PATH=$PATH:$HOME/go/bin
  fi
  echo "InstallExchain completed"
  printLogo
}

GetArchitecture
Prepare
InstallExchain