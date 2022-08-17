#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{
    coins, from_slice, to_binary, Addr, BankMsg, Binary, Deps, DepsMut, Env, MessageInfo, Order,
    Response, StdResult, Storage, SubMsg, Uint128, WasmMsg,
};

use cw2::set_contract_version;
use cw20::{Balance, Cw20CoinVerified, Cw20ExecuteMsg, Cw20ReceiveMsg, Denom};
use cw4::{
    Member, MemberChangedHookMsg, MemberDiff, MemberListResponse, MemberResponse,
    TotalWeightResponse,
};
use cw_storage_plus::Bound;
use cw_utils::{maybe_addr, NativeBalance};

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg, ReceiveMsg, StakedResponse};
use crate::state::{Config, ADMIN, CLAIMS, CONFIG, HOOKS, MEMBERS, STAKE, TOTAL};

// version info for migration info
const CONTRACT_NAME: &str = "crates.io:cw4-stake";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

// Note, you can use StdResult in some functions where you do not
// make use of the custom errors
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    mut deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;
    let api = deps.api;
    ADMIN.set(deps.branch(), maybe_addr(api, msg.admin)?)?;

    // min_bond is at least 1, so 0 stake -> non-membership
    let min_bond = std::cmp::max(msg.min_bond, Uint128::new(1));

    let config = Config {
        denom: msg.denom,
        tokens_per_weight: msg.tokens_per_weight,
        min_bond,
        unbonding_period: msg.unbonding_period,
    };
    CONFIG.save(deps.storage, &config)?;
    TOTAL.save(deps.storage, &0)?;

    Ok(Response::default())
}

// And declare a custom Error variant for the ones where you will want to make use of it
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    let api = deps.api;
    match msg {
        ExecuteMsg::UpdateAdmin { admin } => {
            Ok(ADMIN.execute_update_admin(deps, info, maybe_addr(api, admin)?)?)
        }
        ExecuteMsg::AddHook { addr } => {
            Ok(HOOKS.execute_add_hook(&ADMIN, deps, info, api.addr_validate(&addr)?)?)
        }
        ExecuteMsg::RemoveHook { addr } => {
            Ok(HOOKS.execute_remove_hook(&ADMIN, deps, info, api.addr_validate(&addr)?)?)
        }
        ExecuteMsg::Bond {} => execute_bond(deps, env, Balance::from(info.funds), info.sender),
        ExecuteMsg::Unbond { tokens: amount } => execute_unbond(deps, env, info, amount),
        ExecuteMsg::Claim {} => execute_claim(deps, env, info),
        ExecuteMsg::Receive(msg) => execute_receive(deps, env, info, msg),
    }
}

pub fn execute_bond(
    deps: DepsMut,
    env: Env,
    amount: Balance,
    sender: Addr,
) -> Result<Response, ContractError> {
    let cfg = CONFIG.load(deps.storage)?;

    // ensure the sent denom was proper
    // NOTE: those clones are not needed (if we move denom, we return early),
    // but the compiler cannot see that (yet...)
    let amount = match (&cfg.denom, &amount) {
        (Denom::Native(want), Balance::Native(have)) => must_pay_funds(have, want),
        (Denom::Cw20(want), Balance::Cw20(have)) => {
            if want == &have.address {
                Ok(have.amount)
            } else {
                Err(ContractError::InvalidDenom(want.into()))
            }
        }
        _ => Err(ContractError::MixedNativeAndCw20(
            "Invalid address or denom".to_string(),
        )),
    }?;

    // update the sender's stake
    let new_stake = STAKE.update(deps.storage, &sender, |stake| -> StdResult<_> {
        Ok(stake.unwrap_or_default() + amount)
    })?;

    let messages = update_membership(
        deps.storage,
        sender.clone(),
        new_stake,
        &cfg,
        env.block.height,
    )?;

    Ok(Response::new()
        .add_submessages(messages)
        .add_attribute("action", "bond")
        .add_attribute("amount", amount)
        .add_attribute("sender", sender))
}

pub fn execute_receive(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    wrapper: Cw20ReceiveMsg,
) -> Result<Response, ContractError> {
    // info.sender is the address of the cw20 contract (that re-sent this message).
    // wrapper.sender is the address of the user that requested the cw20 contract to send this.
    // This cannot be fully trusted (the cw20 contract can fake it), so only use it for actions
    // in the address's favor (like paying/bonding tokens, not withdrawls)
    let msg: ReceiveMsg = from_slice(&wrapper.msg)?;
    let balance = Balance::Cw20(Cw20CoinVerified {
        address: info.sender,
        amount: wrapper.amount,
    });
    let api = deps.api;
    match msg {
        ReceiveMsg::Bond {} => {
            execute_bond(deps, env, balance, api.addr_validate(&wrapper.sender)?)
        }
    }
}

pub fn execute_unbond(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    amount: Uint128,
) -> Result<Response, ContractError> {
    // reduce the sender's stake - aborting if insufficient
    let new_stake = STAKE.update(deps.storage, &info.sender, |stake| -> StdResult<_> {
        Ok(stake.unwrap_or_default().checked_sub(amount)?)
    })?;

    // provide them a claim
    let cfg = CONFIG.load(deps.storage)?;
    CLAIMS.create_claim(
        deps.storage,
        &info.sender,
        amount,
        cfg.unbonding_period.after(&env.block),
    )?;

    let messages = update_membership(
        deps.storage,
        info.sender.clone(),
        new_stake,
        &cfg,
        env.block.height,
    )?;

    Ok(Response::new()
        .add_submessages(messages)
        .add_attribute("action", "unbond")
        .add_attribute("amount", amount)
        .add_attribute("sender", info.sender))
}

pub fn must_pay_funds(balance: &NativeBalance, denom: &str) -> Result<Uint128, ContractError> {
    match balance.0.len() {
        0 => Err(ContractError::NoFunds {}),
        1 => {
            let balance = &balance.0;
            let payment = balance[0].amount;
            if balance[0].denom == denom {
                Ok(payment)
            } else {
                Err(ContractError::MissingDenom(denom.to_string()))
            }
        }
        _ => Err(ContractError::ExtraDenoms(denom.to_string())),
    }
}

fn update_membership(
    storage: &mut dyn Storage,
    sender: Addr,
    new_stake: Uint128,
    cfg: &Config,
    height: u64,
) -> StdResult<Vec<SubMsg>> {
    // update their membership weight
    let new = calc_weight(new_stake, cfg);
    let old = MEMBERS.may_load(storage, &sender)?;

    // short-circuit if no change
    if new == old {
        return Ok(vec![]);
    }
    // otherwise, record change of weight
    match new.as_ref() {
        Some(w) => MEMBERS.save(storage, &sender, w, height),
        None => MEMBERS.remove(storage, &sender, height),
    }?;

    // update total
    TOTAL.update(storage, |total| -> StdResult<_> {
        Ok(total + new.unwrap_or_default() - old.unwrap_or_default())
    })?;

    // alert the hooks
    let diff = MemberDiff::new(sender, old, new);
    HOOKS.prepare_hooks(storage, |h| {
        MemberChangedHookMsg::one(diff.clone())
            .into_cosmos_msg(h)
            .map(SubMsg::new)
    })
}

fn calc_weight(stake: Uint128, cfg: &Config) -> Option<u64> {
    if stake < cfg.min_bond {
        None
    } else {
        let w = stake.u128() / (cfg.tokens_per_weight.u128());
        Some(w as u64)
    }
}

pub fn execute_claim(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
) -> Result<Response, ContractError> {
    let release = CLAIMS.claim_tokens(deps.storage, &info.sender, &env.block, None)?;
    if release.is_zero() {
        return Err(ContractError::NothingToClaim {});
    }

    let config = CONFIG.load(deps.storage)?;
    let (amount_str, message) = match &config.denom {
        Denom::Native(denom) => {
            let amount_str = coin_to_string(release, denom.as_str());
            let amount = coins(release.u128(), denom);
            let message = SubMsg::new(BankMsg::Send {
                to_address: info.sender.to_string(),
                amount,
            });
            (amount_str, message)
        }
        Denom::Cw20(addr) => {
            let amount_str = coin_to_string(release, addr.as_str());
            let transfer = Cw20ExecuteMsg::Transfer {
                recipient: info.sender.clone().into(),
                amount: release,
            };
            let message = SubMsg::new(WasmMsg::Execute {
                contract_addr: addr.into(),
                msg: to_binary(&transfer)?,
                funds: vec![],
            });
            (amount_str, message)
        }
    };

    Ok(Response::new()
        .add_submessage(message)
        .add_attribute("action", "claim")
        .add_attribute("tokens", amount_str)
        .add_attribute("sender", info.sender))
}

#[inline]
fn coin_to_string(amount: Uint128, denom: &str) -> String {
    format!("{} {}", amount, denom)
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::Member {
            addr,
            at_height: height,
        } => to_binary(&query_member(deps, addr, height)?),
        QueryMsg::ListMembers { start_after, limit } => {
            to_binary(&list_members(deps, start_after, limit)?)
        }
        QueryMsg::TotalWeight {} => to_binary(&query_total_weight(deps)?),
        QueryMsg::Claims { address } => {
            to_binary(&CLAIMS.query_claims(deps, &deps.api.addr_validate(&address)?)?)
        }
        QueryMsg::Staked { address } => to_binary(&query_staked(deps, address)?),
        QueryMsg::Admin {} => to_binary(&ADMIN.query_admin(deps)?),
        QueryMsg::Hooks {} => to_binary(&HOOKS.query_hooks(deps)?),
    }
}

fn query_total_weight(deps: Deps) -> StdResult<TotalWeightResponse> {
    let weight = TOTAL.load(deps.storage)?;
    Ok(TotalWeightResponse { weight })
}

pub fn query_staked(deps: Deps, addr: String) -> StdResult<StakedResponse> {
    let addr = deps.api.addr_validate(&addr)?;
    let stake = STAKE.may_load(deps.storage, &addr)?.unwrap_or_default();
    let denom = CONFIG.load(deps.storage)?.denom;
    Ok(StakedResponse { stake, denom })
}

fn query_member(deps: Deps, addr: String, height: Option<u64>) -> StdResult<MemberResponse> {
    let addr = deps.api.addr_validate(&addr)?;
    let weight = match height {
        Some(h) => MEMBERS.may_load_at_height(deps.storage, &addr, h),
        None => MEMBERS.may_load(deps.storage, &addr),
    }?;
    Ok(MemberResponse { weight })
}

// settings for pagination
const MAX_LIMIT: u32 = 30;
const DEFAULT_LIMIT: u32 = 10;

fn list_members(
    deps: Deps,
    start_after: Option<String>,
    limit: Option<u32>,
) -> StdResult<MemberListResponse> {
    let limit = limit.unwrap_or(DEFAULT_LIMIT).min(MAX_LIMIT) as usize;
    let addr = maybe_addr(deps.api, start_after)?;
    let start = addr.as_ref().map(Bound::exclusive);

    let members = MEMBERS
        .range(deps.storage, start, None, Order::Ascending)
        .take(limit)
        .map(|item| {
            item.map(|(addr, weight)| Member {
                addr: addr.into(),
                weight,
            })
        })
        .collect::<StdResult<_>>()?;

    Ok(MemberListResponse { members })
}

#[cfg(test)]
mod tests {
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{
        coin, from_slice, CosmosMsg, OverflowError, OverflowOperation, StdError, Storage,
    };
    use cw20::Denom;
    use cw4::{member_key, TOTAL_KEY};
    use cw_controllers::{AdminError, Claim, HookError};
    use cw_utils::Duration;

    use crate::error::ContractError;

    use super::*;

    const INIT_ADMIN: &str = "juan";
    const USER1: &str = "somebody";
    const USER2: &str = "else";
    const USER3: &str = "funny";
    const DENOM: &str = "stake";
    const TOKENS_PER_WEIGHT: Uint128 = Uint128::new(1_000);
    const MIN_BOND: Uint128 = Uint128::new(5_000);
    const UNBONDING_BLOCKS: u64 = 100;
    const CW20_ADDRESS: &str = "wasm1234567890";

    fn default_instantiate(deps: DepsMut) {
        do_instantiate(
            deps,
            TOKENS_PER_WEIGHT,
            MIN_BOND,
            Duration::Height(UNBONDING_BLOCKS),
        )
    }

    fn do_instantiate(
        deps: DepsMut,
        tokens_per_weight: Uint128,
        min_bond: Uint128,
        unbonding_period: Duration,
    ) {
        let msg = InstantiateMsg {
            denom: Denom::Native("stake".to_string()),
            tokens_per_weight,
            min_bond,
            unbonding_period,
            admin: Some(INIT_ADMIN.into()),
        };
        let info = mock_info("creator", &[]);
        instantiate(deps, mock_env(), info, msg).unwrap();
    }

    fn cw20_instantiate(deps: DepsMut, unbonding_period: Duration) {
        let msg = InstantiateMsg {
            denom: Denom::Cw20(Addr::unchecked(CW20_ADDRESS)),
            tokens_per_weight: TOKENS_PER_WEIGHT,
            min_bond: MIN_BOND,
            unbonding_period,
            admin: Some(INIT_ADMIN.into()),
        };
        let info = mock_info("creator", &[]);
        instantiate(deps, mock_env(), info, msg).unwrap();
    }

    fn bond(mut deps: DepsMut, user1: u128, user2: u128, user3: u128, height_delta: u64) {
        let mut env = mock_env();
        env.block.height += height_delta;

        for (addr, stake) in &[(USER1, user1), (USER2, user2), (USER3, user3)] {
            if *stake != 0 {
                let msg = ExecuteMsg::Bond {};
                let info = mock_info(addr, &coins(*stake, DENOM));
                execute(deps.branch(), env.clone(), info, msg).unwrap();
            }
        }
    }

    fn bond_cw20(mut deps: DepsMut, user1: u128, user2: u128, user3: u128, height_delta: u64) {
        let mut env = mock_env();
        env.block.height += height_delta;

        for (addr, stake) in &[(USER1, user1), (USER2, user2), (USER3, user3)] {
            if *stake != 0 {
                let msg = ExecuteMsg::Receive(Cw20ReceiveMsg {
                    sender: addr.to_string(),
                    amount: Uint128::new(*stake),
                    msg: to_binary(&ReceiveMsg::Bond {}).unwrap(),
                });
                let info = mock_info(CW20_ADDRESS, &[]);
                execute(deps.branch(), env.clone(), info, msg).unwrap();
            }
        }
    }

    fn unbond(mut deps: DepsMut, user1: u128, user2: u128, user3: u128, height_delta: u64) {
        let mut env = mock_env();
        env.block.height += height_delta;

        for (addr, stake) in &[(USER1, user1), (USER2, user2), (USER3, user3)] {
            if *stake != 0 {
                let msg = ExecuteMsg::Unbond {
                    tokens: Uint128::new(*stake),
                };
                let info = mock_info(addr, &[]);
                execute(deps.branch(), env.clone(), info, msg).unwrap();
            }
        }
    }

    #[test]
    fn proper_instantiation() {
        let mut deps = mock_dependencies();
        default_instantiate(deps.as_mut());

        // it worked, let's query the state
        let res = ADMIN.query_admin(deps.as_ref()).unwrap();
        assert_eq!(Some(INIT_ADMIN.into()), res.admin);

        let res = query_total_weight(deps.as_ref()).unwrap();
        assert_eq!(0, res.weight);
    }

    fn get_member(deps: Deps, addr: String, at_height: Option<u64>) -> Option<u64> {
        let raw = query(deps, mock_env(), QueryMsg::Member { addr, at_height }).unwrap();
        let res: MemberResponse = from_slice(&raw).unwrap();
        res.weight
    }

    // this tests the member queries
    fn assert_users(
        deps: Deps,
        user1_weight: Option<u64>,
        user2_weight: Option<u64>,
        user3_weight: Option<u64>,
        height: Option<u64>,
    ) {
        let member1 = get_member(deps, USER1.into(), height);
        assert_eq!(member1, user1_weight);

        let member2 = get_member(deps, USER2.into(), height);
        assert_eq!(member2, user2_weight);

        let member3 = get_member(deps, USER3.into(), height);
        assert_eq!(member3, user3_weight);

        // this is only valid if we are not doing a historical query
        if height.is_none() {
            // compute expected metrics
            let weights = vec![user1_weight, user2_weight, user3_weight];
            let sum: u64 = weights.iter().map(|x| x.unwrap_or_default()).sum();
            let count = weights.iter().filter(|x| x.is_some()).count();

            // TODO: more detailed compare?
            let msg = QueryMsg::ListMembers {
                start_after: None,
                limit: None,
            };
            let raw = query(deps, mock_env(), msg).unwrap();
            let members: MemberListResponse = from_slice(&raw).unwrap();
            assert_eq!(count, members.members.len());

            let raw = query(deps, mock_env(), QueryMsg::TotalWeight {}).unwrap();
            let total: TotalWeightResponse = from_slice(&raw).unwrap();
            assert_eq!(sum, total.weight); // 17 - 11 + 15 = 21
        }
    }

    // this tests the member queries
    fn assert_stake(deps: Deps, user1_stake: u128, user2_stake: u128, user3_stake: u128) {
        let stake1 = query_staked(deps, USER1.into()).unwrap();
        assert_eq!(stake1.stake, user1_stake.into());

        let stake2 = query_staked(deps, USER2.into()).unwrap();
        assert_eq!(stake2.stake, user2_stake.into());

        let stake3 = query_staked(deps, USER3.into()).unwrap();
        assert_eq!(stake3.stake, user3_stake.into());
    }

    #[test]
    fn bond_stake_adds_membership() {
        let mut deps = mock_dependencies();
        default_instantiate(deps.as_mut());
        let height = mock_env().block.height;

        // Assert original weights
        assert_users(deps.as_ref(), None, None, None, None);

        // ensure it rounds down, and respects cut-off
        bond(deps.as_mut(), 12_000, 7_500, 4_000, 1);

        // Assert updated weights
        assert_stake(deps.as_ref(), 12_000, 7_500, 4_000);
        assert_users(deps.as_ref(), Some(12), Some(7), None, None);

        // add some more, ensure the sum is properly respected (7.5 + 7.6 = 15 not 14)
        bond(deps.as_mut(), 0, 7_600, 1_200, 2);

        // Assert updated weights
        assert_stake(deps.as_ref(), 12_000, 15_100, 5_200);
        assert_users(deps.as_ref(), Some(12), Some(15), Some(5), None);

        // check historical queries all work
        assert_users(deps.as_ref(), None, None, None, Some(height + 1)); // before first stake
        assert_users(deps.as_ref(), Some(12), Some(7), None, Some(height + 2)); // after first stake
        assert_users(deps.as_ref(), Some(12), Some(15), Some(5), Some(height + 3));
        // after second stake
    }

    #[test]
    fn unbond_stake_update_membership() {
        let mut deps = mock_dependencies();
        default_instantiate(deps.as_mut());
        let height = mock_env().block.height;

        // ensure it rounds down, and respects cut-off
        bond(deps.as_mut(), 12_000, 7_500, 4_000, 1);
        unbond(deps.as_mut(), 4_500, 2_600, 1_111, 2);

        // Assert updated weights
        assert_stake(deps.as_ref(), 7_500, 4_900, 2_889);
        assert_users(deps.as_ref(), Some(7), None, None, None);

        // Adding a little more returns weight
        bond(deps.as_mut(), 600, 100, 2_222, 3);

        // Assert updated weights
        assert_users(deps.as_ref(), Some(8), Some(5), Some(5), None);

        // check historical queries all work
        assert_users(deps.as_ref(), None, None, None, Some(height + 1)); // before first stake
        assert_users(deps.as_ref(), Some(12), Some(7), None, Some(height + 2)); // after first bond
        assert_users(deps.as_ref(), Some(7), None, None, Some(height + 3)); // after first unbond
        assert_users(deps.as_ref(), Some(8), Some(5), Some(5), Some(height + 4)); // after second bond

        // error if try to unbond more than stake (USER2 has 5000 staked)
        let msg = ExecuteMsg::Unbond {
            tokens: Uint128::new(5100),
        };
        let mut env = mock_env();
        env.block.height += 5;
        let info = mock_info(USER2, &[]);
        let err = execute(deps.as_mut(), env, info, msg).unwrap_err();
        assert_eq!(
            err,
            ContractError::Std(StdError::overflow(OverflowError::new(
                OverflowOperation::Sub,
                5000,
                5100
            )))
        );
    }

    #[test]
    fn cw20_token_bond() {
        let mut deps = mock_dependencies();
        cw20_instantiate(deps.as_mut(), Duration::Height(2000));

        // Assert original weights
        assert_users(deps.as_ref(), None, None, None, None);

        // ensure it rounds down, and respects cut-off
        bond_cw20(deps.as_mut(), 12_000, 7_500, 4_000, 1);

        // Assert updated weights
        assert_stake(deps.as_ref(), 12_000, 7_500, 4_000);
        assert_users(deps.as_ref(), Some(12), Some(7), None, None);
    }

    #[test]
    fn cw20_token_claim() {
        let unbonding_period: u64 = 50;
        let unbond_height: u64 = 10;

        let mut deps = mock_dependencies();
        let unbonding = Duration::Height(unbonding_period);
        cw20_instantiate(deps.as_mut(), unbonding);

        // bond some tokens
        bond_cw20(deps.as_mut(), 20_000, 13_500, 500, 1);

        // unbond part
        unbond(deps.as_mut(), 7_900, 4_600, 0, unbond_height);

        // Assert updated weights
        assert_stake(deps.as_ref(), 12_100, 8_900, 500);
        assert_users(deps.as_ref(), Some(12), Some(8), None, None);

        // with proper claims
        let mut env = mock_env();
        env.block.height += unbond_height;
        let expires = unbonding.after(&env.block);
        assert_eq!(
            get_claims(deps.as_ref(), &Addr::unchecked(USER1)),
            vec![Claim::new(7_900, expires)]
        );

        // wait til they expire and get payout
        env.block.height += unbonding_period;
        let res = execute(
            deps.as_mut(),
            env,
            mock_info(USER1, &[]),
            ExecuteMsg::Claim {},
        )
        .unwrap();
        assert_eq!(res.messages.len(), 1);
        match &res.messages[0].msg {
            CosmosMsg::Wasm(WasmMsg::Execute {
                contract_addr,
                msg,
                funds,
            }) => {
                assert_eq!(contract_addr.as_str(), CW20_ADDRESS);
                assert_eq!(funds.len(), 0);
                let parsed: Cw20ExecuteMsg = from_slice(msg).unwrap();
                assert_eq!(
                    parsed,
                    Cw20ExecuteMsg::Transfer {
                        recipient: USER1.into(),
                        amount: Uint128::new(7_900)
                    }
                );
            }
            _ => panic!("Must initiate cw20 transfer"),
        }
    }

    #[test]
    fn raw_queries_work() {
        // add will over-write and remove have no effect
        let mut deps = mock_dependencies();
        default_instantiate(deps.as_mut());
        // Set values as (11, 6, None)
        bond(deps.as_mut(), 11_000, 6_000, 0, 1);

        // get total from raw key
        let total_raw = deps.storage.get(TOTAL_KEY.as_bytes()).unwrap();
        let total: u64 = from_slice(&total_raw).unwrap();
        assert_eq!(17, total);

        // get member votes from raw key
        let member2_raw = deps.storage.get(&member_key(USER2)).unwrap();
        let member2: u64 = from_slice(&member2_raw).unwrap();
        assert_eq!(6, member2);

        // and execute misses
        let member3_raw = deps.storage.get(&member_key(USER3));
        assert_eq!(None, member3_raw);
    }

    fn get_claims(deps: Deps, addr: &Addr) -> Vec<Claim> {
        CLAIMS.query_claims(deps, addr).unwrap().claims
    }

    #[test]
    fn unbond_claim_workflow() {
        let mut deps = mock_dependencies();
        default_instantiate(deps.as_mut());

        // create some data
        bond(deps.as_mut(), 12_000, 7_500, 4_000, 1);
        unbond(deps.as_mut(), 4_500, 2_600, 0, 2);
        let mut env = mock_env();
        env.block.height += 2;

        // check the claims for each user
        let expires = Duration::Height(UNBONDING_BLOCKS).after(&env.block);
        assert_eq!(
            get_claims(deps.as_ref(), &Addr::unchecked(USER1)),
            vec![Claim::new(4_500, expires)]
        );
        assert_eq!(
            get_claims(deps.as_ref(), &Addr::unchecked(USER2)),
            vec![Claim::new(2_600, expires)]
        );
        assert_eq!(get_claims(deps.as_ref(), &Addr::unchecked(USER3)), vec![]);

        // do another unbond later on
        let mut env2 = mock_env();
        env2.block.height += 22;
        unbond(deps.as_mut(), 0, 1_345, 1_500, 22);

        // with updated claims
        let expires2 = Duration::Height(UNBONDING_BLOCKS).after(&env2.block);
        assert_eq!(
            get_claims(deps.as_ref(), &Addr::unchecked(USER1)),
            vec![Claim::new(4_500, expires)]
        );
        assert_eq!(
            get_claims(deps.as_ref(), &Addr::unchecked(USER2)),
            vec![Claim::new(2_600, expires), Claim::new(1_345, expires2)]
        );
        assert_eq!(
            get_claims(deps.as_ref(), &Addr::unchecked(USER3)),
            vec![Claim::new(1_500, expires2)]
        );

        // nothing can be withdrawn yet
        let err = execute(
            deps.as_mut(),
            env2,
            mock_info(USER1, &[]),
            ExecuteMsg::Claim {},
        )
        .unwrap_err();
        assert_eq!(err, ContractError::NothingToClaim {});

        // now mature first section, withdraw that
        let mut env3 = mock_env();
        env3.block.height += 2 + UNBONDING_BLOCKS;
        // first one can now release
        let res = execute(
            deps.as_mut(),
            env3.clone(),
            mock_info(USER1, &[]),
            ExecuteMsg::Claim {},
        )
        .unwrap();
        assert_eq!(
            res.messages,
            vec![SubMsg::new(BankMsg::Send {
                to_address: USER1.into(),
                amount: coins(4_500, DENOM),
            })]
        );

        // second releases partially
        let res = execute(
            deps.as_mut(),
            env3.clone(),
            mock_info(USER2, &[]),
            ExecuteMsg::Claim {},
        )
        .unwrap();
        assert_eq!(
            res.messages,
            vec![SubMsg::new(BankMsg::Send {
                to_address: USER2.into(),
                amount: coins(2_600, DENOM),
            })]
        );

        // but the third one cannot release
        let err = execute(
            deps.as_mut(),
            env3,
            mock_info(USER3, &[]),
            ExecuteMsg::Claim {},
        )
        .unwrap_err();
        assert_eq!(err, ContractError::NothingToClaim {});

        // claims updated properly
        assert_eq!(get_claims(deps.as_ref(), &Addr::unchecked(USER1)), vec![]);
        assert_eq!(
            get_claims(deps.as_ref(), &Addr::unchecked(USER2)),
            vec![Claim::new(1_345, expires2)]
        );
        assert_eq!(
            get_claims(deps.as_ref(), &Addr::unchecked(USER3)),
            vec![Claim::new(1_500, expires2)]
        );

        // add another few claims for 2
        unbond(deps.as_mut(), 0, 600, 0, 30 + UNBONDING_BLOCKS);
        unbond(deps.as_mut(), 0, 1_005, 0, 50 + UNBONDING_BLOCKS);

        // ensure second can claim all tokens at once
        let mut env4 = mock_env();
        env4.block.height += 55 + UNBONDING_BLOCKS + UNBONDING_BLOCKS;
        let res = execute(
            deps.as_mut(),
            env4,
            mock_info(USER2, &[]),
            ExecuteMsg::Claim {},
        )
        .unwrap();
        assert_eq!(
            res.messages,
            vec![SubMsg::new(BankMsg::Send {
                to_address: USER2.into(),
                // 1_345 + 600 + 1_005
                amount: coins(2_950, DENOM),
            })]
        );
        assert_eq!(get_claims(deps.as_ref(), &Addr::unchecked(USER2)), vec![]);
    }

    #[test]
    fn add_remove_hooks() {
        // add will over-write and remove have no effect
        let mut deps = mock_dependencies();
        default_instantiate(deps.as_mut());

        let hooks = HOOKS.query_hooks(deps.as_ref()).unwrap();
        assert!(hooks.hooks.is_empty());

        let contract1 = String::from("hook1");
        let contract2 = String::from("hook2");

        let add_msg = ExecuteMsg::AddHook {
            addr: contract1.clone(),
        };

        // non-admin cannot add hook
        let user_info = mock_info(USER1, &[]);
        let err = execute(
            deps.as_mut(),
            mock_env(),
            user_info.clone(),
            add_msg.clone(),
        )
        .unwrap_err();
        assert_eq!(err, HookError::Admin(AdminError::NotAdmin {}).into());

        // admin can add it, and it appears in the query
        let admin_info = mock_info(INIT_ADMIN, &[]);
        let _ = execute(
            deps.as_mut(),
            mock_env(),
            admin_info.clone(),
            add_msg.clone(),
        )
        .unwrap();
        let hooks = HOOKS.query_hooks(deps.as_ref()).unwrap();
        assert_eq!(hooks.hooks, vec![contract1.clone()]);

        // cannot remove a non-registered contract
        let remove_msg = ExecuteMsg::RemoveHook {
            addr: contract2.clone(),
        };
        let err = execute(deps.as_mut(), mock_env(), admin_info.clone(), remove_msg).unwrap_err();
        assert_eq!(err, HookError::HookNotRegistered {}.into());

        // add second contract
        let add_msg2 = ExecuteMsg::AddHook {
            addr: contract2.clone(),
        };
        let _ = execute(deps.as_mut(), mock_env(), admin_info.clone(), add_msg2).unwrap();
        let hooks = HOOKS.query_hooks(deps.as_ref()).unwrap();
        assert_eq!(hooks.hooks, vec![contract1.clone(), contract2.clone()]);

        // cannot re-add an existing contract
        let err = execute(deps.as_mut(), mock_env(), admin_info.clone(), add_msg).unwrap_err();
        assert_eq!(err, HookError::HookAlreadyRegistered {}.into());

        // non-admin cannot remove
        let remove_msg = ExecuteMsg::RemoveHook { addr: contract1 };
        let err = execute(deps.as_mut(), mock_env(), user_info, remove_msg.clone()).unwrap_err();
        assert_eq!(err, HookError::Admin(AdminError::NotAdmin {}).into());

        // remove the original
        let _ = execute(deps.as_mut(), mock_env(), admin_info, remove_msg).unwrap();
        let hooks = HOOKS.query_hooks(deps.as_ref()).unwrap();
        assert_eq!(hooks.hooks, vec![contract2]);
    }

    #[test]
    fn hooks_fire() {
        let mut deps = mock_dependencies();
        default_instantiate(deps.as_mut());

        let hooks = HOOKS.query_hooks(deps.as_ref()).unwrap();
        assert!(hooks.hooks.is_empty());

        let contract1 = String::from("hook1");
        let contract2 = String::from("hook2");

        // register 2 hooks
        let admin_info = mock_info(INIT_ADMIN, &[]);
        let add_msg = ExecuteMsg::AddHook {
            addr: contract1.clone(),
        };
        let add_msg2 = ExecuteMsg::AddHook {
            addr: contract2.clone(),
        };
        for msg in vec![add_msg, add_msg2] {
            let _ = execute(deps.as_mut(), mock_env(), admin_info.clone(), msg).unwrap();
        }

        // check firing on bond
        assert_users(deps.as_ref(), None, None, None, None);
        let info = mock_info(USER1, &coins(13_800, DENOM));
        let res = execute(deps.as_mut(), mock_env(), info, ExecuteMsg::Bond {}).unwrap();
        assert_users(deps.as_ref(), Some(13), None, None, None);

        // ensure messages for each of the 2 hooks
        assert_eq!(res.messages.len(), 2);
        let diff = MemberDiff::new(USER1, None, Some(13));
        let hook_msg = MemberChangedHookMsg::one(diff);
        let msg1 = SubMsg::new(hook_msg.clone().into_cosmos_msg(contract1.clone()).unwrap());
        let msg2 = SubMsg::new(hook_msg.into_cosmos_msg(contract2.clone()).unwrap());
        assert_eq!(res.messages, vec![msg1, msg2]);

        // check firing on unbond
        let msg = ExecuteMsg::Unbond {
            tokens: Uint128::new(7_300),
        };
        let info = mock_info(USER1, &[]);
        let res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_users(deps.as_ref(), Some(6), None, None, None);

        // ensure messages for each of the 2 hooks
        assert_eq!(res.messages.len(), 2);
        let diff = MemberDiff::new(USER1, Some(13), Some(6));
        let hook_msg = MemberChangedHookMsg::one(diff);
        let msg1 = SubMsg::new(hook_msg.clone().into_cosmos_msg(contract1).unwrap());
        let msg2 = SubMsg::new(hook_msg.into_cosmos_msg(contract2).unwrap());
        assert_eq!(res.messages, vec![msg1, msg2]);
    }

    #[test]
    fn only_bond_valid_coins() {
        let mut deps = mock_dependencies();
        default_instantiate(deps.as_mut());

        // cannot bond with 0 coins
        let info = mock_info(USER1, &[]);
        let err = execute(deps.as_mut(), mock_env(), info, ExecuteMsg::Bond {}).unwrap_err();
        assert_eq!(err, ContractError::NoFunds {});

        // cannot bond with incorrect denom
        let info = mock_info(USER1, &[coin(500, "FOO")]);
        let err = execute(deps.as_mut(), mock_env(), info, ExecuteMsg::Bond {}).unwrap_err();
        assert_eq!(err, ContractError::MissingDenom(DENOM.to_string()));

        // cannot bond with 2 coins (even if one is correct)
        let info = mock_info(USER1, &[coin(1234, DENOM), coin(5000, "BAR")]);
        let err = execute(deps.as_mut(), mock_env(), info, ExecuteMsg::Bond {}).unwrap_err();
        assert_eq!(err, ContractError::ExtraDenoms(DENOM.to_string()));

        // can bond with just the proper denom
        // cannot bond with incorrect denom
        let info = mock_info(USER1, &[coin(500, DENOM)]);
        execute(deps.as_mut(), mock_env(), info, ExecuteMsg::Bond {}).unwrap();
    }

    #[test]
    fn ensure_bonding_edge_cases() {
        // use min_bond 0, tokens_per_weight 500
        let mut deps = mock_dependencies();
        do_instantiate(
            deps.as_mut(),
            Uint128::new(100),
            Uint128::zero(),
            Duration::Height(5),
        );

        // setting 50 tokens, gives us Some(0) weight
        // even setting to 1 token
        bond(deps.as_mut(), 50, 1, 102, 1);
        assert_users(deps.as_ref(), Some(0), Some(0), Some(1), None);

        // reducing to 0 token makes us None even with min_bond 0
        unbond(deps.as_mut(), 49, 1, 102, 2);
        assert_users(deps.as_ref(), Some(0), None, None, None);
    }
}
