use cosmwasm_schema::{cw_serde, QueryResponses};
use super::state::Number;

#[cw_serde]
pub struct InstantiateMsg {
    pub count: Number,
}

#[cw_serde]
pub enum ExecuteMsg {
    Increment { count: Number },
    Write { count: Number },
    Read { count: Number },
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    // GetCount returns the current count as a json-encoded number
    #[returns(GetCountResponse)]
    GetCount {},
}

// We define a custom struct for each query response
#[cw_serde]
pub struct GetCountResponse {
    pub count: Number,
}
