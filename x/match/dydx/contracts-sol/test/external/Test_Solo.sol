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

import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { I_Solo } from "../../external/dydx/I_Solo.sol";


/**
 * @title Test_Solo
 * @author dYdX
 *
 * Interface for calling the Solo margin smart contract.
 */
/* solium-disable-next-line camelcase */
contract Test_Solo is
  I_Solo
{
    // ============ Test Data ============

    mapping(address => bool) internal _GLOBAL_OPERATORS_;
    mapping(address => mapping(address => bool)) internal _LOCAL_OPERATORS_;
    mapping(uint256 => address) public _TOKEN_ADDRESSES_;

    // ============ Test Data Setter Functions ============

    function setIsLocalOperator(
        address owner,
        address operator,
        bool approved
    )
        external
        returns (bool)
    {
        return _LOCAL_OPERATORS_[owner][operator] = approved;
    }

    function setIsGlobalOperator(
        address operator,
        bool approved
    )
        external
        returns (bool)
    {
        return _GLOBAL_OPERATORS_[operator] = approved;
    }

    function setTokenAddress(
        uint256 marketId,
        address tokenAddress
    )
        external
    {
        _TOKEN_ADDRESSES_[marketId] = tokenAddress;
    }

    // ============ Getter Functions ============

    /**
     * Return true if a particular address is approved as an operator for an owner's accounts.
     * Approved operators can act on the accounts of the owner as if it were the operator's own.
     *
     * @param  owner     The owner of the accounts
     * @param  operator  The possible operator
     * @return           True if operator is approved for owner's accounts
     */
    function getIsLocalOperator(
        address owner,
        address operator
    )
        external
        view
        returns (bool)
    {
        return _LOCAL_OPERATORS_[owner][operator];
    }

    /**
     * Return true if a particular address is approved as a global operator. Such an address can
     * act on any account as if it were the operator's own.
     *
     * @param  operator  The address to query
     * @return           True if operator is a global operator
     */
    function getIsGlobalOperator(
        address operator
    )
        external
        view
        returns (bool)
    {
        return _GLOBAL_OPERATORS_[operator];
    }

    /**
     * @notice Get the ERC20 token address for a market.
     *
     * @param  marketId  The market to query
     * @return           The token address
     */
    function getMarketTokenAddress(
        uint256 marketId
    )
        public
        view
        returns (address)
    {
        return _TOKEN_ADDRESSES_[marketId];
    }

    // ============ State-Changing Functions ============

    /**
     * @notice The main entry-point to Solo that allows users and contracts to manage accounts.
     *  Takes one or more actions on one or more accounts. The msg.sender must be the owner or
     *  operator of all accounts except for those being liquidated, vaporized, or traded with.
     *  One call to operate() is considered a singular "operation". Account collateralization is
     *  ensured only after the completion of the entire operation.
     *
     * @param  accounts  A list of all accounts that will be used in this operation. Cannot contain
     *                   duplicates. In each action, the relevant account will be referred to by its
     *                   index in the list.
     * @param  actions   An ordered list of all actions that will be taken in this operation. The
     *                   actions will be processed in order.
     */
    function operate(
        I_Solo.AccountInfo[] memory accounts,
        I_Solo.ActionArgs[] memory actions
    )
        public // public instead of external to avoid UnimplementedFeatureError
    {
        // Expect exactly one account and one action.
        require(accounts.length == 1, "Expected one account");
        require(actions.length == 1, "Expected one action");

        I_Solo.AccountInfo memory account = accounts[0];
        I_Solo.ActionArgs memory action = actions[0];

        // Get the ERC20 token.
        IERC20 token = IERC20(getMarketTokenAddress(action.primaryMarketId));

        // Compare account and action parameters.
        require(account.number == action.accountId, "Account ID mismatch");

        // Check amount parameters.
        I_Solo.AssetAmount memory amount = action.amount;
        require(
            amount.denomination == I_Solo.AssetDenomination.Wei,
            "Expected amount denomination to be Wei"
        );
        uint256 amountToTransfer;
        if (amount.ref == I_Solo.AssetReference.Delta) {
            amountToTransfer = amount.value;
        } else {
            require(amount.value == 0, "When using AssetReference.Target, expect value to be zero");
            require(
                action.actionType == I_Solo.ActionType.Withdraw,
                "When using AssetReference.Target, expect action to be withdrawal"
            );

            // Assume that the whole token balance belongs to the sender.
            amountToTransfer = token.balanceOf(address(this));
        }

        if (action.actionType == I_Solo.ActionType.Withdraw) {
            require(!amount.sign, "Expected amount to be negative");

            // Perform token transfer.
            token.transfer(action.otherAddress, amountToTransfer);
        } else if (action.actionType == I_Solo.ActionType.Deposit) {
            require(amount.sign, "Expected amount to be positive");

            // Perform token transfer.
            token.transferFrom(action.otherAddress, address(this), amountToTransfer);
        } else {
            revert("Expected action type to be Withdraw or Deposit");
        }
    }
}
