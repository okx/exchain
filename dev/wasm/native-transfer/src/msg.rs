use cosmwasm_schema::{cw_serde};

#[cw_serde]
pub struct InstantiateMsg {}

#[cw_serde]
pub enum ExecuteMsg {
    /// Transfer is a base message to move tokens to another account without triggering actions
    Transfer { recipient: String },
}

// #[cw_serde]
// #[derive(QueryResponses)]
// pub enum QueryMsg {
//     // GetCount returns the current count as a json-encoded number
//     #[returns(GetCountResponse)]
//     GetCount {},
// }

// We define a custom struct for each query response
// #[cw_serde]
// pub struct GetCountResponse {
//     pub count: i32,
// }
