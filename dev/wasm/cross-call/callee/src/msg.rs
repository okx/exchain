use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::Uint256;

#[cw_serde]
pub struct InstantiateMsg {
}

#[cw_serde]
// pub enum ExecuteMsg {
//     /// Releasing all funds in the contract to the beneficiary. This is the only "proper" action of this demo contract.
//     Release {},
// }
pub enum ExecuteMsg {
    Add {delta:Uint256},
    Subtract {},
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    #[returns(Uint256)]
    GetCounter {},
}

#[cw_serde]
pub struct VerifierResponse {
    pub verifier: String,
}

#[cw_serde]
pub struct MigrateMsg {}
