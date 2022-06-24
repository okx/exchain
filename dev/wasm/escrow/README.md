# Escrow

This is a simple single-use escrow contract. It creates a contract that can hold some
native tokens and gives the power to an arbiter to release them to a pre-defined
beneficiary. They can release all tokens, or only a fraction. If an optional
timeout is reached, the tokens can no longer be released, rather they can only
be returned to the original funder. Tokens can be added to the contract at any
time without causing any errors, or losing access to them.

This contract is mainly considered as a simple tutorial example. In the real
world, you would probably want one contract to manage many escrows and allow
some global configuration options on it. It is generally simpler to rely on
some well-known address for handling all escrows securely than checking each
deployed escrow is using the proper wasm code.

As of v0.2.0, this was rebuilt from
[`cosmwasm-template`](https://github.com/confio/cosmwasm-template),
which is the recommended way to create any contracts.

## Using this project

If you want to get acquainted more with this contract, you should check out
[Developing](./Developing.md), which explains more on how to run tests and develop code.
[Publishing](./Publishing.md) contains useful information on how to publish your contract
to the world, once you are ready to deploy it on a running blockchain. And
[Importing](./Importing.md) contains information about pulling in other contracts or crates
that have been published.

But more than anything, there is an [online tutorial](https://www.cosmwasm.com/docs/getting-started/intro),
which leads you step-by-step on how to modify this particular contract.
