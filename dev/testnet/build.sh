build(){
    branch=$1
    memo=$2
    cwd=$(cd "$(dirname $0)";pwd)
    rm -rf $cwd/images/$memo

    git clone https://github.com/okex/exchain.git -b $branch $cwd/images/$memo/exchain

    cd $cwd/images/$memo/exchain && make build 

    mv $cwd/images/$memo/build/exchaind $cwd/images/$memo
    mv $cwd/images/$memo/build/exchaincli $cwd/images/$memo

    # cd cwd
}

build "v1.1.4.1" "test"