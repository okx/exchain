# Simple usage with a mounted data directory:
# > docker build -t okchain .
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.okchaind:/root/.okchaind -v ~/.okchaincli:/root/.okchaincli okchain okchaind init mynode
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.okchaind:/root/.okchaind -v ~/.okchaincli:/root/.okchaincli okchain okchaind start
FROM golang:1.13-buster

# Install minimum necessary dependencies, remove packages
RUN apt update
RUN apt install -y curl git build-essential vim

# Set working directory for the build
WORKDIR /go/src/github.com/okex/okchain

# Add source files
COPY . .

# Build OKChain
RUN GOPROXY=http://goproxy.cn make install

# Install libgo_cosmwasm.so to a shared directory where it is readable by all users
# See https://github.com/CosmWasm/wasmd/issues/43#issuecomment-608366314
RUN cp /go/pkg/mod/github.com/okex/go-cosmwasm@v*/api/libgo_cosmwasm.so /lib/x86_64-linux-gnu/

WORKDIR /root

RUN cp /go/bin/okchaind /usr/bin/okchaind
RUN cp /go/bin/okchaincli /usr/bin/okchaincli

# Run okchaind by default, omit entrypoint to ease using container with okchaincli
CMD ["okchaind"]
