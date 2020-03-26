maDEP := $(shell command -v dep 2> /dev/null)
SUM := $(shell which shasum)

COMMIT := $(shell git rev-parse HEAD)
CAT := $(if $(filter $(OS),Windows_NT),type,cat)

# this should be the same as the version in go.mod
Version=0.10.0
Tendermint=0.32.7
CosmosSDK=0.37.4
ServerName=okchaind
ClientName=okchaincli
#StartBlockHeight=2000

# process linker flags
ifeq ($(VERSION),)
    VERSION = $(COMMIT)
endif

build_tags = netgo

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

ldflags = -X github.com/okex/okchain/vendor/github.com/cosmos/cosmos-sdk/version.Version=$(Version) \
  -X github.com/okex/okchain/vendor/github.com/cosmos/cosmos-sdk/version.ServerName=$(ServerName) \
  -X github.com/okex/okchain/vendor/github.com/cosmos/cosmos-sdk/version.ClientName=$(ClientName) \
  -X github.com/okex/okchain/vendor/github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
  -X github.com/okex/okchain/vendor/github.com/cosmos/cosmos-sdk/version.CosmosSDK=$(CosmosSDK) \
  -X github.com/okex/okchain/vendor/github.com/cosmos/cosmos-sdk/version.Tendermint=$(Tendermint) \
  -X github.com/okex/okchain/vendor/github.com/cosmos/cosmos-sdk/version.VendorDirHash=$(shell $(SUM) -a 256 go.sum | cut -d ' ' -f1) \
  -X "github.com/okex/okchain/vendor/github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags)" \
  -X "github.com/okex/okchain/vendor/github.com/tendermint/tendermint/types.startBlockHeightStr=$(StartBlockHeight)"

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -ldflags '$(ldflags)'


all: check install

install: okchain

okchain:
	go install -v $(BUILD_FLAGS) -tags "$(BUILD_TAGS)" ./cmd/okchaind
	go install -v $(BUILD_FLAGS) -tags "$(BUILD_TAGS)" ./cmd/okchaincli

check:
	cd unittest; make $@

checko:
	cd unittest; make $@

checks:
	cd unittest; make $@

checkt:
	cd unittest; make $@

checkd:
	cd unittest; make $@

checkdi:
	cd unittest; make $@

verify:
	cd unittest; make $@

profile:
	cd unittest; make $@


get_vendor_deps:
	@echo "--> Generating vendor directory via dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -vendor-only

update_vendor_deps:
	@echo "--> Running dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -update


cli:
	go install -v $(BUILD_FLAGS) -tags "$(BUILD_TAGS)" ./cmd/okchaincli

server:
	go install -v $(BUILD_FLAGS) -tags "$(BUILD_TAGS)" ./cmd/okchaind


tiger:
	go install -v $(BUILD_FLAGS) -tags "$(BUILD_TAGS)" ./cmd/tiger

launchcmd:
	go install -v $(BUILD_FLAGS) -tags "$(BUILD_TAGS)" ./cmd/launch



format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s


.PHONY: build
