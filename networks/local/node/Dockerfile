FROM golang:stretch as build-env

# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y $PACKAGES

WORKDIR /okbchain
# Add source files
COPY . .

# build okbchain
RUN make build-linux

# Final image
FROM golang:1.20 as final

WORKDIR /okbchaind
# Copy over binaries from the build-env
COPY --from=build-env /okbchain/build/okbchaind /usr/bin/okbchaind
COPY --from=build-env /okbchain/build/okbchaincli /usr/bin/okbchaincli
COPY --from=build-env /okbchain/networks/local/node/wrapper.sh /usr/bin/wrapper.sh

EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

