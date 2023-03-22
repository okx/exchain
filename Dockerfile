# Simple usage with a mounted data directory:
# > docker build -t okbchain .
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.okbchaind:/root/.okbchaind -v ~/.okbchaincli:/root/.okbchaincli okbchain okbchaind init mynode
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.okbchaind:/root/.okbchaind -v ~/.okbchaincli:/root/.okbchaincli okbchain okbchaind start
FROM golang:1.20.2-alpine AS build-env

# Install minimum necessary dependencies, remove packages
RUN apk add --no-cache curl make git libc-dev bash gcc linux-headers eudev-dev

# Set working directory for the build
WORKDIR /go/src/github.com/okx/okbchain

# Add source files
COPY . .

ENV GO111MODULE=on \
    GOPROXY=http://goproxy.cn
# Build okbchain
RUN make install

# Final image
FROM alpine:edge

WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/okbchaind /usr/bin/okbchaind
COPY --from=build-env /go/bin/okbchaincli /usr/bin/okbchaincli

# Run okbchaind by default, omit entrypoint to ease using container with okbchaincli
CMD ["okbchaind"]
