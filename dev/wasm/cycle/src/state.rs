use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::Uint256;
use cw_storage_plus::Item;

pub type Number = u128;

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct State {
    pub count: Number,
}

pub const STATE: Item<State> = Item::new("state");
