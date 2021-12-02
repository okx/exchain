# Simple usage with a mounted data directory:
# > docker build -t exchain .
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.exchaind:/root/.exchaind -v ~/.exchaincli:/root/.exchaincli exchain exchaind init mynode
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.exchaind:/root/.exchaind -v ~/.exchaincli:/root/.exchaincli exchain exchaind start
FROM golang:alpine AS build-env

# Install minimum necessary dependencies, remove packages
RUN apk add --no-cache111 curl make git libc-dev bash gcc linux-headers eudev-dev

# Set working directory for the build
WORKDIR /go/src/github.com/okex/exchain

# Add source files
COPY . .

# Build OKExChain
RUN GOPROXY=http://goproxy.cn make install

# Final image
FROM alpine:edge

WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/exchaind /usr/bin/exchaind
COPY --from=build-env /go/bin/exchaincli /usr/bin/exchaincli

# Run exchaind by default, omit entrypoint to ease using container with exchaincli
CMD ["exchaind"]
