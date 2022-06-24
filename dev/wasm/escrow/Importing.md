# Importing

In [Publishing](./Publishing.md), we discussed how you can publish your contract to the world.
This looks at the flip-side, how can you use someone else's contract (which is the same
question as how they will use your contract). Let's go through the various stages.

## Getting the Code

Before using remote code, you most certainly want to verify it is honest.
There are two ways to get the code of another contract, either by cloning the git repo
or by downloading the cargo crate. You should be familiar with using git already.
However, the rust publishing system doesn't rely on git tags (they are optional),
so to make sure you are looking at the proper code, I would suggest getting the
actual code of the tagged crate.

```sh
cargo install cargo-download
cargo download cw-escrow==0.1.0 > crate.tar.gz
tar xzvf crate.tar.gz
cd cw-escrow-0.1.0
```

(alternate, simpler approach, but seems to be broken):

```sh
cargo install cargo-clone
cargo clone cw-escrow --vers 0.1.0
```

## Verifying Artifacts

The simplest audit of the repo is to simply check that the artifacts in the repo
are correct. You can use the same commands you do when developing, with the one
exception that the `.cargo/config` file is not present on downloaded crates,
so you will have to run the full commands.

First, make a git commit here, so we can quickly see any diffs:

```sh
git init .
echo target > .gitignore
git add .
git commit -m 'From crates.io'
```

To validate the tests:

```sh
cargo build --release --target wasm32-unknown-unknown
cargo test
```

To generate the schema:

```sh
cargo run --example schema
```

And to generate the `contract.wasm` and `hash.txt`:

```sh
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/code/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/rust-optimizer:0.9.0
```

Make sure the values you generate match what was uploaded with a simple `git diff`.
If there is any discrepancy, please raise an issue on the repo, and please add an issue
to the cawesome-wasm list if the package is listed there (it should be validated before
adding, but just in case).

In the future, we will produce a script to do this automatic verification steps that can
be run by many individuals to quickly catch any fake uploaded wasm hashes in a
decentralized manner.

## Reviewing

Once you have done the quick programatic checks, it is good to give at least a quick
look through the code. A glance at `examples/schema.rs` to make sure it is outputing
all relevant structs from `contract.rs`, and also ensure `src/lib.rs` is just the
default wrapper (nothing funny going on there). After this point, we can dive into
the contract code itself. Check the flows for the handle methods, any invariants and
permission checks that should be there, and a reasonable data storage format.

You can dig into the contract as far as you want, but it is important to make sure there
are no obvious backdoors at least.

## Decentralized Verification

It's not very practical to do a deep code review on every dependency you want to use,
which is a big reason for the popularity of code audits in the blockchain world. We trust
some experts review in lieu of doing the work ourselves. But wouldn't it be nice to do this
in a decentralized manner and peer-review each other's contracts? Bringing in deeper domain
knowledge and saving fees.

Luckily, there is an amazing project called [crev](https://github.com/crev-dev/cargo-crev/blob/master/cargo-crev/README.md)
that provides `A cryptographically verifiable code review system for the cargo (Rust) package manager`.

I highly recommend that CosmWasm contract developers get set up with this. At minimum, we
can all add a review on a package that programmatically checked out that the json schemas
and wasm bytecode do match the code, and publish our claim, so we don't all rely on some
central server to say it validated this. As we go on, we can add deeper reviews on standard
packages.

If you want to use `cargo-crev`, please follow their
[getting started guide](https://github.com/crev-dev/cargo-crev/blob/master/cargo-crev/src/doc/getting_started.md)
and once you have made your own *proof repository* with at least one *trust proof*,
please make a PR to the [`cawesome-wasm`]() repo with a link to your repo and
some public name or pseudonym that people know you by. This allows people who trust you
to also reuse your proofs.

There is a [standard list of proof repos](https://github.com/crev-dev/cargo-crev/wiki/List-of-Proof-Repositories)
with some strong rust developers in there. This may cover dependencies like `serde` and `snafu`
but will not hit any CosmWasm-related modules, so we look to bootstrap a very focused
review community.
