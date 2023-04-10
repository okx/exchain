use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::{Addr, Uint128};
use cw_storage_plus::Item;

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct State {
    pub count: Uint128,
    pub owner: Addr,
}

pub const STATE: Item<State> = Item::new("state");
