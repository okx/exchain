# Simple usage with a mounted data directory:
# > docker build -t exchain .
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.exchaind:/root/.exchaind -v ~/.exchaincli:/root/.exchaincli exchain exchaind init mynode
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.exchaind:/root/.exchaind -v ~/.exchaincli:/root/.exchaincli exchain exchaind start
FROM golang:alpine AS build-env

# Install minimum necessary dependencies, remove packages
RUN apk add --no-cache curl make git libc-dev bash gcc linux-headers eudev-dev libstdc++ g++ bzip2-dev gflags-dev lz4-dev snappy-dev zlib-dev zstd-dev perl

# Set working directory for the build
WORKDIR /go/src/github.com/okex/exchain

# Add source files
COPY . .

# Build OKExChain
RUN make rocksdb
RUN GOPROXY=http://goproxy.cn make mainnet WITH_ROCKSDB=true

# Final image
FROM alpine:edge

RUN apk add --no-cache bzip2-dev gflags lz4-libs snappy zlib zstd-libs

WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /usr/local/lib/librocksdb.so* /usr/lib/
COPY --from=build-env /go/bin/exchaind /usr/bin/exchaind
COPY --from=build-env /go/bin/exchaincli /usr/bin/exchaincli

# Run exchaind by default, omit entrypoint to ease using container with exchaincli
CMD ["exchaind", "--db_backend=rocksdb"]
