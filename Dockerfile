# Simple usage with a mounted data directory:
# > docker build -t okchain .
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.okexchaind:/root/.okexchaind -v ~/.okexchaincli:/root/.okexchaincli okchain okexchaind init mynode
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.okexchaind:/root/.okexchaind -v ~/.okexchaincli:/root/.okexchaincli okchain okexchaind start
FROM golang:alpine AS build-env

# Install minimum necessary dependencies, remove packages
RUN apk add --no-cache curl make git libc-dev bash gcc linux-headers eudev-dev

# Set working directory for the build
WORKDIR /go/src/github.com/okex/okchain

# Add source files
COPY . .

# Build OKChain
RUN GOPROXY=http://goproxy.cn make install

# Final image
FROM alpine:edge

WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/okexchaind /usr/bin/okexchaind
COPY --from=build-env /go/bin/okexchaincli /usr/bin/okexchaincli

# Run okexchaind by default, omit entrypoint to ease using container with okexchaincli
CMD ["okexchaind"]
