use cosmwasm_std::{entry_point, from_slice, to_binary, AllBalanceResponse, BankMsg, Deps, DepsMut, Env, Event, MessageInfo, QueryResponse, Response, StdError, StdResult, Binary, Uint256, WasmMsg, coin, to_vec};

use crate::errors::HackError;
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg, VerifierResponse, MigrateMsg};
use crate::state::{State, CONFIG_KEY, CONFIG_KEY1, State1};

#[entry_point]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, HackError> {
    // This adds some unrelated event attribute for testing purposes
    let d = Uint256::from(0u32);
    deps.storage.set(
        CONFIG_KEY1,
        &to_vec(&State1 {
            counter: d,
        })?,
    );
    Ok(Response::new())
}

#[entry_point]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, HackError> {
    match msg {
        ExecuteMsg::Add { delta } => try_add(deps,delta),
        ExecuteMsg::Subtract {} => try_sub(deps),
    }
}

pub fn try_add(deps: DepsMut,delta:Uint256) -> Result<Response, HackError> {
    let data = deps
        .storage
        .get(CONFIG_KEY1)
        .ok_or_else(|| StdError::not_found("State1"))?;
    let state: State1 = from_slice(&data)?;
    let mut counter= Uint256::from(state.counter);

    counter += delta;
    deps.storage.set(
        CONFIG_KEY1,
        &to_vec(&State1 {
            counter: counter,
        })?,
    );

    Ok(Response::new().add_attribute("Added", counter).add_attribute("Changed", counter))
}

pub fn try_sub(deps: DepsMut) -> Result<Response, HackError> {
    let data = deps
        .storage
        .get(CONFIG_KEY1)
        .ok_or_else(|| StdError::not_found("State1"))?;
    let state: State1 = from_slice(&data)?;
    let mut counter= Uint256::from(state.counter);
    if counter == Uint256::zero(){
        HackError::TooLow {};
    }
    counter -= Uint256::from(1u128);
    deps.storage.set(
        CONFIG_KEY1,
        &to_vec(&State1 {
            counter: counter,
        })?,
    );
    Ok(Response::new().add_attribute("Changed", counter))
}

#[entry_point]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<QueryResponse> {
    match msg {
        QueryMsg::GetCounter {} => to_binary(&query_counter(deps)?),
    }
}

fn query_counter(deps: Deps) -> StdResult<Uint256> {
    let data = deps
        .storage
        .get(CONFIG_KEY1)
        .ok_or_else(|| StdError::not_found("State1"))?;
    let state: State1 = from_slice(&data)?;
    Ok(state.counter)
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn migrate(deps: DepsMut, _env: Env, _msg: MigrateMsg) -> StdResult<Response> {
    Ok(Response::default())
}


#[cfg(test)]
mod tests {
    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_dependencies_with_balances, mock_env, mock_info, MOCK_CONTRACT_ADDR, mock_dependencies_with_balance};
    use cosmwasm_std::{Api as _, from_binary};
    // import trait Storage to get access to read
    use cosmwasm_std::{attr, coins, Addr, Storage, SubMsg};

    #[test]
    fn proper_initialization() {
        let mut deps = mock_dependencies();

        let creator = String::from("creator");
        let expected_state = State1 {
            counter: Uint256::from(0u32),
        };

        let msg = InstantiateMsg {};
        let info = mock_info(creator.as_str(), &[]);
        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());

        // it worked, let's check the state
        let data = deps.storage.get(CONFIG_KEY1).expect("no data stored");
        let state: State1 = from_slice(&data).unwrap();
        assert_eq!(state, expected_state);
    }

    #[test]
    fn add() {
        let mut deps = mock_dependencies_with_balance(&coins(2, "token"));

        let msg = InstantiateMsg {};
        let creator = String::from("creator");
        let info = mock_info(creator.as_str(), &[]);
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // beneficiary can release it
        let info = mock_info("anyone", &[]);
        let msg = ExecuteMsg::Add {delta:Uint256::from(1u128)};
        let _res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();

        // should increase counter by 1
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetCounter {}).unwrap();
        let value: Uint256 = from_binary(&res).unwrap();
        assert_eq!(Uint256::from(1u128), value);
    }
}
