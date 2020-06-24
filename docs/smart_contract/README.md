
* [The Guidelines for deploy an existing Smart Contract and interact with it afterwards](#The-Guidelines-for-deploy-an-existing-Smart-Contract-and-interact-with-it-afterwards)
   * [1、Prepare](#1prepare)
      * [1.1、Set up okchaincli](#11set-up-okchaincli)
         * [1.1.1、Build okchaincli](#111build-okchaincli)
         * [1.1.2、Set up okchaincli env](#112set-up-okchaincli-env)
      * [1.2、Prepare okchain account](#12prepare-okchain-account)
        * [1.2.1、Create okchain accounts](#121create-okchain-accounts)
         * [1.2.2、Require Test Tokens](#122require-test-tokens)
      * [1.3、 Prepare wasm contract file](#13-prepare-wasm-contract-file)
   * [2、Install Contract](#2install-contract)
   * [3、Instantiate Contract](#3instantiate-contract)
   * [4、Invoke Contract](#4invoke-contract)
   * [5、Query Contract](#5query-contract)
* [Creating your own Smart Contract](#Creating-your-own-Smart-Contract)
   * [1、Implementing the Smart Contract](#1implementing-the-smart-contract)
   * [2、Testing the Smart Contract (rust)](#2testing-the-smart-contract-rust)
   * [3、Production Builds](#3production-builds)


# The Guidelines for deploy an existing Smart Contract and interact with it afterwards

A smart contract is a computer program or a transaction protocol which is intended to automatically execute, control or document legally relevant events and actions according to the terms of a contract or an agreement. To better describe the usage of smart contracts on okchain , we will use the erc20 example to show the whole process.

## 1、Prepare

### 1.1、Set up okchaincli

As cosmwasm is under repaid development and have't release a stable version, so we recommend to use our official okchain-wasm test-net to run your smart contract.

Here is the instructions to help you connect to our okchain-wasm test-net.

#### 1.1.1、Build okchaincli

Because we don't have a stable cosmwasm version yet, so we decide not to provide a binary download now, you will need to build the okchaincli by yourself. Before you run the next commands, make sure you have install a right version of golang and configure it rightly (we recommend a version of go1.12.7, as for how to configure golang, please google it and there will be a lots of results)

here is the instructions to build okchaincli:

~~~bash
git clone https://github.com/okex/okchain.git -b okchain-wasm
cd okchain
make install
~~~

#### 1.1.2、Set up okchaincli env

Before you  connect to our okchain-wasm test-net and play with it, we recommend you to configure your okchaincli, it will save you lots of time to type configuration command.

~~~bash
okchaincli config chain-id okchain 
okchaincli config output json    
okchaincli config indent true
okchaincli config node tcp://3.112.102.224:26657
okchaincli config trust-node true
~~~

### 1.2、Prepare okchain account

So we are ready to play with okchain-wasm test-net, before do that,  you will need some accounts to represent yourself to Interact with others.

#### 1.2.1、Create okchain accounts

If you don't have any okchain accounts in okchain-wasm test-net yet, you can follow these command to add by yourself:

~~~bash
$okchaincli keys add <name> [flags]
~~~

usage:

~~~bash
$okchaincli keys add yourAccountName
~~~

The next command will show your accounts list you have created:

~~~bash
$okchaincli keys list
[
  {
    "name": "account1",
    "type": "local",
    "address": "okchain1gsn3jf86x253z4990tf8hpsy6cqk9rxk5tll0y",
    "pubkey": "okchainpub1addwnpepqv82z3zyw2rt897ed5kc0r8v0v3ul4qll3usx35fsp4ld4peslxru753cq0"
  },
  {
    "name": "account2",
    "type": "local",
    "address": "okchain1jt8kk0jyvdnvzmfxgrdm40smv3p752hw6qn8ay",
    "pubkey": "okchainpub1addwnpepqvry6upv856cg65h6nk0zneezmcphxzvu85gw9vvxsxy409hdyg8cewu9cu"
  }
]
~~~

#### 1.2.2、Require Test Tokens

Before any operation, make sure the accounts list above have enough token in okchain-wasm test-net. 

You can contract us to ask for test tokens:

Wechat: 

Telegram:  `Okchain`

**Warning**: The test token here is different from the **Okchain test-net**, the token here is just for **Okchain-wasm test-net**!

### 1.3、 Prepare wasm contract file

One of the most significant tokens is known as ERC-20, which has emerged as the technical standard used for all smart contracts on the Ethereum blockchain for token implementation. As it is so famous in blockchain industry, so we decided to provide an erc-20 implementation to show the powerful of wasm smart contract. You can download the erc20 contract wasm file: [contract.wasm](https://raw.githubusercontent.com/CosmWasm/cosmwasm-examples/master/erc20/contract.wasm), the rust code can be found in this link: [erc20](https://github.com/CosmWasm/cosmwasm-examples/tree/master/erc20).

The contract provide the following function: `init`, `Approve`, `Transfer`, `TransferFrom`, `Burn`, `Balance`, `Allowance`, which is the main function of the erc20 protocol. The following instructions will show you how to use it in our okchain-wasm test-net.

## 2、Install Contract

~~~bash
$okchaincli tx wasm store -h
Upload a wasm binary

Usage:
  okchaincli tx wasm store [wasm file] --source [source] --builder [builder] [flags]
~~~

Usage:

~~~bash
okchaincli tx wasm store  ./erc20.wasm --from captain --gas 2000000 --fees 2okt -y -b block
~~~

Get `code_id` from output:

~~~bash
"code_id": 1
~~~

## 3、Instantiate Contract

~~~bash
$okchaincli tx wasm instantiate -h
Instantiate a wasm contract

Usage:
  okchaincli tx wasm instantiate [code_id_int64] [json_encoded_init_args] [flags]
~~~

Usage:

~~~bash
okchaincli tx wasm instantiate 1 "{\"name\":\"Test okchain token\",\"symbol\":\"TOKT\",\"decimals\":10,\"initial_balances\":[{\"address\":\"okchain1gsn3jf86x253z4990tf8hpsy6cqk9rxk5tll0y\",\"amount\":\"10000000\"}]}"  --from account1 --label "First ERC20 token in okchain" --gas 200000 --fees 4okt -y -b block
~~~

Get the `contract address` from output：

~~~bash
"contract_address": okchain18vd8fpwxzck93qlwghaj6arh4p7c5n897czf0h
~~~

## 4、Invoke Contract

~~~bash
$okchaincli tx wasm execute -h
Execute a command on a wasm contract

Usage:
  okchaincli tx wasm execute [contract_addr_bech32] [json_encoded_send_args] [flags]
~~~

Usage:

~~~bash
// Transfer from account1 to account2
okchaincli tx wasm execute okchain18vd8fpwxzck93qlwghaj6arh4p7c5n897czf0h "{\"transfer\":{\"recipient\":\"okchain1jt8kk0jyvdnvzmfxgrdm40smv3p752hw6qn8ay\",\"amount\":\"500000\"}}" --from account1 --gas 2000000 --fees 4okt -y -b block
~~~

## 5、Query Contract

* Query contract detail

~~~bash
$okchaincli query wasm contract -h
Prints out metadata of a contract given its address

Usage:
  okchaincli query wasm contract [bech32_address] [flags]
~~~

Usage:

~~~bash
$okchaincli query wasm contract okchain18vd8fpwxzck93qlwghaj6arh4p7c5n897czf0h	
~~~

* Query account balance

~~~bash
$okchaincli query wasm contract-state -h
Querying commands for the wasm module

Usage:
  okchaincli query wasm contract-state [flags]
  okchaincli query wasm contract-state [command]

Available Commands:
  all         Prints out all internal state of a contract given its address
  raw         Prints out internal state for key of a contract given its address
  smart       Calls contract with given address  with query data and prints the returned result
~~~

Usage:

~~~bash
// query account1 balance
okchaincli query wasm contract-state  smart okchain18vd8fpwxzck93qlwghaj6arh4p7c5n897czf0h "{\"balance\":{\"address\":\"okchain1gsn3jf86x253z4990tf8hpsy6cqk9rxk5tll0y\"}}" 

// query account2 balance
okchaincli query wasm contract-state  smart okchain18vd8fpwxzck93qlwghaj6arh4p7c5n897czf0h "{\"balance\":{\"address\":\"okchain1jt8kk0jyvdnvzmfxgrdm40smv3p752hw6qn8ay\"}}"
~~~


# Creating your own Smart Contract

If you want to get started building you own, the simplest way is to go to the [cosmwasm-template](https://github.com/CosmWasm/cosmwasm-template) repository and follow the instructions. This will give you a simple contract along with tests, and a properly configured build environment. From there you can edit the code to add your desired logic and publish it as an independent repo.

## 1、Implementing the Smart Contract

If you start from the [cosmwasm-template](https://github.com/CosmWasm/cosmwasm-template), you may notice that all of the Wasm exports are taken care of by `lib.rs`, which should shouldn't need to modify. What you need to do is simply look in `contract.rs` and implement `init` ,`handle` and `query` functions, defining your custom `InitMsg` , `HandleMsg` and `QueryMsg` structs for parsing your custom message types (as json):

~~~rust
pub fn init<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
    msg: InitMsg,
) -> StdResult<InitResponse> {}

pub fn handle<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
    msg: HandleMsg,
) -> StdResult<HandleResponse> {}

pub fn query<S: Storage, A: Api, Q: Querier>(
    deps: &Extern<S, A, Q>,
    msg: QueryMsg,
) -> StdResult<Binary> {}
~~~

## 2、Testing the Smart Contract (rust)

For quick unit tests and useful error messages, it is often helpful to compile the code using native build system and then test all code except for the `extern "C"` functions (which should just be small wrappers around the real logic).

If you have non-trivial logic in the contract, please write tests using rust's standard tooling. If you run `cargo test`, it will compile into native code using the `debug` profile, and you get the normal test environment you know and love. Notably, you can add plenty of requirements to `[dev-dependencies]` in `Cargo.toml` and they will be available for your testing joy. As long as they are only used in `#[cfg(test)]` blocks, they will never make it into the (release) Wasm builds and have no overhead on the production artifact.

## 3、Production Builds

The above build process (`cargo wasm`) works well to produce wasm output for testing. However, it is quite large, around 1.5 MB likely, and not suitable for posting to the blockchain. Furthermore, it is very helpful if we have reproducible build step so others can prove the on-chain wasm code was generated from the published rust code.

For that, we have a separate repo, [cosmwasm-opt](https://github.com/CosmWasm/cosmwasm-opt) that provides a [docker image](https://hub.docker.com/r/CosmWasm/cosmwasm-opt/tags) for building. For more info, look at [cosmwasm-opt README](https://github.com/CosmWasm/cosmwasm-opt/blob/master/README.md#usage), but the quickstart guide is:

~~~bash
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/code/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/rust-optimizer:0.8.0
~~~

It will output a highly size-optimized build as `contract.wasm` in `$CODE`. With our example contract, the size went down to 126kB (from 1.6MB from `cargo wasm`). If we didn't use serde-json, this would be much smaller still...
