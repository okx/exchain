use cosmwasm_std::{entry_point, from_slice, to_binary, AllBalanceResponse, BankMsg, Deps, DepsMut, Env, Event, MessageInfo, QueryResponse, Response, StdError, StdResult, Binary, Uint256, WasmMsg, coin, to_vec, ContractResult};

use crate::errors::HackError;
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg, VerifierResponse, MigrateMsg};
use crate::state::{State, CONFIG_KEY, CONFIG_KEY1, CONFIG_KEY2, State1};

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
    deps.storage.set(
        CONFIG_KEY2,
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
        ExecuteMsg::Call { delta, addr } => call_ret(deps, delta, addr, &env),
        ExecuteMsg::DelegateCall { delta, addr } => del_call_ret(deps, delta, addr, &env),
    }
}

pub fn call_ret(deps: DepsMut, delta:Uint256, callee_addr:String, _env: &Env) -> Result<Response, HackError> {
    let msg_str = format!("{{\"add\":{{\"delta\":\"{}\"}}}}", delta);
    let msg_b = Binary(msg_str.into_bytes());
    let send_msg = WasmMsg::Execute {
        contract_addr: callee_addr,
        msg: msg_b,
        funds: vec![]
    };
    let result = deps.api.call(_env, &send_msg);
    match result {
        Ok(data1) => {
            let pret = String::from_utf8(data1.clone()).unwrap();
            deps.api.debug(pret.as_str());
            let result: ContractResult<Response> = from_slice(data1.as_slice()).unwrap_or_else(|_| {
                ContractResult::Err(format!("this is contract call err about unwrap"))
            });
            if result.is_err() {
                deps.api.debug(format!("this is contract err: {:?}", result).as_str());
            } else {
                let deserialized = result.unwrap();
                let parsed_num: Result<Uint256, _> = deserialized.attributes[0].value.clone().parse();
                deps.api.debug(format!("this is contract: {:?}", parsed_num).as_str());
                match parsed_num {
                    Ok(num) => {
                        let data = deps
                            .storage
                            .get(CONFIG_KEY2)
                            .ok_or_else(|| StdError::not_found("State1"))?;
                        let state: State1 = from_slice(&data)?;
                        let mut counter= Uint256::from(state.counter);

                        counter += num;
                        deps.storage.set(
                            CONFIG_KEY2,
                            &to_vec(&State1 {
                                counter: counter,
                            })?,
                        );
                    },
                    Err(_) => {
                        deps.api.debug("Failed to parse number.");
                    }
                }
            }
        }
        Err(StdError::GenericErr {msg: mssage, .. } )=> {
            deps.api.debug(format!("this is contract err: {:?}", mssage).as_str());
        }
        Err(err) => {
            deps.api.debug("this is contract err");
        }
    }
    Ok(Response::new())
}

pub fn del_call_ret(deps: DepsMut, delta:Uint256, callee_addr:String, _env: &Env) -> Result<Response, HackError> {
    let msg_str = format!("{{\"add\":{{\"delta\":\"{}\"}}}}", delta);
    let msg_b = Binary(msg_str.into_bytes());
    let send_msg = WasmMsg::Execute {
        contract_addr: callee_addr,
        msg: msg_b,
        funds: vec![]
    };
    let result = deps.api.delegate_call(_env, &send_msg);
    match result {
        Ok(data1) => {
            let pret = String::from_utf8(data1.clone()).unwrap();
            deps.api.debug(pret.as_str());
            let result: ContractResult<Response> = from_slice(data1.as_slice()).unwrap_or_else(|_| {
                ContractResult::Err(format!("this is contract call err about unwrap"))
            });
            if result.is_err() {
                deps.api.debug(format!("this is contract err: {:?}", result).as_str());
            } else {
                let deserialized = result.unwrap();
                let parsed_num: Result<Uint256, _> = deserialized.attributes[0].value.clone().parse();
                deps.api.debug(format!("this is contract: {:?}", parsed_num).as_str());
                match parsed_num {
                    Ok(num) => {
                        let data = deps
                            .storage
                            .get(CONFIG_KEY2)
                            .ok_or_else(|| StdError::not_found("State1"))?;
                        let state: State1 = from_slice(&data)?;
                        let mut counter= Uint256::from(state.counter);

                        counter += num;
                        deps.storage.set(
                            CONFIG_KEY2,
                            &to_vec(&State1 {
                                counter: counter,
                            })?,
                        );
                    },
                    Err(_) => {
                        deps.api.debug("Failed to parse number.");
                    }
                }
            }
        }
        Err(StdError::GenericErr {msg: mssage, .. } )=> {
            deps.api.debug(format!("this is contract err: {:?}", mssage).as_str());
        }
        Err(err) => {
            deps.api.debug("this is contract err");
        }
    }
    Ok(Response::new())
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
        .get(CONFIG_KEY2)
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

        let msg = InstantiateMsg {
            addr: String::from("to call")
        };
        let info = mock_info(creator.as_str(), &[]);
        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());

        // it worked, let's check the state
        let data = deps.storage.get(CONFIG_KEY1).expect("no data stored");
        let state: State1 = from_slice(&data).unwrap();
        assert_eq!(state, expected_state);
    }

    #[test]
    fn call() {
        let mut deps = mock_dependencies_with_balance(&coins(2, "token"));

        let msg = InstantiateMsg {
            addr: String::from("")
        };
        let creator = String::from("creator");
        let info = mock_info(creator.as_str(), &[]);
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // beneficiary can release it
        let info = mock_info("anyone", &[]);
        let msg = ExecuteMsg::Call {delta:Uint256::from(1u128), addr:String::from("to_call")};
        let _res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();

        // should increase counter by 1
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetCounter {}).unwrap();
        let value: Uint256 = from_binary(&res).unwrap();
        assert_eq!(Uint256::from(0u128), value);
    }

    #[test]
    fn del_call() {
        let mut deps = mock_dependencies_with_balance(&coins(2, "token"));

        let msg = InstantiateMsg {
            addr: String::from("")
        };
        let creator = String::from("creator");
        let info = mock_info(creator.as_str(), &[]);
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // beneficiary can release it
        let info = mock_info("anyone", &[]);
        let msg = ExecuteMsg::DelegateCall {delta:Uint256::from(1u128), addr:String::from("to_call")};
        let _res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();

        // should increase counter by 1
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetCounter {}).unwrap();
        let value: Uint256 = from_binary(&res).unwrap();
        assert_eq!(Uint256::from(0u128), value);
    }
}
