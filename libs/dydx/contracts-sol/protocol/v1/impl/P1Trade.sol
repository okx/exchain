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
import { P1FinalSettlement } from "./P1FinalSettlement.sol";
import { BaseMath } from "../../lib/BaseMath.sol";
import { Require } from "../../lib/Require.sol";
import { I_P1Trader } from "../intf/I_P1Trader.sol";
import { P1BalanceMath } from "../lib/P1BalanceMath.sol";
import { P1Types } from "../lib/P1Types.sol";


/**
 * @title P1Trade
 * @author dYdX
 *
 * @notice Contract for settling trades between two accounts. A "trade" in this context may refer
 *  to any approved transfer of balances, as determined by the smart contracts implementing the
 *  I_P1Trader interface and approved as global operators on the PerpetualV1 contract.
 */
contract P1Trade is
    P1FinalSettlement
{
    using SafeMath for uint120;
    using P1BalanceMath for P1Types.Balance;

    // ============ Structs ============

    struct TradeArg {
        uint256 takerIndex;
        uint256 makerIndex;
        address trader;
        bytes data;
    }

    // ============ Events ============

    event LogTrade(
        address indexed maker,
        address indexed taker,
        address trader,
        uint256 marginAmount,
        uint256 positionAmount,
        bool isBuy, // from taker's perspective
        bytes32 makerBalance,
        bytes32 takerBalance
    );

    // ============ Functions ============

    /**
     * @notice Submits one or more trades between any number of accounts.
     * @dev Emits the LogIndex event, one LogAccountSettled event for each account in `accounts`,
     *  and the LogTrade event for each trade in `trades`.
     *
     * @param  accounts  The sorted list of accounts that are involved in trades.
     * @param  trades    The list of trades to execute in-order.
     */
    function trade(
        address[] memory accounts,
        TradeArg[] memory trades
    )
        public
        noFinalSettlement
        nonReentrant
    {
        _verifyAccounts(accounts);
        P1Types.Context memory context = _loadContext();
        P1Types.Balance[] memory initialBalances = _settleAccounts(context, accounts);
        P1Types.Balance[] memory currentBalances = new P1Types.Balance[](initialBalances.length);

        uint256 i;
        for (i = 0; i < initialBalances.length; i++) {
            currentBalances[i] = initialBalances[i].copy();
        }

        bytes32 traderFlags = 0;
        for (i = 0; i < trades.length; i++) {
            TradeArg memory tradeArg = trades[i];

            require(
                _GLOBAL_OPERATORS_[tradeArg.trader],
                "trader is not global operator"
            );

            address maker = accounts[tradeArg.makerIndex];
            address taker = accounts[tradeArg.takerIndex];

            P1Types.TradeResult memory tradeResult = I_P1Trader(tradeArg.trader).trade(
                msg.sender,
                maker,
                taker,
                context.price,
                tradeArg.data,
                traderFlags
            );

            traderFlags |= tradeResult.traderFlags;

            // If the accounts are equal, no need to update balances.
            if (maker == taker) {
                continue;
            }

            // Modify currentBalances in-place. Note that `isBuy` is from the taker's perspective.
            P1Types.Balance memory makerBalance = currentBalances[tradeArg.makerIndex];
            P1Types.Balance memory takerBalance = currentBalances[tradeArg.takerIndex];
            if (tradeResult.isBuy) {
                makerBalance.addToMargin(tradeResult.marginAmount);
                makerBalance.subFromPosition(tradeResult.positionAmount);
                takerBalance.subFromMargin(tradeResult.marginAmount);
                takerBalance.addToPosition(tradeResult.positionAmount);
            } else {
                makerBalance.subFromMargin(tradeResult.marginAmount);
                makerBalance.addToPosition(tradeResult.positionAmount);
                takerBalance.addToMargin(tradeResult.marginAmount);
                takerBalance.subFromPosition(tradeResult.positionAmount);
            }

            // Store the new balances in storage.
            _BALANCES_[maker] = makerBalance;
            _BALANCES_[taker] = takerBalance;

            emit LogTrade(
                maker,
                taker,
                tradeArg.trader,
                tradeResult.marginAmount,
                tradeResult.positionAmount,
                tradeResult.isBuy,
                makerBalance.toBytes32(),
                takerBalance.toBytes32()
            );
        }

        _verifyAccountsFinalBalances(
            context,
            accounts,
            initialBalances,
            currentBalances
        );
    }

    /**
     * @dev Verify that `accounts` contains at least one address and that the contents are unique.
     *  We verify uniqueness by requiring that the array is sorted.
     */
    function _verifyAccounts(
        address[] memory accounts
    )
        private
        pure
    {
        require(
            accounts.length > 0,
            "Accounts must have non-zero length"
        );

        // Check that accounts are unique
        address prevAccount = accounts[0];
        for (uint256 i = 1; i < accounts.length; i++) {
            address account = accounts[i];
            require(
                account > prevAccount,
                "Accounts must be sorted and unique"
            );
            prevAccount = account;
        }
    }

    /**
     * Verify that account balances at the end of the tx are allowable given the initial balances.
     *
     * We require that for every account, either:
     * 1. The account meets the collateralization requirement; OR
     * 2. All of the following are true:
     *   a) The absolute value of the account position has not increased;
     *   b) The sign of the account position has not flipped positive to negative or vice-versa.
     *   c) The account's collateralization ratio has not worsened;
     */
    function _verifyAccountsFinalBalances(
        P1Types.Context memory context,
        address[] memory accounts,
        P1Types.Balance[] memory initialBalances,
        P1Types.Balance[] memory currentBalances
    )
        private
        pure
    {
        for (uint256 i = 0; i < accounts.length; i++) {
            P1Types.Balance memory currentBalance = currentBalances[i];
            (uint256 currentPos, uint256 currentNeg) =
                currentBalance.getPositiveAndNegativeValue(context.price);

            // See P1Settlement._isCollateralized().
            bool isCollateralized =
                currentPos.mul(BaseMath.base()) >= currentNeg.mul(context.minCollateral);

            if (isCollateralized) {
                continue;
            }

            address account = accounts[i];
            P1Types.Balance memory initialBalance = initialBalances[i];
            (uint256 initialPos, uint256 initialNeg) =
                initialBalance.getPositiveAndNegativeValue(context.price);

            Require.that(
                currentPos != 0,
                "account is undercollateralized and has no positive value",
                account
            );
            Require.that(
                currentBalance.position <= initialBalance.position,
                "account is undercollateralized and absolute position size increased",
                account
            );

            // Note that currentBalance.position can't be zero at this point since that would imply
            // either currentPos is zero or the account is well-collateralized.

            Require.that(
                currentBalance.positionIsPositive == initialBalance.positionIsPositive,
                "account is undercollateralized and position changed signs",
                account
            );
            Require.that(
                initialNeg != 0,
                "account is undercollateralized and was not previously",
                account
            );

            // Note that at this point:
            //   Absolute position size must have decreased and not changed signs.
            //   Initial margin/position must be one of -/-, -/+, or +/-.
            //   Current margin/position must now be either -/+ or +/-.
            //
            // Which implies one of the following [intial] -> [current] configurations:
            //   [-/-] -> [+/-]
            //   [-/+] -> [-/+]
            //   [+/-] -> [+/-]

            // Check that collateralization increased.
            // In the case of [-/-] initial, initialPos == 0 so the following will pass. Otherwise:
            // at this point, either initialNeg and currentNeg represent the margin values, or
            // initialPos and currentPos do. Since the margin is multiplied by the base value in
            // getPositiveAndNegativeValue(), it is safe to use baseDivMul() to divide the margin
            // without any rounding. This is important to avoid the possibility of overflow.
            Require.that(
                currentBalance.positionIsPositive
                    ? currentNeg.baseDivMul(initialPos) <= initialNeg.baseDivMul(currentPos)
                    : initialPos.baseDivMul(currentNeg) <= currentPos.baseDivMul(initialNeg),
                "account is undercollateralized and collateralization decreased",
                account
            );
        }
    }
}
