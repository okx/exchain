#!/usr/bin/env bash

set -eo pipefail

protoc_gen_gocosmos() {
  if ! grep "github.com/gogo/protobuf => github.com/regen-network/protobuf" go.mod &>/dev/null ; then
    echo -e "\tPlease run this command from somewhere inside the cosmos-sdk folder."
    return 1
  fi

  go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@latest 2>/dev/null
}

#protoc_gen_gocosmos
#go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@v1.3.3-alpha.regen.1
#go install github.com/gogo/protobuf/gogoproto
#GO111MODULE=on go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0 github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.16.0
#go install github.com/bufbuild/buf/cmd/buf@v0.30.0
#protoc -I "x/vmbridge/proto" -I "third_party/proto" --gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. --grpc-gateway_out=logtostderr=true:. x/vmbridge/proto/vmbridge/wasm/v1/tx.proto

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  buf protoc \
  -I "proto" \
  -I "third_party/proto" \
  --gocosmos_out=plugins=interfacetype+grpc,\
Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  --grpc-gateway_out=logtostderr=true:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')

done
#
## command to generate docs using protoc-gen-doc
buf protoc \
-I "proto" \
-I "third_party/proto" \
--doc_out=./docs/proto \
--doc_opt=./docs/proto/protodoc-markdown.tmpl,proto-docs.md \
$(find "$(pwd)/proto" -maxdepth 5 -name '*.proto')

# move proto files to the right places
cp -r github.com/CosmWasm/wasmd/* ./
rm -rf github.com
