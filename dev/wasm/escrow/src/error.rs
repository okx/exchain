use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized")]
    Unauthorized {},

    #[error("Escrow expired (end_height {end_height:?} end_time {end_time:?})")]
    Expired {
        end_height: Option<u64>,
        end_time: Option<u64>,
    },

    #[error("Escrow not expired")]
    NotExpired {},
}
