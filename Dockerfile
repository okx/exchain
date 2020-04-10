# Simple usage with a mounted data directory:
# > docker build -t okchain .
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.okchaind:/root/.okchaind -v ~/.okchaincli:/root/.okchaincli okchain okchaind init mynode
# > docker run -it -p 36657:36657 -p 36656:36656 -v ~/.okchaind:/root/.okchaind -v ~/.okchaincli:/root/.okchaincli okchain okchaind start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev

# Set working directory for the build
WORKDIR /go/src/github.com/okex/okchain

# Add source files
COPY . .

# Install minimum necessary dependencies, build OKChain, remove packages
RUN apk add --no-cache $PACKAGES && \
    export GO111MODULE=off && \
    make install

# Final image
FROM alpine:edge

WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/okchaind /usr/bin/okchaind
COPY --from=build-env /go/bin/okchaincli /usr/bin/okchaincli

# Run okchaind by default, omit entrypoint to ease using container with okchaincli
CMD ["okchaind"]
