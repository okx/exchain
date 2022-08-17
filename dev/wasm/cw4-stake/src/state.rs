use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::{Addr, Uint128};
use cw20::Denom;
use cw4::TOTAL_KEY;
use cw_controllers::{Admin, Claims, Hooks};
use cw_storage_plus::{Item, Map, SnapshotMap, Strategy};
use cw_utils::Duration;

pub const CLAIMS: Claims = Claims::new("claims");

#[derive(Serialize, Deserialize, Clone, PartialEq, JsonSchema, Debug)]
pub struct Config {
    /// denom of the token to stake
    pub denom: Denom,
    pub tokens_per_weight: Uint128,
    pub min_bond: Uint128,
    pub unbonding_period: Duration,
}

pub const ADMIN: Admin = Admin::new("admin");
pub const HOOKS: Hooks = Hooks::new("cw4-hooks");
pub const CONFIG: Item<Config> = Item::new("config");
pub const TOTAL: Item<u64> = Item::new(TOTAL_KEY);

pub const MEMBERS: SnapshotMap<&Addr, u64> = SnapshotMap::new(
    cw4::MEMBERS_KEY,
    cw4::MEMBERS_CHECKPOINTS,
    cw4::MEMBERS_CHANGELOG,
    Strategy::EveryBlock,
);

pub const STAKE: Map<&Addr, Uint128> = Map::new("stake");
