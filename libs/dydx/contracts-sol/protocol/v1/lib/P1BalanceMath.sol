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
import { P1Types } from "./P1Types.sol";
import { BaseMath } from "../../lib/BaseMath.sol";
import { SafeCast } from "../../lib/SafeCast.sol";
import { SignedMath } from "../../lib/SignedMath.sol";


/**
 * @title P1BalanceMath
 * @author dYdX
 *
 * @dev Library for manipulating P1Types.Balance structs.
 */
library P1BalanceMath {
    using BaseMath for uint256;
    using SafeCast for uint256;
    using SafeMath for uint256;
    using SignedMath for SignedMath.Int;
    using P1BalanceMath for P1Types.Balance;

    // ============ Constants ============

    uint256 private constant FLAG_MARGIN_IS_POSITIVE = 1 << (8 * 31);
    uint256 private constant FLAG_POSITION_IS_POSITIVE = 1 << (8 * 15);

    // ============ Functions ============

    /**
     * @dev Create a copy of the balance struct.
     */
    function copy(
        P1Types.Balance memory balance
    )
        internal
        pure
        returns (P1Types.Balance memory)
    {
        return P1Types.Balance({
            marginIsPositive: balance.marginIsPositive,
            positionIsPositive: balance.positionIsPositive,
            margin: balance.margin,
            position: balance.position
        });
    }

    /**
     * @dev In-place add amount to balance.margin.
     */
    function addToMargin(
        P1Types.Balance memory balance,
        uint256 amount
    )
        internal
        pure
    {
        SignedMath.Int memory signedMargin = balance.getMargin();
        signedMargin = signedMargin.add(amount);
        balance.setMargin(signedMargin);
    }

    /**
     * @dev In-place subtract amount from balance.margin.
     */
    function subFromMargin(
        P1Types.Balance memory balance,
        uint256 amount
    )
        internal
        pure
    {
        SignedMath.Int memory signedMargin = balance.getMargin();
        signedMargin = signedMargin.sub(amount);
        balance.setMargin(signedMargin);
    }

    /**
     * @dev In-place add amount to balance.position.
     */
    function addToPosition(
        P1Types.Balance memory balance,
        uint256 amount
    )
        internal
        pure
    {
        SignedMath.Int memory signedPosition = balance.getPosition();
        signedPosition = signedPosition.add(amount);
        balance.setPosition(signedPosition);
    }

    /**
     * @dev In-place subtract amount from balance.position.
     */
    function subFromPosition(
        P1Types.Balance memory balance,
        uint256 amount
    )
        internal
        pure
    {
        SignedMath.Int memory signedPosition = balance.getPosition();
        signedPosition = signedPosition.sub(amount);
        balance.setPosition(signedPosition);
    }

    /**
     * @dev Returns the positive and negative values of the margin and position together, given a
     *  price, which is used as a conversion rate between the two currencies.
     *
     *  No rounding occurs here--the returned values are "base values" with extra precision.
     */
    function getPositiveAndNegativeValue(
        P1Types.Balance memory balance,
        uint256 price
    )
        internal
        pure
        returns (uint256, uint256)
    {
        uint256 positiveValue = 0;
        uint256 negativeValue = 0;

        // add value of margin
        if (balance.marginIsPositive) {
            positiveValue = uint256(balance.margin).mul(BaseMath.base());
        } else {
            negativeValue = uint256(balance.margin).mul(BaseMath.base());
        }

        // add value of position
        uint256 positionValue = uint256(balance.position).mul(price);
        if (balance.positionIsPositive) {
            positiveValue = positiveValue.add(positionValue);
        } else {
            negativeValue = negativeValue.add(positionValue);
        }

        return (positiveValue, negativeValue);
    }

    /**
     * @dev Returns a compressed bytes32 representation of the balance for logging.
     */
    function toBytes32(
        P1Types.Balance memory balance
    )
        internal
        pure
        returns (bytes32)
    {
        uint256 result =
            uint256(balance.position)
            | (uint256(balance.margin) << 128)
            | (balance.marginIsPositive ? FLAG_MARGIN_IS_POSITIVE : 0)
            | (balance.positionIsPositive ? FLAG_POSITION_IS_POSITIVE : 0);
        return bytes32(result);
    }

    // ============ Helper Functions ============

    /**
     * @dev Returns a SignedMath.Int version of the margin in balance.
     */
    function getMargin(
        P1Types.Balance memory balance
    )
        internal
        pure
        returns (SignedMath.Int memory)
    {
        return SignedMath.Int({
            value: balance.margin,
            isPositive: balance.marginIsPositive
        });
    }

    /**
     * @dev Returns a SignedMath.Int version of the position in balance.
     */
    function getPosition(
        P1Types.Balance memory balance
    )
        internal
        pure
        returns (SignedMath.Int memory)
    {
        return SignedMath.Int({
            value: balance.position,
            isPositive: balance.positionIsPositive
        });
    }

    /**
     * @dev In-place modify the signed margin value of a balance.
     */
    function setMargin(
        P1Types.Balance memory balance,
        SignedMath.Int memory newMargin
    )
        internal
        pure
    {
        balance.margin = newMargin.value.toUint120();
        balance.marginIsPositive = newMargin.isPositive;
    }

    /**
     * @dev In-place modify the signed position value of a balance.
     */
    function setPosition(
        P1Types.Balance memory balance,
        SignedMath.Int memory newPosition
    )
        internal
        pure
    {
        balance.position = newPosition.value.toUint120();
        balance.positionIsPositive = newPosition.isPositive;
    }
}
