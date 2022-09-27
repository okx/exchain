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
import { P1TraderConstants } from "./P1TraderConstants.sol";
import { BaseMath } from "../../lib/BaseMath.sol";
import { Math } from "../../lib/Math.sol";
import { P1Getters } from "../impl/P1Getters.sol";
import { P1BalanceMath } from "../lib/P1BalanceMath.sol";
import { P1Types } from "../lib/P1Types.sol";


/**
 * @title P1Liquidation
 * @author dYdX
 *
 * @notice Contract allowing accounts to be liquidated by other accounts.
 */
contract P1Liquidation is
    P1TraderConstants
{
    using SafeMath for uint256;
    using Math for uint256;
    using P1BalanceMath for P1Types.Balance;

    // ============ Structs ============

    struct TradeData {
        uint256 amount;
        bool isBuy; // from taker's perspective
        bool allOrNothing; // if true, will revert if maker's position is less than the amount
    }

    // ============ Events ============

    event LogLiquidated(
        address indexed maker,
        address indexed taker,
        uint256 amount,
        bool isBuy, // from taker's perspective
        uint256 oraclePrice
    );

    // ============ Immutable Storage ============

    // address of the perpetual contract
    address public _PERPETUAL_V1_;

    // ============ Constructor ============

    constructor (
        address perpetualV1
    )
        public
    {
        _PERPETUAL_V1_ = perpetualV1;
    }

    // ============ External Functions ============

    /**
     * @notice Allows an account below the minimum collateralization to be liquidated by another
     *  account. This allows the account to be partially or fully subsumed by the liquidator.
     * @dev Emits the LogLiquidated event.
     *
     * @param  sender  The address that called the trade() function on PerpetualV1.
     * @param  maker   The account to be liquidated.
     * @param  taker   The account of the liquidator.
     * @param  price   The current oracle price of the underlying asset.
     * @param  data    A struct of type TradeData.
     * @return         The amounts to be traded, and flags indicating that a liquidation occurred.
     */
    function trade(
        address sender,
        address maker,
        address taker,
        uint256 price,
        bytes calldata data,
        bytes32 /* traderFlags */
    )
        external
        returns (P1Types.TradeResult memory)
    {
        address perpetual = _PERPETUAL_V1_;

        require(
            msg.sender == perpetual,
            "msg.sender must be PerpetualV1"
        );

        require(
            P1Getters(perpetual).getIsGlobalOperator(sender),
            "Sender is not a global operator"
        );

        TradeData memory tradeData = abi.decode(data, (TradeData));
        P1Types.Balance memory makerBalance = P1Getters(perpetual).getAccountBalance(maker);

        _verifyTrade(
            tradeData,
            makerBalance,
            perpetual,
            price
        );

        // Bound the execution amount by the size of the maker position.
        uint256 amount = Math.min(tradeData.amount, makerBalance.position);

        // When partially liquidating the maker, maintain the same position/margin ratio.
        // Ensure the collateralization of the maker does not decrease.
        uint256 marginAmount;
        if (tradeData.isBuy) {
            marginAmount = uint256(makerBalance.margin).getFractionRoundUp(
                amount,
                makerBalance.position
            );
        } else {
            marginAmount = uint256(makerBalance.margin).getFraction(amount, makerBalance.position);
        }

        emit LogLiquidated(
            maker,
            taker,
            amount,
            tradeData.isBuy,
            price
        );

        return P1Types.TradeResult({
            marginAmount: marginAmount,
            positionAmount: amount,
            isBuy: tradeData.isBuy,
            traderFlags: TRADER_FLAG_LIQUIDATION
        });
    }

    // ============ Helper Functions ============

    function _verifyTrade(
        TradeData memory tradeData,
        P1Types.Balance memory makerBalance,
        address perpetual,
        uint256 price
    )
        private
        view
    {
        require(
            _isUndercollateralized(makerBalance, perpetual, price),
            "Cannot liquidate since maker is not undercollateralized"
        );
        require(
            !tradeData.allOrNothing || makerBalance.position >= tradeData.amount,
            "allOrNothing is set and maker position is less than amount"
        );
        require(
            tradeData.isBuy == makerBalance.positionIsPositive,
            "liquidation must not increase maker's position size"
        );

        // Disallow liquidating in the edge case where both the position and margin are negative.
        //
        // This case is not handled correctly by P1Trade. If an account is in this situation, the
        // margin should first be set to zero via a deposit, then the account should be deleveraged.
        require(
            makerBalance.marginIsPositive || makerBalance.margin == 0 ||
                makerBalance.positionIsPositive || makerBalance.position == 0,
            "Cannot liquidate when maker position and margin are both negative"
        );
    }

    function _isUndercollateralized(
        P1Types.Balance memory balance,
        address perpetual,
        uint256 price
    )
        private
        view
        returns (bool)
    {
        uint256 minCollateral = P1Getters(perpetual).getMinCollateral();
        (uint256 positive, uint256 negative) = balance.getPositiveAndNegativeValue(price);

        // See P1Settlement.sol for discussion of overflow risk.
        return positive.mul(BaseMath.base()) < negative.mul(minCollateral);
    }
}
