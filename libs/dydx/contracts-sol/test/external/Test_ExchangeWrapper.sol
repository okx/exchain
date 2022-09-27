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
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/SafeERC20.sol";
import { I_ExchangeWrapper } from "../../external/I_ExchangeWrapper.sol";


/**
 * @title Test_ExchangeWrapper
 * @author dYdX
 *
 * ExchangeWrapper for testing.
 */
/* solium-disable-next-line camelcase */
contract Test_ExchangeWrapper is
    I_ExchangeWrapper
{
    using SafeERC20 for IERC20;

    // ============ Constants ============

    // Arbitrary address to send tokens to (they are burned).
    address public constant EXCHANGE_ADDRESS = address(0x1);

    // ============ Structs ============

    struct Order {
        uint256 amount;
    }

    // ============ Test Data ============

    uint256 public _MAKER_AMOUNT_;
    uint256 public _TAKER_AMOUNT_;

    // ============ Test Data Setter Functions ============

    function setMakerAmount(
        uint256 makerAmount
    )
        external
    {
        _MAKER_AMOUNT_ = makerAmount;
    }

    function setTakerAmount(
        uint256 takerAmount
    )
        external
    {
        _TAKER_AMOUNT_ = takerAmount;
    }

    // ============ Getter Functions ============

    function getExchangeCost(
        address /* makerToken */,
        address /* takerToken */,
        uint256 /* desiredMakerToken */,
        bytes calldata /* orderData */
    )
        external
        view
        returns (uint256)
    {
        return _TAKER_AMOUNT_;
    }

    // ============ State-Changing Functions ============

    function exchange(
        address /* tradeOriginator */,
        address receiver,
        address makerToken,
        address takerToken,
        uint256 requestedFillAmount,
        bytes calldata orderData
    )
        external
        returns (uint256)
    {
        Order memory order = abi.decode(orderData, (Order));

        require(
            order.amount == requestedFillAmount,
            "amount mistmatch"
        );

        IERC20(takerToken).safeTransfer(EXCHANGE_ADDRESS, requestedFillAmount);
        IERC20(makerToken).safeApprove(receiver, _MAKER_AMOUNT_);

        return _MAKER_AMOUNT_;
    }
}
