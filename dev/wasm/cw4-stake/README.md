# CW4 Stake

This is a second implementation of the [cw4 spec](../../packages/cw4/README.md).
It fulfills all elements of the spec, including the raw query lookups,
and is designed to be used as a backing storage for 
[cw3 compliant contracts](../../packages/cw3/README.md).

It provides a similar API to [`cw4-group`] (which handles elected membership),
but rather than appointing members (by admin or multisig), their
membership and weight are based on the number of tokens they have staked.
This is similar to many DAOs.

Only one denom can be bonded with both `min_bond` as the minimum amount
that must be sent by one address to enter, as well as `tokens_per_weight`,
which can be used to normalize the weight (eg. if the token is uatom
and you want 1 weight per ATOM, you can set `tokens_per_weight = 1_000_000`).

There is also an unbonding period (`Duration`) which sets how long the
tokens are frozen before being released. These frozen tokens can neither
be used for voting, nor claimed by the original owner. Only after the period
can you get your tokens back. This liquidity loss is the "skin in the game"
provided by staking to this contract.

## Instantiation

**TODO**

To create it, you must pass in a list of members, as well as an optional
`admin`, if you wish it to be mutable.

```rust
pub struct InstantiateMsg {
    /// denom of the token to stake
    pub stake: String,
    pub tokens_per_weight: u64,
    pub min_bond: Uint128,
    pub unbonding_period: Duration,
}
```

Members are defined by an address and a weight. This is transformed
and stored under their `CanonicalAddr`, in a format defined in
[cw4 raw queries](../../packages/cw4/README.md#raw).

Note that 0 *is an allowed weight*. This doesn't give any voting rights, 
but it does define this address is part of the group, which may be
meaningful in some circumstances.

The weights of the members will be computed as the funds they send 
(in tokens) divided by `tokens_per_weight`, rounded down to the nearest
whole number (i.e. using integer division). If the total sent is less than
`min_bond`, the stake will remain, but they will not be counted as a
member. If `min_bond` is higher than `tokens_per_weight`, you cannot
have any member with 0 weight.

## Messages

Most messages and queries are defined by the 
[cw4 spec](../../packages/cw4/README.md). Please refer to it for more info.

The following messages have been added to handle un/staking tokens:

`Bond{}` - bond all staking tokens sent with the message and update membership weight

`Unbond{tokens}` - starts the unbonding process for the given number 
  of tokens. The sender immediately loses weight from these tokens,
  and can claim them back to his wallet after `unbonding_period`

`Claim{}` -  used to claim your native tokens that you previously "unbonded"
  after the contract-defined waiting period (eg. 1 week)

And the corresponding queries:

`Claims{address}` - Claims shows the tokens in process of unbonding
    for this address

`Staked{address}` - Show the number of tokens currently staked by this address.
