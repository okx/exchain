<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Breaking" for breaking API changes.

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [Unreleased]


# [](https://github.com/okex/exchain/compare/v1.1.6...v) (2022-02-07)



## [1.1.6](https://github.com/okex/exchain/compare/v1.1.5...v1.1.6) (2022-01-29)



## [1.1.5](https://github.com/okex/exchain/compare/v1.1.4...v1.1.5) (2022-01-20)



## [1.1.4](https://github.com/okex/exchain/compare/v1.1.3...v1.1.4) (2022-01-12)



## [1.1.3](https://github.com/okex/exchain/compare/v1.1.2...v1.1.3) (2022-01-12)



## [1.1.2](https://github.com/okex/exchain/compare/v1.1.1...v1.1.2) (2022-01-11)



## [1.1.1](https://github.com/okex/exchain/compare/v1.1.0...v1.1.1) (2022-01-10)



# [1.1.0](https://github.com/okex/exchain/compare/v1.0.3...v1.1.0) (2022-01-05)



## [1.0.3](https://github.com/okex/exchain/compare/v1.0.2...v1.0.3) (2022-01-04)



## [1.0.2](https://github.com/okex/exchain/compare/v1.0.1...v1.0.2) (2021-12-29)



## [1.0.1](https://github.com/okex/exchain/compare/v1.0.0...v1.0.1) (2021-12-28)



# [1.0.0](https://github.com/okex/exchain/compare/v0.19.17...v1.0.0) (2021-12-22)



## [0.19.17](https://github.com/okex/exchain/compare/v0.19.16...v0.19.17) (2021-11-29)



## [0.19.16](https://github.com/okex/exchain/compare/v0.19.15...v0.19.16) (2021-11-21)



## [0.19.15](https://github.com/okex/exchain/compare/v0.19.14...v0.19.15) (2021-11-16)



## [0.19.14](https://github.com/okex/exchain/compare/v0.19.13...v0.19.14) (2021-11-13)



## [0.19.13](https://github.com/okex/exchain/compare/v0.19.12...v0.19.13) (2021-11-01)



## [0.19.12](https://github.com/okex/exchain/compare/v0.19.11...v0.19.12) (2021-10-18)



## [0.19.11](https://github.com/okex/exchain/compare/v0.19.10...v0.19.11) (2021-10-14)



## [0.19.10](https://github.com/okex/exchain/compare/v0.19.9...v0.19.10) (2021-10-08)



## [0.19.9](https://github.com/okex/exchain/compare/v0.19.8...v0.19.9) (2021-09-30)



## [0.19.8](https://github.com/okex/exchain/compare/v0.19.7...v0.19.8) (2021-09-29)



## [0.19.7](https://github.com/okex/exchain/compare/v0.19.6...v0.19.7) (2021-09-27)



## [0.19.6](https://github.com/okex/exchain/compare/v0.19.5...v0.19.6) (2021-09-16)



## [0.19.5](https://github.com/okex/exchain/compare/v0.19.2...v0.19.5) (2021-09-09)



## [0.19.2](https://github.com/okex/exchain/compare/v0.19.1...v0.19.2) (2021-08-31)



## [0.19.1](https://github.com/okex/exchain/compare/v0.18.18...v0.19.1) (2021-08-24)

### Features

* Support 0x prefixed address format ([\#973](https://github.com/okex/exchain/pull/973))

* Add logger for tx_pool ([\#983](https://github.com/okex/exchain/pull/983))


### Documentation

* Provide issue template to  report running error ([\#983](https://github.com/okex/exchain/pull/979))

### Code Refactoring

* Deprecated homestead signer after `testnet` block height arrive at `5300000` or `mainnet` arrive at `5150000`. ([\#977](https://github.com/okex/exchain/pull/977))

* Rewrite account nonce query code ([\#989](https://github.com/okex/exchain/pull/989))


### Chores

* Change mercury height in `makefile` ([\#991](https://github.com/okex/exchain/pull/991))


## [0.18.18](https://github.com/okex/exchain/compare/v0.18.17...v0.18.18) (2021-08-18)


### Features

* Push pending tx to kafka ([\#942](https://github.com/okex/exchain/pull/942))


### Bug Fixes

* Fix query failed when using height of pending block ([\#975](https://github.com/okex/exchain/pull/975))


### Code Refactoring

* Change the default param of RPC flag in exchain client ([\#972](https://github.com/okex/exchain/pull/972))


### Chores

* Rewrite the guide in `README.md` ([\#970](https://github.com/okex/exchain/pull/970))

* Optimize `makefile` to avoid set params of genesisHeight in manually ([\#971](https://github.com/okex/exchain/pull/971))



## [0.18.17](https://github.com/okex/exchain/compare/v0.18.16...v0.18.17) (2021-08-16)


### Features

* Support the batch call in a websocket request ([\#957](https://github.com/okex/exchain/pull/957))

* Add pruning flag ([\#964](https://github.com/okex/exchain/pull/964))

* Add local-replay pprof port ([\#965](https://github.com/okex/exchain/pull/965))


### Bug Fixes

* Fix inconsistent error result with ethereum rpc ([\#969](https://github.com/okex/exchain/pull/969))


### Code Refactoring

* Optimize `eth-api` check whether is enabling `txpool` feature ([\#966](https://github.com/okex/exchain/pull/966))

* Remove the traverse of accounts when eth_call ([\#967](https://github.com/okex/exchain/pull/967))



### Chores

* Add `dev/dump.sh` to filter `oec.log` ([\#963](https://github.com/okex/exchain/pull/963))



## [0.18.16](https://github.com/okex/exchain/compare/v0.18.15...v0.18.16) (2021-08-12)


### Features

* Add perform log ([\#959](https://github.com/okex/exchain/pull/959))


### Code Refactoring

* Prune and compact app store ([\#955](https://github.com/okex/exchain/pull/955))

* Expand iavl cache size ([\#960](https://github.com/okex/exchain/pull/960))


### Chores

* Add `dev/start.sh` to start an exchain node ([\#956](https://github.com/okex/exchain/pull/956))



## [0.18.15](https://github.com/okex/exchain/compare/v0.18.14...v0.18.15) (2021-08-10)


### Features

* Add explains of index for dropped txs ([\#953](https://github.com/okex/exchain/pull/953))


### Bug Fixes

* Fix `txpool` unlock mutex is not work ([\#951](https://github.com/okex/exchain/pull/951))

* Fix `txpool` is not return error when is full  ([\#952](https://github.com/okex/exchain/pull/952))



## [0.18.14](https://github.com/okex/exchain/compare/v0.18.13...v0.18.14) (2021-08-06)


### Features

* Add recommend gas price ([\#940](https://github.com/okex/exchain/pull/940))



## [0.18.13](https://github.com/okex/exchain/compare/v0.18.12...v0.18.13) (2021-08-05)


### Bug Fixes

* Fix key mismatch for delete option of black/white list ([\#938](https://github.com/okex/exchain/pull/938))

* Fix `eth_getCode` failed after the contract is blocked ([\#941](https://github.com/okex/exchain/pull/941))



## [0.18.12](https://github.com/okex/exchain/compare/v0.18.11...v0.18.12) (2021-08-03)


### Features

* Add `nonce` support for cosmos tx ([\#931](https://github.com/okex/exchain/pull/931))

* Add `rpc.disable-api` flag ([\#934](https://github.com/okex/exchain/pull/934))



## [0.18.11](https://github.com/okex/exchain/compare/v0.18.10...v0.18.11) (2021-08-01)


### Features

* Add `uuid` in case of duplicate file names ([\#928](https://github.com/okex/exchain/pull/928))


### Bug Fixes

* Fix keys recover bug ([\#926](https://github.com/okex/exchain/pull/926))


## [0.18.10](https://github.com/okex/exchain/compare/v0.18.9...v0.18.10) (2021-07-12)


### Features

* Add `txpool` api ([\#904](https://github.com/okex/exchain/pull/904))

* Observe rpc call duration and record to prometheus ([\#917](https://github.com/okex/exchain/pull/917))

* Add parse app tx by `exchain-amino` ([\#920](https://github.com/okex/exchain/pull/920))


### Bug Fixes

* Open eth_call lru cache when `watcher` enable ([\#909](https://github.com/okex/exchain/pull/909))


### Documentation

* Fix incorrect info of circleci and tag in `README` ([\#919](https://github.com/okex/exchain/pull/919))


### Tests

* Update ut of evm `endblock` ([\#910](https://github.com/okex/exchain/pull/910))



## [0.18.9](https://github.com/okex/exchain/compare/v0.18.8...v0.18.9) (2021-06-22)


### Bug Fixes

* Fix address not found in whitelist after first enable `watchdb` ([\#908](https://github.com/okex/exchain/pull/908))



## [0.18.8](https://github.com/okex/exchain/compare/v0.18.7...v0.18.8) (2021-06-15)


### Features

* Add switch to deleteAccount and deleteState ([\#900](https://github.com/okex/exchain/pull/900))


### Bug Fixes

* Fix query evm tx failed ([\#902](https://github.com/okex/exchain/pull/902))

* Fix generate bloom filter failed occasionally ([\#903](https://github.com/okex/exchain/pull/903))


### Performance Improvements

* Add state lru cache to optimize RPC API `eth_call`  ([\#901](https://github.com/okex/exchain/pull/901))


## [0.18.7](https://github.com/okex/exchain/compare/v0.18.6...v0.18.7) (2021-05-21)


### Features

* Add `txpool` ([\#864](https://github.com/okex/exchain/pull/864))

* Add `watch.db` to operate `account` and `state`  ([\#858](https://github.com/okex/exchain/pull/858))

* Add logical of `evm state` roll back ([\#870](https://github.com/okex/exchain/pull/870))

* Add a resolve in `eth_unsubscribe` ([\#875](https://github.com/okex/exchain/pull/875))

* Add check of watch db enabled in handler ([\#877](https://github.com/okex/exchain/pull/877))

* Add rpc monitor: register metrics to prometheus ([\#884](https://github.com/okex/exchain/pull/884))

* Add mutex in websocket connection, in case of concurrent write ([\#887](https://github.com/okex/exchain/pull/887))


### Bug Fixes

* Fix failed to get block height via the keyword `block.number` ([\#878](https://github.com/okex/exchain/pull/878))

* Fix function error in keywords of `block.number/hash/timestamp` ([\#879](https://github.com/okex/exchain/pull/879))

* Fix invalid state set when the value is zero ([\#881](https://github.com/okex/exchain/pull/881))


### Code Refactoring

* Refactor RPC `eth_call` to optimize qps ([\#857](https://github.com/okex/exchain/pull/857))

* Compatible with `ether client` ([\#888](https://github.com/okex/exchain/pull/888))

* Ensure to concurrent safe `cdc` and config singleton ([\#897](https://github.com/okex/exchain/pull/897))


### Performance Improvements

* Optimize rpc eth_getBalance method ([\#891](https://github.com/okex/exchain/pull/891))

* Add lru cache for `GetCodeByHash` ([\#893](https://github.com/okex/exchain/pull/893))

* Optimize code to avoid re-init cdc ([\#895](https://github.com/okex/exchain/pull/895))

* Add state lru cache to optimize `eth_call` ([\#898](https://github.com/okex/exchain/pull/898))



## [0.18.6](https://github.com/okex/exchain/compare/v0.18.5...v0.18.6) (2021-05-17)


### Build

* Update cosmos version to `v0.39.2-exchain5` and tendermint version to `v0.33.9-exchain4` ([\#869](https://github.com/okex/exchain/pull/869))



## [0.18.5](https://github.com/okex/exchain/compare/v0.18.4...v0.18.5) (2021-05-14)


### Features

* Add gas limit buffer ([\#864](https://github.com/okex/exchain/pull/864))


### Bug Fixes

*  Fix the `newFilter` don't return `log` ([\#856](https://github.com/okex/exchain/pull/856))


### Build

* Update cosmos version to `v0.39.3-0.20210514032300-327d9c09e6b0` ([\#864](https://github.com/okex/exchain/pull/864))



## [0.18.4](https://github.com/okex/exchain/compare/v0.18.3...v0.18.4) (2021-05-10)


### Performance Improvements

* Limit the RPC connection number ([\#853](https://github.com/okex/exchain/pull/853))

* Limit the websocket connection number ([\#855](https://github.com/okex/exchain/pull/855))


### Build

* Update cosmos version to `v0.39.2-exchain4` and tendermint version to `v0.33.9-exchain3` ([\#855](https://github.com/okex/exchain/pull/855))



## [0.18.3](https://github.com/okex/exchain/compare/v0.18.2...v0.18.3) (2021-04-26)


### Bug Fixes

* Fix inconsistent bytecode via call `eth_getCode` ([\#847](https://github.com/okex/exchain/pull/847))


### Performance Improvements

* Update estimateGas upper to 130%  ([\#860](https://github.com/okex/exchain/pull/860))

* Update default gas price  ([\#862](https://github.com/okex/exchain/pull/862))


## [0.18.2](https://github.com/okex/exchain/compare/v0.18.1...v0.18.2) (2021-04-25)


### Build

* Update cosmos version to `v0.39.2-exchain3` and tendermint version to `v0.33.9-exchain2` ([\#837](https://github.com/okex/exchain/pull/837))



## [0.18.1](https://github.com/okex/exchain/compare/v0.18.0...v0.18.1) (2021-04-25)


### Performance Improvements

* `getTransactionReceipt` check that `ContractAddress` is not `0x00000000000000000000` ([\#843](https://github.com/okex/exchain/pull/843))



# [0.18.0](https://github.com/okex/exchain/compare/v0.17.5...v0.18.0) (2021-04-23)


### Features

*  Add v018 migrate ([\#833]((https://github.com/okex/exchain/pull/833)))


### Performance Improvements

* Enhance websocket handle error ([\#838](https://github.com/okex/exchain/pull/833))



## [0.17.5](https://github.com/okex/exchain/compare/v0.17.4...v0.17.5) (2021-04-22)


### Build

* Update cosmos version to `v0.39.2-exchain2` and tendermint version to `v0.33.9-exchain1` ([\#837](https://github.com/okex/exchain/pull/837))



## [0.17.4](https://github.com/okex/exchain/compare/v0.17.3...v0.17.4) (2021-04-20)


### Code Refactoring

* Set the default gas price of `eth_call` to `1 Gwei` ([\#825](https://github.com/okex/exchain/pull/825))

* Remove redundant code ([\#826](https://github.com/okex/exchain/pull/826))



## [0.17.3](https://github.com/okex/exchain/compare/v0.17.2...v0.17.3) (2021-04-15)


### Performance Improvements

* Optimize error of ErrTxTooLarge ([\#820](https://github.com/okex/exchain/pull/820))

* Max world state num ([\#819](https://github.com/okex/exchain/pull/819))



## [0.17.2](https://github.com/okex/exchain/compare/v0.17.1...v0.17.2) (2021-04-13)


### Features

* Support the function of calling web3 api in websocket ([\#795](https://github.com/okex/exchain/pull/795))


### Code Refactoring

* Change nonce format ([\#785](https://github.com/okex/exchain/pull/785))


### Performance Improvements

* Optimize performance of the func eth_call ([\#793](https://github.com/okex/exchain/pull/793))

* Enhance websocket ([\#792](https://github.com/okex/exchain/pull/793))



## [0.17.1](https://github.com/okex/exchain/compare/v0.17.0...v0.17.1) (2021-04-12)


### BREAKING CHANGES

* Rename to exchain ([\#816](https://github.com/okex/exchain/pull/816))



# [0.17.0](https://github.com/okex/exchain/compare/v0.16.9...v0.17.0) (2021-04-11)


### Features

* Add pruning for block state db and add export fro app db ([\#811](https://github.com/okex/exchain/pull/811))


### Build

* Update cosmos-sdk version to `v0.39.2-exchain1` ([\#814](https://github.com/okex/exchain/pull/814))


