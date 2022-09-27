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

import { SafeMath } from "@openzeppelin/contracts/math/SafeMath.sol";
import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/SafeERC20.sol";
import { P1Settlement } from "./P1Settlement.sol";
import { BaseMath } from "../../lib/BaseMath.sol";
import { Math } from "../../lib/Math.sol";
import { P1BalanceMath } from "../lib/P1BalanceMath.sol";
import { P1Types } from "../lib/P1Types.sol";


/**
 * @title P1FinalSettlement
 * @author dYdX
 *
 * @notice Functions regulating the smart contract's behavior during final settlement.
 */
contract P1FinalSettlement is
    P1Settlement
{
    using SafeMath for uint256;

    // ============ Events ============

    event LogWithdrawFinalSettlement(
        address indexed account,
        uint256 amount,
        bytes32 balance
    );

    // ============ Modifiers ============

    /**
    * @dev Modifier to ensure the function is not run after final settlement has been enabled.
    */
    modifier noFinalSettlement() {
        require(
            !_FINAL_SETTLEMENT_ENABLED_,
            "Not permitted during final settlement"
        );
        _;
    }

    /**
    * @dev Modifier to ensure the function is only run after final settlement has been enabled.
    */
    modifier onlyFinalSettlement() {
        require(
            _FINAL_SETTLEMENT_ENABLED_,
            "Only permitted during final settlement"
        );
        _;
    }

    // ============ Functions ============

    /**
     * @notice Withdraw the number of margin tokens equal to the value of the account at the time
     *  that final settlement occurred.
     * @dev Emits the LogAccountSettled and LogWithdrawFinalSettlement events.
     */
    function withdrawFinalSettlement()
        external
        onlyFinalSettlement
        nonReentrant
    {
        // Load the context using the final settlement price.
        P1Types.Context memory context = P1Types.Context({
            price: _FINAL_SETTLEMENT_PRICE_,
            minCollateral: _MIN_COLLATERAL_,
            index: _GLOBAL_INDEX_
        });

        // Apply funding changes.
        P1Types.Balance memory balance = _settleAccount(context, msg.sender);

        // Determine the account net value.
        // `positive` and `negative` are base values with extra precision.
        (uint256 positive, uint256 negative) = P1BalanceMath.getPositiveAndNegativeValue(
            balance,
            context.price
        );

        // No amount is withdrawable.
        if (positive < negative) {
            return;
        }

        // Get the account value, which is rounded down to the nearest token amount.
        uint256 accountValue = positive.sub(negative).div(BaseMath.base());

        // Get the number of tokens in the Perpetual Contract.
        uint256 contractBalance = IERC20(_TOKEN_).balanceOf(address(this));

        // Determine the maximum withdrawable amount.
        uint256 amountToWithdraw = Math.min(contractBalance, accountValue);

        // Update the user's balance.
        uint120 remainingMargin = accountValue.sub(amountToWithdraw).toUint120();
        balance = P1Types.Balance({
            marginIsPositive: remainingMargin != 0,
            positionIsPositive: false,
            margin: remainingMargin,
            position: 0
        });
        _BALANCES_[msg.sender] = balance;

        // Send the tokens.
        SafeERC20.safeTransfer(
            IERC20(_TOKEN_),
            msg.sender,
            amountToWithdraw
        );

        // Emit the log.
        emit LogWithdrawFinalSettlement(
            msg.sender,
            amountToWithdraw,
            balance.toBytes32()
        );
    }
}
