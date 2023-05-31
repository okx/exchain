#[cfg(not(feature = "library"))]
use hex::encode as hex_encode;
use sha3::{Digest, Keccak256};
use cosmwasm_std::entry_point;
use cosmwasm_std::{
    Deps, to_binary,Binary, DepsMut, Env, MessageInfo, Response, StdResult, Uint128,CosmosMsg
};

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, InstantiateMsg,QueryMsg,CallToEvmMsg,MigrateMsg};
use crate::state:: COUNTER;


#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, ContractError> {

    COUNTER.save(deps.storage, &Uint128::zero())?;
    Ok(Response::new()
        .add_attribute("method", "instantiate"))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    env: Env,
    _info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response<CallToEvmMsg>, ContractError> {
    match msg {
        //Add the count in the wasm contract
        ExecuteMsg::Add {delta} => try_add(deps,delta),
        //Add the count in the evm contract
        ExecuteMsg::AddCounterForEvm {evm_contract, delta} => try_add_counter_for_evm(env,evm_contract, delta),
    }
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn migrate(_deps: DepsMut, _env: Env, _msg: MigrateMsg) -> StdResult<Response> {
    Ok(Response::default())
}

pub fn try_add_counter_for_evm(env:Env ,the_evm_contract:String, delta:Uint128) -> Result<Response<CallToEvmMsg>, ContractError> {

    //Splicing "calldata" of EVM
    //sign_data:EVM function selector
    let mut sign_data = short_signature("add(uint128)").to_vec();
    sign_data.append(&mut encode(delta.u128()));

    //Vec<u8> => hex;
    let evm_calldata = hex_encode(&sign_data);


    //The specific message "CosmosMsg::Custom(CallToEvmMsg)" can trigger an evm transaction
    let message = CosmosMsg::Custom(CallToEvmMsg {
        sender: env.contract.address.to_string(), //wasm contract address(from)
        evmaddr: the_evm_contract,//evm contract address(to)
        calldata: evm_calldata, //calldata
        value: Uint128::zero(), //The native token you want to send
    });

    Ok(Response::new().add_message(message))
}

pub fn try_add(deps:DepsMut ,delta:Uint128) -> Result<Response<CallToEvmMsg>, ContractError> {
    
    let _res =COUNTER.update(deps.storage, | old| -> StdResult<_> {
        Ok(old.checked_add(delta).unwrap())
    });

    Ok(Response::new().add_attribute("add",delta))
}

pub fn short_signature(func_name: &str) -> [u8; 4] {
	let mut result = [0u8; 4];
	result.copy_from_slice(&Keccak256::digest(func_name)[..4]);
	result
}

pub fn encode(delta:u128)->Vec<u8>{
    let mut a = [0u8; 16].to_vec();
    a.append(&mut delta.to_be_bytes().to_vec());
    a
}


#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetCounter {} => to_binary(&query_counter(deps)?),
    }
}

pub fn query_counter(deps: Deps) -> StdResult<Uint128> {

    let info = COUNTER.may_load(deps.storage).unwrap_or_default();
    Ok(info.unwrap())
}