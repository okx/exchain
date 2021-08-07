maDEP := $(shell command -v dep 2> /dev/null)
SUM := $(shell which shasum)

COMMIT := $(shell git rev-parse HEAD)
CAT := $(if $(filter $(OS),Windows_NT),type,cat)
export GO111MODULE=on

GithubTop=github.com


Version=v0.18.14
CosmosSDK=v0.39.2
Tendermint=v0.33.9
Iavl=v0.14.3
Name=exchain
ServerName=exchaind
ClientName=exchaincli
# the height of the 1st block is GenesisHeight+1
GenesisHeight=0
MercuryHeight=0

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


ldflags = -X $(GithubTop)/cosmos/cosmos-sdk/version.Version=$(Version) \
	-X $(GithubTop)/cosmos/cosmos-sdk/version.Name=$(Name) \
  -X $(GithubTop)/cosmos/cosmos-sdk/version.ServerName=$(ServerName) \
  -X $(GithubTop)/cosmos/cosmos-sdk/version.ClientName=$(ClientName) \
  -X $(GithubTop)/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
  -X $(GithubTop)/cosmos/cosmos-sdk/version.CosmosSDK=$(CosmosSDK) \
  -X $(GithubTop)/cosmos/cosmos-sdk/version.Tendermint=$(Tendermint) \
  -X "$(GithubTop)/cosmos/cosmos-sdk/version.BuildTags=$(build_tags)" \
  -X $(GithubTop)/tendermint/tendermint/types.startBlockHeightStr=$(GenesisHeight) \
  -X $(GithubTop)/cosmos/cosmos-sdk/types.MILESTONE_MERCURY_HEIGHT=$(MercuryHeight) \


ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -ldflags '$(ldflags)'  -gcflags "all=-N -l"
BUILD_TESTNET_FLAGS := $(BUILD_FLAGS)

all: install

install: exchain

exchain:
	go install -v $(BUILD_FLAGS) -tags "$(BUILD_TAGS)" ./cmd/exchaind
	go install -v $(BUILD_FLAGS) -tags "$(BUILD_TAGS)" ./cmd/exchaincli

testnet:
	go install -v $(BUILD_TESTNET_FLAGS) -tags "$(BUILD_TAGS)" ./cmd/exchaind
	go install -v $(BUILD_TESTNET_FLAGS) -tags "$(BUILD_TAGS)" ./cmd/exchaincli

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./app/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/backend/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/common/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/dex/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/distribution/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/genutil/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/gov/...
#	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/order/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/params/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/staking/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/token/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/upgrade/...
#	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./vendor/github.com/cosmos/cosmos-sdk/x/mint/...

get_vendor_deps:
	@echo "--> Generating vendor directory via dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -vendor-only

update_vendor_deps:
	@echo "--> Running dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -update

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download
.PHONY: go-mod-cache

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

cli:
	go install -v $(BUILD_FLAGS) -tags "$(BUILD_TAGS)" ./cmd/exchaincli

server:
	go install -v $(BUILD_FLAGS) -tags "$(BUILD_TAGS)" ./cmd/exchaind

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s

build:
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -o build/exchaind.exe ./cmd/exchaind
	go build $(BUILD_FLAGS) -o build/exchaincli.exe ./cmd/exchaincli
else
	go build $(BUILD_FLAGS) -o build/exchaind ./cmd/exchaind
	go build $(BUILD_FLAGS) -o build/exchaincli ./cmd/exchaincli
endif

build-linux:
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

build-docker-exchainnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/exchaind/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/exchaind:Z exchain/node testnet --v 4 -o . --starting-ip-address 192.168.10.2 --keyring-backend=test ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down


.PHONY: build
