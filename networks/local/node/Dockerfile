FROM golang:stretch as build-env

# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y $PACKAGES

WORKDIR /okexchain
# Add source files
COPY . .

# build OKExChain
RUN make build-linux

# Final image
FROM golang:1.14 as final

WORKDIR /okexchaind
# Copy over binaries from the build-env
COPY --from=build-env /okexchain/build/okexchaind /usr/bin/okexchaind
COPY --from=build-env /okexchain/build/okexchaincli /usr/bin/okexchaincli
COPY --from=build-env /okexchain/networks/local/node/wrapper.sh /usr/bin/wrapper.sh

EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

