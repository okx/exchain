use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use cosmwasm_std::{CosmosMsg, CustomMsg,Uint128};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct InstantiateMsg {
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum ExecuteMsg {
    Add {delta:Uint128},
    AddCounterForEvm {evm_contract:String,delta:Uint128}
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum QueryMsg {
    GetCounter {}
}

#[derive(Serialize, Deserialize, Clone, PartialEq, JsonSchema, Debug)]
#[serde(rename_all = "snake_case")]
pub struct CallToEvmMsg {
    pub sender: String, 
    pub evmaddr: String,
    pub calldata: String,
    pub value: Uint128,
}

#[derive(Serialize, Deserialize, Clone, PartialEq, JsonSchema, Debug)]
#[serde(rename_all = "snake_case")]
pub struct MigrateMsg{
    pub sender: String

}

impl Into<CosmosMsg<CallToEvmMsg>> for CallToEvmMsg {
    fn into(self) -> CosmosMsg<CallToEvmMsg> {
        CosmosMsg::Custom(self)
    }
}

impl CustomMsg for CallToEvmMsg {}

pub struct CallToEvmMsgResponse {
    pub response: String,
}


