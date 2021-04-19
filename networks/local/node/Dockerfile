FROM golang:stretch as build-env

# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y $PACKAGES

WORKDIR /exchain
# Add source files
COPY . .

# build ExChain
RUN make build-linux

# Final image
FROM golang:1.14 as final

WORKDIR /exchaind
# Copy over binaries from the build-env
COPY --from=build-env /exchain/build/exchaind /usr/bin/exchaind
COPY --from=build-env /exchain/build/exchaincli /usr/bin/exchaincli
COPY --from=build-env /exchain/networks/local/node/wrapper.sh /usr/bin/wrapper.sh

EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

