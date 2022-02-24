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

# [](https://github.com/okex/exchain/compare/v1.1.8...v) (2022-02-23)


### Code refactoring

* Remove the deprecated ethermint transaction handler from the code and remove `handleSimulation` ([\#1553](https://github.com/okex/exchain/pull/1553))


### Chores

* Support running high tps test network locally, the highest tps can reach 5800([\#1585](https://github.com/okex/exchain/pull/1585))



## [1.1.8](https://github.com/okex/exchain/compare/v1.1.7...v1.1.8) (2022-02-21)


### Features

* Tax rewards is distributed to `treasures` proposal ([\#1523](https://github.com/okex/exchain/pull/1523))

* Add flag `r` to set the number of RPC nodes when initialize the `testnet`  ([\#1535](https://github.com/okex/exchain/pull/1535))

* Add flag `goleveldb` for cmd of data query ([\#1538](https://github.com/okex/exchain/pull/1538))

* Support different account type for `eth_getBalanceBatch` ([\#1540](https://github.com/okex/exchain/pull/1540))

* Add flag `delta-version` to reduce using wrong delta-version cause panic ([\#1546](https://github.com/okex/exchain/pull/1546))

* Refine function `handleSimulation`  ([\#1549](https://github.com/okex/exchain/pull/1549))

* Support `prerun` compatible with parallel tx ([\#1555](https://github.com/okex/exchain/pull/1555))

* Add `GetType` to transaction interface  ([\#1561](https://github.com/okex/exchain/pull/1561))

* Add `ante` statistics log ([\#1572](https://github.com/okex/exchain/pull/1572))


### Bug fixes

* Fix repair state on start when state machine broken ([\#1563](https://github.com/okex/exchain/pull/1563))


### Styles

* Format all code by `gofmt` ([\#1536](https://github.com/okex/exchain/pull/1536))


### Code refactoring

* Full check for cosmos tx when using `wtx` ([\#1545](https://github.com/okex/exchain/pull/1545))

* Reorganize `ante` decorator ([\#1562](https://github.com/okex/exchain/pull/1562))

* Refactor delete `watchdb` data in `x/evm/watcher` ([\#1575](https://github.com/okex/exchain/pull/1575))

* Update baseApp `checkTx` height ([\#1581](https://github.com/okex/exchain/pull/1581))


### Performance Improvements

* Optimize `NodeToNodeJson` by using unsafe pointer in `dds` ([\#1551](https://github.com/okex/exchain/pull/1551))

* Optimize for `handleMsgEthereumTx` refund ([\#1552](https://github.com/okex/exchain/pull/1552))

* Optimize ethTx marshal in `x/evm/watcher` ([\#1560](https://github.com/okex/exchain/pull/1560))



## [1.1.7](https://github.com/okex/exchain/compare/v1.1.6...v1.1.7) (2022-02-07)


### Features

* Support for debugging trace tx ([\#1427](https://github.com/okex/exchain/pull/1427))

* Support check flag conflict of `start` command ([\#1504](https://github.com/okex/exchain/pull/1504))

* Add `wtx` whitelist ([\#1532](https://github.com/okex/exchain/pull/1532))


### Bug fixes

* Fix `nodeid` and `wtx` statistics ([\#1525](https://github.com/okex/exchain/pull/1525))

* Fix RPC-API `eth_getTransactionCount` method ([\#1530](https://github.com/okex/exchain/pull/1530))

* Fix `tx_index` DB type ([\#1531](https://github.com/okex/exchain/pull/1531))


### Tests

* Add amino ut of state delta ([\#1520](https://github.com/okex/exchain/pull/1520))

* Add `tm-db` ut ([\#1529](https://github.com/okex/exchain/pull/1529))



## [1.1.6](https://github.com/okex/exchain/compare/v1.1.5...v1.1.6) (2022-01-29)


### Features

* Support handle paralleled-tx error when run `txDecode` failed  ([\#1453](https://github.com/okex/exchain/pull/1453))

* Extract current node key or from specific file ([\#1478](https://github.com/okex/exchain/pull/1478))

* Add `GetUnsafeValue` in database interface ([\#1495](https://github.com/okex/exchain/pull/1495))

* Add `jobChan` to control commit `watchData` ([\#1494](https://github.com/okex/exchain/pull/1494))


### Bug fixes

* Fix bug to set global height when repair state ([\#1479](https://github.com/okex/exchain/pull/1479))

* Fix rocksdb out of memory ([\#1483](https://github.com/okex/exchain/pull/1483))

* Fix bug in `x/auth` when using `checkTx` ([\#1485](https://github.com/okex/exchain/pull/1485))

* Fix bug that don't set `query-fast` flag when the node start with archive mode ([\#1496](https://github.com/okex/exchain/pull/1496))

* Fix `dds` check before using `DeltaMap` ([\#1509](https://github.com/okex/exchain/pull/1509))


### Code refactoring

* Reorganize `amino` code ([\#1373](https://github.com/okex/exchain/pull/1373))

* Refine `dds` encode and decode code ([\#1482](https://github.com/okex/exchain/pull/1482))

* Refactor  `CheckTx` and `broadcast` in `mempool` ([\#1484](https://github.com/okex/exchain/pull/1484))

* Change `dds` encoder from`json` to `amino` ([\#1510](https://github.com/okex/exchain/pull/1510))


### Performance Improvements

* Put marshal into upload routine ([\#1472](https://github.com/okex/exchain/pull/1472))

* Get rid of remote `abci` client ([\#1481](https://github.com/okex/exchain/pull/1481))

* Update `tm-db` to `v0.5.2-oec1` for rocksdb `mmmap` options ([\#1491](https://github.com/okex/exchain/pull/1491))

* Deduplicate in `x/evm/watcher`  ([\#1514](https://github.com/okex/exchain/pull/1514))

* Delete unused code in `wtx` ([\#1516](https://github.com/okex/exchain/pull/1516))


### Tests

* Fix TestTxProofs in `libs/tendermint/lite/proxy` ([\#1502](https://github.com/okex/exchain/pull/1502))


### Chores

* Add script to test `oec` multiple nodes upgrading ([\#1477](https://github.com/okex/exchain/pull/1477))

* Refactor rocksdb m1 repair patch file path ([\#1511](https://github.com/okex/exchain/pull/1511))


### Breaking Changes

* Move `tmdb` to `libs/tm-db` ([\#1493](https://github.com/okex/exchain/pull/1493))



## [1.1.5](https://github.com/okex/exchain/compare/v1.1.4...v1.1.5) (2022-01-20)


### Features

* Add `flatkv` storage for reading performance ([\#1357](https://github.com/okex/exchain/pull/1357))

* Add `wtx` to reduce checking-signatures time of trustable tx  ([\#1429](https://github.com/okex/exchain/pull/1429))

* Add delta redis-db ([\#1430](https://github.com/okex/exchain/pull/1430))

* Add `blockhash` in the `executionTask` of `prerun` ([\#1448](https://github.com/okex/exchain/pull/1448))


### Bug fixes

* Fix `store` and `state` heights unmatched ([\#1432](https://github.com/okex/exchain/pull/1432))

* Fix mismatch hash of pending tx ([\#1435](https://github.com/okex/exchain/pull/1435))

* Removed `GetValidator` caching to fix concurrency error ([\#1447](https://github.com/okex/exchain/pull/1447))

* Fix missing version when ac enabled ([\#1451](https://github.com/okex/exchain/pull/1451))

* Fix incorrect nonce ([\#1454](https://github.com/okex/exchain/pull/1454))

* Fix `prerun` panic when reset deliverState ([\#1457](https://github.com/okex/exchain/pull/1457))

* Fix wrong `AppHash` when node has been restarted ([\#1469](https://github.com/okex/exchain/pull/1469))


### Performance Improvements

* RPC optimize by using `amino` encoding data ([\#1326](https://github.com/okex/exchain/pull/1326))

* Clear unused code out in the directory `x` ([\#1412](https://github.com/okex/exchain/pull/1412))


### Tests

* Add ut of `delta` ([\#1355](https://github.com/okex/exchain/pull/1355))

* Add sub-process ut template ([\#1397](https://github.com/okex/exchain/pull/1397))

* Benchmark test about encoding tx performance among `go-amino`,`rlp`,`exchain-amino` and `json` ([\#1426](https://github.com/okex/exchain/pull/1426))

* Fix ut in the BeginBlock when using `prerun`  ([\#1463](https://github.com/okex/exchain/pull/1463))


### Chores

* Update `dockerfile` ([\#1456](https://github.com/okex/exchain/pull/1456))



## [1.1.4](https://github.com/okex/exchain/compare/v1.1.3...v1.1.4) (2022-01-12)


### Bug fixes

* Fix tx decode in query to make compatible with all tx encode ([\#1424](https://github.com/okex/exchain/pull/1424))



## [1.1.3](https://github.com/okex/exchain/compare/v1.1.2...v1.1.3) (2022-01-12)


### Features

* Support build rocksdb on `m1` ([\#1416](https://github.com/okex/exchain/pull/1416))


### Bug fixes 

* Fix `txpool` bug  before venus height ([\#1421](https://github.com/okex/exchain/pull/1421))


### Code refactoring

* Change default `MercuryHeight` and `VenusHeight` to 1 to make local code execution consistent with `mainnet` ([\#1413](https://github.com/okex/exchain/pull/1413))


### Tests

* Add ut related to `txhash` ([\#1418](https://github.com/okex/exchain/pull/1418))


## [1.1.2](https://github.com/okex/exchain/compare/v1.1.1...v1.1.2) (2022-01-11)


### Features

* Support `gen_proof_mem` when using replay ([\#1405](https://github.com/okex/exchain/pull/1405))

* Add 2 hours buffer of `VenusHeight` upgrade to close `txPool` ([\#1415](https://github.com/okex/exchain/pull/1415))


## [1.1.1](https://github.com/okex/exchain/compare/v1.1.0...v1.1.1) (2022-01-10)


### Features

* Support `bech32` address convert to other kinds ([\#1384](https://github.com/okex/exchain/pull/1384))

* Add check log of producer and consumer in `wathcdb` ([\#1393](https://github.com/okex/exchain/pull/1393))

* `RLP` encode is not allowed until venus height ([\#1407](https://github.com/okex/exchain/pull/1407))

* `amino` encode for `MsgEthereumTx` is not supported after venus height ([\#1410](https://github.com/okex/exchain/pull/1410))


### Bug fixes

* Fix bug that can't find contract address when using counter contract ([\#1392](https://github.com/okex/exchain/pull/1392))

* Fix `chain-id` check ([\#1399](https://github.com/okex/exchain/pull/1399))


### Code refactoring

* Update `txpool` encoder ([\#1396](https://github.com/okex/exchain/pull/1396))

* Update `dds` flag code ([\#1398](https://github.com/okex/exchain/pull/1398))

* Refine `RLP` encode and decode function ([\#1400](https://github.com/okex/exchain/pull/1400))

* Change `GetTxEncoder` from `RLP` to `amino` in function `doCall` ([\#1401](https://github.com/okex/exchain/pull/1401))


### Tests

* Add contract case in `watchdata` ut ([\#1387](https://github.com/okex/exchain/pull/1387))



# [1.1.0](https://github.com/okex/exchain/compare/v1.0.3...v1.1.0) (2022-01-05)


### Features

* Add flag to close log analyzer ([\#1321](https://github.com/okex/exchain/pull/1321))


### Bug fixes

* Fix bug when pre-run with watch data happen  ([\#1386](https://github.com/okex/exchain/pull/1386))


### Code refactoring

* Change tx hash to `Keccak256` for ethereum compatibility ([\#1350](https://github.com/okex/exchain/pull/1350))


### Tests

* Add `watchdata` ut ([\#1364](https://github.com/okex/exchain/pull/1364))

* Close the `WAL` ut of consensus ([\#1371](https://github.com/okex/exchain/pull/1371))

* Fix the  unit test `txhash` ([\#1385](https://github.com/okex/exchain/pull/1385))



## [1.0.3](https://github.com/okex/exchain/compare/v1.0.2...v1.0.3) (2022-01-04)


### Features

* RPC-API `personal_newaccount` is compatible with eth ([\#1333](https://github.com/okex/exchain/pull/1333))

* Calculate delta-payload hash when marshal and unmarshal ([\#1340](https://github.com/okex/exchain/pull/1340))

* Add standard mode in testing `OIP20` contract  ([\#1359](https://github.com/okex/exchain/pull/1359)) 

* Add `dds` statistics ([\#1368](https://github.com/okex/exchain/pull/1368))

* Support subscribe log event from `kafka` ([\#1376](https://github.com/okex/exchain/pull/1376))


### Bug fixes

* Fix `oec` client bugs about deploy contract when switch to fast-sync mode from consensus mode ([\#1354](https://github.com/okex/exchain/pull/1354))

* Fix REST-API  `get blocked list` ([\#1366](https://github.com/okex/exchain/pull/1366))

* Fix `MarshalResponseEndBlockToAmino` bug ([\#1372](https://github.com/okex/exchain/pull/1372))


### Documentation

* Add "How to build a private chain" in `README.md` ([\#1369](https://github.com/okex/exchain/pull/1369))


### Code refactoring

* Restructure `inner tx `([\#1335](https://github.com/okex/exchain/pull/1335))

* Rename `lrp` flag ([\#1378](https://github.com/okex/exchain/pull/1378))


### Tests

* Fix ut ci and tendermint module ut ([\#1334](https://github.com/okex/exchain/pull/1334))

* Disable `TestReactorBroadcastTxMessage` ut function ([\#1365](https://github.com/okex/exchain/pull/1365))



## [1.0.2](https://github.com/okex/exchain/compare/v1.0.1...v1.0.2) (2021-12-29)


### Features

* Get chain-id from `gendoc` ([\#1346](https://github.com/okex/exchain/pull/1346))

* Add dds statistic ([\#1352](https://github.com/okex/exchain/pull/1352))


### Chores

* Refine `rdb` install  ([\#1344](https://github.com/okex/exchain/pull/1344))

* Update `testnet.sh` ([\#1349](https://github.com/okex/exchain/pull/1349))



## [1.0.1](https://github.com/okex/exchain/compare/v1.0.0...v1.0.1) (2021-12-28)


### Features

* Add flag to init redis params  ([\#1338](https://github.com/okex/exchain/pull/1338))


### Bug fixes

* Fix bug that don't upload when is `fastSync` ([\#1339](https://github.com/okex/exchain/pull/1339))


### Code refactoring

* Refactor `prerun` code  ([\#1341](https://github.com/okex/exchain/pull/1341))



# [1.0.0](https://github.com/okex/exchain/compare/v0.19.17...v1.0.0) (2021-12-22)


### Features

* Delta sync ([\#1153](https://github.com/okex/exchain/pull/1153))

* Support submit an update contract method blocked list proposal  ([\#1182](https://github.com/okex/exchain/pull/1182))

* Support redis dds in `state-delta` ([\#1249](https://github.com/okex/exchain/pull/1249))

* Validate `Chain-id` and `GenesisHeight` while app inits and starts ([\#1258](https://github.com/okex/exchain/pull/1258))

* Support proactively run tx in `exchaincli` ([\#1271](https://github.com/okex/exchain/pull/1271))

* Add flag `flagMnemonic` to hard-code the mnemonic of first 4 validators when start testnet ([\#1278](https://github.com/okex/exchain/pull/1278))

* Add dds tools ([\#1279](https://github.com/okex/exchain/pull/1279))

* Add `compress` to compress delta bytes ([\#1287](https://github.com/okex/exchain/pull/1287))

* Add `delta-version` to separate  different encoding delta data  ([\#1297](https://github.com/okex/exchain/pull/1297))

* Add RPC-API `eth_getBalanceBatch` ([\#1288](https://github.com/okex/exchain/pull/1288))

* Add auto test tool ([\#1302](https://github.com/okex/exchain/pull/1302))

* Add `dds` log ([\#1315](https://github.com/okex/exchain/pull/1315))


### Bug fixes

* Fix `mempool` mutex ([\#1238](https://github.com/okex/exchain/pull/1238))

* Fix the bug of abnormal gasUsed statistics ([\#1266](https://github.com/okex/exchain/pull/1266))

* Fix bug for using `gov` params  from `cosmos-sdk/x/gov` ([\#1298](https://github.com/okex/exchain/pull/1298))

* Fix bug for contract-method delete multiply for change cache struct case to error ([\#1309](https://github.com/okex/exchain/pull/1309))

* Fix bug that enable `fastquery` when don't use delta  ([\#1311](https://github.com/okex/exchain/pull/1311))


### Code refactoring

* Move repair state on `exchaindcli` start   ([\#1205](https://github.com/okex/exchain/pull/1205))

* Change flag about `state-delta` ([\#1247](https://github.com/okex/exchain/pull/1247))

* Re-organize `state-delta` code ([\#1277](https://github.com/okex/exchain/pull/1277))

* Change delta from `p2p` to `deltaContext`  ([\#1280](https://github.com/okex/exchain/pull/1280))

* Refine d`ds  ([\#1284](https://github.com/okex/exchain/pull/1284))

* Get rid of panic when using `dds`  ([\#1295](https://github.com/okex/exchain/pull/1295))

* Cache `account`, `contract` and `code` multiply ([\#1300](https://github.com/okex/exchain/pull/1300))

* Refactor `runtx`  ([\#1306](https://github.com/okex/exchain/pull/1306))

* Refactor consensus test case ([\#1314](https://github.com/okex/exchain/pull/1314))

* Refactor delta download ([\#1317](https://github.com/okex/exchain/pull/1317))


### Performance Improvements

* Reduce the sleep duration while time for `ApplyBlock` is less than `CommitTimeout` ([\#1221](https://github.com/okex/exchain/pull/1221))

* Amino codec optimize ([\#1234](https://github.com/okex/exchain/pull/1234))

* GC optimize about amino codec and `iavl` amino ([\#1246](https://github.com/okex/exchain/pull/1246))

* GC optimize about keccak cache  ([\#1261](https://github.com/okex/exchain/pull/1261))


### Chores

* Update `install-rocksdb.sh` ([\#1256](https://github.com/okex/exchain/pull/1256))



## [0.19.17](https://github.com/okex/exchain/compare/v0.19.16...v0.19.17) (2021-11-29)


### Features

* Add checkTxCnt and `mempool` txs count to the log ([\#1224](https://github.com/okex/exchain/pull/1224))

* Add the flag of max gas used per block to replay ([\#1231](https://github.com/okex/exchain/pull/1231))

* Register `oec` config to dynamic config ([\#1232](https://github.com/okex/exchain/pull/1232))

* Query the current blocked list of contract addresses during evm calling  ([\#1276](https://github.com/okex/exchain/pull/1276))


### Bug fixes

* Fix panic when remove a removed tx ([\#1213](https://github.com/okex/exchain/pull/1213))

* Fix `mempool` full issue ([\#1220](https://github.com/okex/exchain/pull/1220))

* Use `Stringer` instead of `fmt` to reduce unnecessary string construction at non-debug levels ([\#1222](https://github.com/okex/exchain/pull/1222))

* Fix not display proposer address when the node is not proposer ([\#1223](https://github.com/okex/exchain/pull/1223))

* Fix `eth_getTransactionReceipt` gasUsed ([\#1230](https://github.com/okex/exchain/pull/1230))

* Fix `mempool` mutex ([\#1244](https://github.com/okex/exchain/pull/1244))


### Code refactoring

* Close old pruning logic when open ac ([\#1211](https://github.com/okex/exchain/pull/1211))

* Update node cache size dynamically  ([\#1225](https://github.com/okex/exchain/pull/1225))

* Refactor `abci` mutex ([\#1237](https://github.com/okex/exchain/pull/1237))


### Performance Improvements

* Optimising API 'eth_getTransactionReceiptsByBlock' ([\#1217](https://github.com/okex/exchain/pull/1217))

* Optimize node mode  ([\#1226](https://github.com/okex/exchain/pull/1226))

* Optimize pruning ([\#1239](https://github.com/okex/exchain/pull/1239))


### Tests

* Benchmark `iavl` ([\#1219](https://github.com/okex/exchain/pull/1219))


### Chores

* Update `install-rocksdb.sh` ([\#1270](https://github.com/okex/exchain/pull/1270))



## [0.19.16](https://github.com/okex/exchain/compare/v0.19.15...v0.19.16) (2021-11-21)


### Features

* Add amino custom marshall and unmarshal function in the base types  ([\#1146](https://github.com/okex/exchain/pull/1146))

* Add db to `gasused` for calculation ([\#1165](https://github.com/okex/exchain/pull/1165))

* Add `MaxTxNumPerBlock` and `MaxGasUsedPerBlock` in dynamic config to forced flushing `mempool` ([\#1204](https://github.com/okex/exchain/pull/1204))

* Set the default value of `p2p.seeds` ([\#1206](https://github.com/okex/exchain/pull/1206))

* Display block producer address ([\#1212](https://github.com/okex/exchain/pull/1212))


### Bug fixes

* Fix `mempool` config isn't set  ([\#1209](https://github.com/okex/exchain/pull/1209))


### Code refactoring

* Change RPC prometheus histogram buckets for monitor ([\#1207](https://github.com/okex/exchain/pull/1207))



## [0.19.15](https://github.com/okex/exchain/compare/v0.19.14...v0.19.15) (2021-11-16)


### Features

* Support `FlagRocksdbOpts` flag ([\#1203](https://github.com/okex/exchain/pull/1203))

* Add log of `iavl-height` ([\#1208](https://github.com/okex/exchain/pull/1208))


### Bug fixes

* Fix p2p sanity error when address ID is same  ([\#1193](https://github.com/okex/exchain/pull/1193))

* Fix panic when `checkTx` read map ([\#1196](https://github.com/okex/exchain/pull/1196))

* Fix without gracefully exit ([\#1197](https://github.com/okex/exchain/pull/1197))

* Fix the usage of `sed` in Ubuntu when using `start.sh` ([\#1198](https://github.com/okex/exchain/pull/1198))

* Fix consensus issue with repaired state ([\#1199](https://github.com/okex/exchain/pull/1199))

* Add the block producer flag to control consensus time ([\#1201](https://github.com/okex/exchain/pull/1201))


### Code refactoring

* Use string literals instead of `analyzer.RunFuncName` ([\#1195](https://github.com/okex/exchain/pull/1195))



### Performance Improvements

* Optimize cache signCache



## [0.19.14](https://github.com/okex/exchain/compare/v0.19.13...v0.19.14) (2021-11-13)


### Features

* Support batch transaction search ([\#1156](https://github.com/okex/exchain/pull/1156))

* Convert pin keys to const type stored in `libs/cosmos/baseapp/const.go` ([\#1159](https://github.com/okex/exchain/pull/1159))

* Support `innertx` ([\#1161](https://github.com/okex/exchain/pull/1161))

* Add the time of reading database ([\#1170](https://github.com/okex/exchain/pull/1170))

* Get `evm` execute trace and save the trace to database ([\#1171](https://github.com/okex/exchain/pull/1171))

* Add node mode for flags management ([\#1173](https://github.com/okex/exchain/pull/1173))

* Add pruning log and invalid-tx log ([\#1189](https://github.com/okex/exchain/pull/1189))


### Bug fixes

* Fix BlockGasMeter show about `ParallelTx` ([\#1154](https://github.com/okex/exchain/pull/1154))

* Fix `gasUsed` when tx isn't used in `x/evm` ([\#1164](https://github.com/okex/exchain/pull/1164))

* Fix mismatched block hash from block filter ([\#1166](https://github.com/okex/exchain/pull/1166))

* Fix bug export height(-1) state failed ([\#1179](https://github.com/okex/exchain/pull/1179))

* Fix `iavl` log not initialized before used ([\#1181](https://github.com/okex/exchain/pull/1181))

* Fix `eth_getBalance` result with refund gas in fast-query mode ([\#1187](https://github.com/okex/exchain/pull/1187))


### Code refactoring

* Refine `abci` tracer ([\#1160](https://github.com/okex/exchain/pull/1160))


### Performance Improvements

* Optimize holding active objects list in `state objects`  ([\#1178](https://github.com/okex/exchain/pull/1178))

* Disable `trace.GoRId` for better performance ([\#1184](https://github.com/okex/exchain/pull/1184))

* Disable `gid` by default ([\#1190](https://github.com/okex/exchain/pull/1190))


### Chores

* Add CI ut ([\#1162](https://github.com/okex/exchain/pull/1162))



## [0.19.13](https://github.com/okex/exchain/compare/v0.19.12...v0.19.13) (2021-11-01)


### Features

* Support paralleled tx in `x/evm` ([\#1100](https://github.com/okex/exchain/pull/1100))

* Support paralleled-tx when replay tx ([\#1112](https://github.com/okex/exchain/pull/1112))

* Add keccak256Hash Cache to storage `address-key` hash ([\#1120](https://github.com/okex/exchain/pull/1120))

* Add flag to close checkTx mutex ([\#1129](https://github.com/okex/exchain/pull/1129))

* Hide `stream` flags ([\#1136](https://github.com/okex/exchain/pull/1136))

* Add pprof to `applyblock` ([\#1147](https://github.com/okex/exchain/pull/1147))


### Bug fixes

* Bloom must not support paralleled-tx in `x/evm` ([\#1119](https://github.com/okex/exchain/pull/1119))

* Fix `leveldb` api compatible in `x/evm` ([\#1122](https://github.com/okex/exchain/pull/1122))

* Fix `saveCommitOrphans` bug when it is executed before commit event ([\#1131](https://github.com/okex/exchain/pull/1131))

* Fix concurrent read write map ([\#1155](https://github.com/okex/exchain/pull/1155))


### Code refactoring

* Sub-command `covert`  covert data between`goleveldb` and `rocksdb` ([\#1109](https://github.com/okex/exchain/pull/1109))

* Refactor `addressRecord` in `mempool`  ([\#1149](https://github.com/okex/exchain/pull/1149))

* Refine timing in `x/analyzer` ([\#1151](https://github.com/okex/exchain/pull/1151))


### Tests

* Add ut of async-commit ([\#1137](https://github.com/okex/exchain/pull/1137))


### Chores

* Change path about `cosmos` and `tendermint` in makefile ([\#1142](https://github.com/okex/exchain/pull/1142))


### Breaking Changes

* Ship all dependence, including `cosmos-sdk`, `tendermint`, `iavl`, `tendermint-db` ([\#1128](https://github.com/okex/exchain/pull/1128))



## [0.19.12](https://github.com/okex/exchain/compare/v0.19.11...v0.19.12) (2021-10-18)


### Features

* Add `antehandle` analysis ([\#1079](https://github.com/okex/exchain/pull/1079))

* Add command to export eth keystore file in `exchaincli` ([\#1084](https://github.com/okex/exchain/pull/1084))

* Dump pprof automatically when ApplyBlock elapsed time is too long. ([\#1087](https://github.com/okex/exchain/pull/1087))

* Run replay with dump pprof ([\#1089](https://github.com/okex/exchain/pull/1089))

* Use flag control log print ([\#1098](https://github.com/okex/exchain/pull/1098))

* Add sub-command `covert` to covert data from `goleveldb` to `rocksdb` ([\#1114](https://github.com/okex/exchain/pull/1114))


### Bug fixes

* `x/evm` start `x/analyzer` log when transaction is `checktx` ([\#1091](https://github.com/okex/exchain/pull/1091))

* Fix eth tx multiple signature ([\#1092](https://github.com/okex/exchain/pull/1092))


### Code refactoring

* repair state with `start-height` instead of `commit-inyerval` ([\#1095](https://github.com/okex/exchain/pull/1095))


### Performance Improvements

* Compact `rocksdb` ([\#1083](https://github.com/okex/exchain/pull/1083))

* `watchdb` is compatible with `rocksdb` ([\#1088](https://github.com/okex/exchain/pull/1088))


## [0.19.11](https://github.com/okex/exchain/compare/v0.19.10...v0.19.11) (2021-10-14)


### Features

* Add more timestamp log into `x/analyzer` and count times  that `db`  is written or read ([\#1060](https://github.com/okex/exchain/pull/1060))

* Add features check whether address is blocked ([\#1066](https://github.com/okex/exchain/pull/1066))

* Count `evm` execute time ([\#1072](https://github.com/okex/exchain/pull/1072))


### Bug fixes

* Fix wrong gas consume in `antehandler` ([\#1085](https://github.com/okex/exchain/pull/1085))

* `x/analyzer` start log when transaction is `checktx` ([\#1093](https://github.com/okex/exchain/pull/1093))


### Performance Improvements

* Delete useless code in deliverTx ([\#1076](https://github.com/okex/exchain/pull/1076))


### Chores

* Support make `rocksdb` in makefile ([\#1074](https://github.com/okex/exchain/pull/1074))



## [0.19.10](https://github.com/okex/exchain/compare/v0.19.9...v0.19.10) (2021-10-08)


### Features

* Support `rocksdb` ([\#1055](https://github.com/okex/exchain/pull/1055))

* Add trace log to catch problems when happen enhancement ([\#1057](https://github.com/okex/exchain/pull/1057))

* Add `DEBUG` argument to makefile to switch between debug and release mode ([\#1058](https://github.com/okex/exchain/pull/1058))

* Add more timestamp log into `x/analyzer` ([\#1059](https://github.com/okex/exchain/pull/1059))



## [0.19.9](https://github.com/okex/exchain/compare/v0.19.8...v0.19.9) (2021-09-30)


### Features

* Support async commit to `iavl tree` ([\#1048](https://github.com/okex/exchain/pull/1048))



## [0.19.8](https://github.com/okex/exchain/compare/v0.19.7...v0.19.8) (2021-09-29)


### Features

* Add log into `x/evm commitStatedb` and `x/analyzer` ([\#1041](https://github.com/okex/exchain/pull/1041))


### Bug fixes

* Fix current local replay ([\#1043](https://github.com/okex/exchain/pull/1043))



## [0.19.7](https://github.com/okex/exchain/compare/v0.19.6...v0.19.7) (2021-09-27)


### Features

* Add query blocked contracts in rest api ([\#1027](https://github.com/okex/exchain/pull/1027))

* Add `iaviewer` ([\#1030](https://github.com/okex/exchain/pull/1030))


### Bug fixes 

* Fix RPC API `eth_getCode` occur error after get `blocknumber` ([\#1029](https://github.com/okex/exchain/pull/1029))


### Chores

* Enable `x/evm` in `start.sh`  ([\#1040](https://github.com/okex/exchain/pull/1040))



## [0.19.6](https://github.com/okex/exchain/compare/v0.19.5...v0.19.6) (2021-09-16)


### Features

* Support `eip-1898` feature ([\#1024](https://github.com/okex/exchain/pull/1024))


### Bug fixes

* Fix sending cosmos tx to pending pool ([\#1023](https://github.com/okex/exchain/pull/1023))

* Fix local replay ([\#1061](https://github.com/okex/exchain/pull/1061))


### Code Refactoring

* Rename sub-command `repair-data` to `repair-state` ([\#1020](https://github.com/okex/exchain/pull/1020))


### Performance Improvements

* Auto dump pprof when cpu is high ([\#1018](https://github.com/okex/exchain/pull/1018))


### Chores

* Enable `golangci-lint` ([\#1019](https://github.com/okex/exchain/pull/1019))



## [0.19.5](https://github.com/okex/exchain/compare/v0.19.2...v0.19.5) (2021-09-09)


### Features

* Add pending pool in `txpool` ([\#997](https://github.com/okex/exchain/pull/997))

* Add sub-command to repair state ([\#1008](https://github.com/okex/exchain/pull/1008))


### Bug Fixes

* Fix `txpool` doesn't drop tx when broadcast error is `ErrInvalidSequence` ([\#995](https://github.com/okex/exchain/pull/995))

* Fix replay block panic when `haltheight` is not set ([\#1013](https://github.com/okex/exchain/pull/1013))


### Chores

* Produce testnet script ([\#1003](https://github.com/okex/exchain/pull/1003))



## [0.19.2](https://github.com/okex/exchain/compare/v0.19.1...v0.19.2) (2021-08-31)


### Features

* Support dynamic config ([\#982](https://github.com/okex/exchain/pull/982))

* Add `eth_multiCall` RPC-API to perform multiple raw contract call multiple raw contract call ([\#998](https://github.com/okex/exchain/pull/998))


### Performance Improvements

*  Optimize compact in pruning ([\#993](https://github.com/okex/exchain/pull/993))


### Chores

* Produce snapshot script ([\#996](https://github.com/okex/exchain/pull/996))



## [0.19.1](https://github.com/okex/exchain/compare/v0.18.18...v0.19.1) (2021-08-24)

### Features

* Support 0x prefixed address format ([\#973](https://github.com/okex/exchain/pull/973))

* Add logger for `txpool` ([\#983](https://github.com/okex/exchain/pull/983))


### Documentation

* Provide issue template to  report running error ([\#979](https://github.com/okex/exchain/pull/979))


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

* Expand `iavl` cache size ([\#960](https://github.com/okex/exchain/pull/960))


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

* Enhance websocket handle error ([\#838](https://github.com/okex/exchain/pull/838))



## [0.17.5](https://github.com/okex/exchain/compare/v0.17.4...v0.17.5) (2021-04-22)


### Build

* Update cosmos version to `v0.39.2-exchain2` and tendermint version to `v0.33.9-exchain1` ([\#837](https://github.com/okex/exchain/pull/837))



## [0.17.4](https://github.com/okex/exchain/compare/v0.17.3...v0.17.4) (2021-04-20)


### Code Refactoring

* Set the default gas price of `eth_call` to `1 Gwei` ([\#825](https://github.com/okex/exchain/pull/825))

* Remove redundant code ([\#826](https://github.com/okex/exchain/pull/826))



## [0.17.3](https://github.com/okex/exchain/compare/v0.17.2...v0.17.3) (2021-04-15)


### Performance Improvements

* Optimize error of `ErrTxTooLarge` ([\#820](https://github.com/okex/exchain/pull/820))

* Max world state num ([\#819](https://github.com/okex/exchain/pull/819))



## [0.17.2](https://github.com/okex/exchain/compare/v0.17.1...v0.17.2) (2021-04-13)


### Features

* The websocket interface supports the query of the web3 library ([\#792](https://github.com/okex/exchain/pull/792))

* Support the function of calling web3 api in websocket ([\#795](https://github.com/okex/exchain/pull/795))


### Code Refactoring

* Change nonce format ([\#785](https://github.com/okex/exchain/pull/785))


### Performance Improvements

* Optimize performance of the function `eth_call` ([\#793](https://github.com/okex/exchain/pull/793))



## [0.17.1](https://github.com/okex/exchain/compare/v0.17.0...v0.17.1) (2021-04-12)


### BREAKING CHANGES

* Rename to `exchain` ([\#816](https://github.com/okex/exchain/pull/816))



# [0.17.0](https://github.com/okex/exchain/compare/v0.16.9...v0.17.0) (2021-04-11)


### Features

* Add pruning for block state db and add export from app db ([\#811](https://github.com/okex/exchain/pull/811))


### Build

* Update cosmos-sdk version to `v0.39.2-exchain1` ([\#814](https://github.com/okex/exchain/pull/814))


