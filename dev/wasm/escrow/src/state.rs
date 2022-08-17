use cosmwasm_std::{Addr, Env, Storage};
use cosmwasm_storage::{singleton, singleton_read, ReadonlySingleton, Singleton};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

static CONFIG_KEY: &[u8] = b"config";

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct State {
    pub arbiter: Addr,
    pub recipient: Addr,
    pub source: Addr,
    pub end_height: Option<u64>,
    pub end_time: Option<u64>,
}

impl State {
    pub fn is_expired(&self, env: &Env) -> bool {
        if let Some(end_height) = self.end_height {
            if env.block.height > end_height {
                return true;
            }
        }

        if let Some(end_time) = self.end_time {
            if env.block.time.nanos() > end_time * 1000 {
                return true;
            }
        }
        false
    }
}

pub fn config(storage: &mut dyn Storage) -> Singleton<State> {
    singleton(storage, CONFIG_KEY)
}

pub fn config_read(storage: &dyn Storage) -> ReadonlySingleton<State> {
    singleton_read(storage, CONFIG_KEY)
}
