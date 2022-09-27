/*

    Copyright 2020 dYdX Trading Inc.

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

*/

pragma solidity 0.5.16;
pragma experimental ABIEncoderV2;

import { P1Types } from "../lib/P1Types.sol";


/**
 * @title I_PerpetualV1
 * @author dYdX
 *
 * @notice Interface for PerpetualV1.
 */
interface I_PerpetualV1 {

    // ============ Structs ============

    struct TradeArg {
        uint256 takerIndex;
        uint256 makerIndex;
        address trader;
        bytes data;
    }

    // ============ State-Changing Functions ============

    /**
     * @notice Submits one or more trades between any number of accounts.
     *
     * @param  accounts  The sorted list of accounts that are involved in trades.
     * @param  trades    The list of trades to execute in-order.
     */
    function trade(
        address[] calldata accounts,
        TradeArg[] calldata trades
    )
        external;

    /**
     * @notice Withdraw the number of margin tokens equal to the value of the account at the time
     *  that final settlement occurred.
     */
    function withdrawFinalSettlement()
        external;

    /**
     * @notice Deposit some amount of margin tokens from the msg.sender into an account.
     *
     * @param  account  The account for which to credit the deposit.
     * @param  amount   the amount of tokens to deposit.
     */
    function deposit(
        address account,
        uint256 amount
    )
        external;

    /**
     * @notice Withdraw some amount of margin tokens from an account to a destination address.
     *
     * @param  account      The account for which to debit the withdrawal.
     * @param  destination  The address to which the tokens are transferred.
     * @param  amount       The amount of tokens to withdraw.
     */
    function withdraw(
        address account,
        address destination,
        uint256 amount
    )
        external;

    /**
     * @notice Grants or revokes permission for another account to perform certain actions on behalf
     *  of the sender.
     *
     * @param  operator  The account that is approved or disapproved.
     * @param  approved  True for approval, false for disapproval.
     */
    function setLocalOperator(
        address operator,
        bool approved
    )
        external;

    // ============ Account Getters ============

    /**
     * @notice Get the balance of an account, without accounting for changes in the index.
     *
     * @param  account  The address of the account to query the balances of.
     * @return          The balances of the account.
     */
    function getAccountBalance(
        address account
    )
        external
        view
        returns (P1Types.Balance memory);

    /**
     * @notice Gets the most recently cached index of an account.
     *
     * @param  account  The address of the account to query the index of.
     * @return          The index of the account.
     */
    function getAccountIndex(
        address account
    )
        external
        view
        returns (P1Types.Index memory);

    /**
     * @notice Gets the local operator status of an operator for a particular account.
     *
     * @param  account   The account to query the operator for.
     * @param  operator  The address of the operator to query the status of.
     * @return           True if the operator is a local operator of the account, false otherwise.
     */
    function getIsLocalOperator(
        address account,
        address operator
    )
        external
        view
        returns (bool);

    // ============ Global Getters ============

    /**
     * @notice Gets the global operator status of an address.
     *
     * @param  operator  The address of the operator to query the status of.
     * @return           True if the address is a global operator, false otherwise.
     */
    function getIsGlobalOperator(
        address operator
    )
        external
        view
        returns (bool);

    /**
     * @notice Gets the address of the ERC20 margin contract used for margin deposits.
     *
     * @return The address of the ERC20 token.
     */
    function getTokenContract()
        external
        view
        returns (address);

    /**
     * @notice Gets the current address of the price oracle contract.
     *
     * @return The address of the price oracle contract.
     */
    function getOracleContract()
        external
        view
        returns (address);

    /**
     * @notice Gets the current address of the funder contract.
     *
     * @return The address of the funder contract.
     */
    function getFunderContract()
        external
        view
        returns (address);

    /**
     * @notice Gets the most recently cached global index.
     *
     * @return The most recently cached global index.
     */
    function getGlobalIndex()
        external
        view
        returns (P1Types.Index memory);

    /**
     * @notice Gets minimum collateralization ratio of the protocol.
     *
     * @return The minimum-acceptable collateralization ratio, returned as a fixed-point number with
     *  18 decimals of precision.
     */
    function getMinCollateral()
        external
        view
        returns (uint256);

    /**
     * @notice Gets the status of whether final-settlement was initiated by the Admin.
     *
     * @return True if final-settlement was enabled, false otherwise.
     */
    function getFinalSettlementEnabled()
        external
        view
        returns (bool);

    // ============ Public Getters ============

    /**
     * @notice Gets whether an address has permissions to operate an account.
     *
     * @param  account   The account to query.
     * @param  operator  The address to query.
     * @return           True if the operator has permission to operate the account,
     *                   and false otherwise.
     */
    function hasAccountPermissions(
        address account,
        address operator
    )
        external
        view
        returns (bool);

    // ============ Authorized Getters ============

    /**
     * @notice Gets the price returned by the oracle.
     * @dev Only able to be called by global operators.
     *
     * @return The price returned by the current price oracle.
     */
    function getOraclePrice()
        external
        view
        returns (uint256);
}
