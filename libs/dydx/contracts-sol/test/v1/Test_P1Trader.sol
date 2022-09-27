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

import { I_P1Trader } from "../../protocol/v1/intf/I_P1Trader.sol";
import { P1Types } from "../../protocol/v1/lib/P1Types.sol";


/**
 * @title Test_P1Trader
 * @author dYdX
 *
 * @notice I_P1Trader implementation for testing.
 */
/* solium-disable-next-line camelcase */
contract Test_P1Trader is
    I_P1Trader
{
    P1Types.TradeResult public _TRADE_RESULT_;
    P1Types.TradeResult public _TRADE_RESULT_2_;

    // Special testing-only trader flag that will cause the second result to be returned.
    bytes32 constant public TRADER_FLAG_RESULT_2 = bytes32(~uint256(0));

    function trade(
        address, // sender
        address, // maker
        address, // taker
        uint256, // price
        bytes calldata, // data
        bytes32 traderFlags
    )
        external
        returns (P1Types.TradeResult memory)
    {
        if (traderFlags == TRADER_FLAG_RESULT_2) {
            return _TRADE_RESULT_2_;
        }
        return _TRADE_RESULT_;
    }

    function setTradeResult(
        uint256 marginAmount,
        uint256 positionAmount,
        bool isBuy,
        bytes32 traderFlags
    )
        external
    {
        _TRADE_RESULT_ = P1Types.TradeResult({
            marginAmount: marginAmount,
            positionAmount: positionAmount,
            isBuy: isBuy,
            traderFlags: traderFlags
        });
    }

    /**
     * Sets a second trade result which can be triggered by the trader flags of the first trade.
     */
    function setSecondTradeResult(
        uint256 marginAmount,
        uint256 positionAmount,
        bool isBuy,
        bytes32 traderFlags
    )
        external
    {
        _TRADE_RESULT_2_ = P1Types.TradeResult({
            marginAmount: marginAmount,
            positionAmount: positionAmount,
            isBuy: isBuy,
            traderFlags: traderFlags
        });
    }
}
