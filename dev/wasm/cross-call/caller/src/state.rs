use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::{Addr, Uint256};

pub const CONFIG_KEY: &[u8] = b"config";

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct State {
    pub verifier: Addr,
    pub beneficiary: Addr,
    pub funder: Addr,
}

pub const CONFIG_KEY1: &[u8] = b"counter";
pub const CONFIG_KEY2: &[u8] = b"counter1";
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct State1 {
    pub counter: Uint256,
}