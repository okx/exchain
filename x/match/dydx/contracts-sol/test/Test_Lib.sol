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

import { BaseMath } from "../protocol/lib/BaseMath.sol";
import { Math } from "../protocol/lib/Math.sol";
import { ReentrancyGuard } from "../protocol/lib/ReentrancyGuard.sol";
import { Require } from "../protocol/lib/Require.sol";
import { SafeCast } from "../protocol/lib/SafeCast.sol";
import { SignedMath } from "../protocol/lib/SignedMath.sol";
import { Storage } from "../protocol/lib/Storage.sol";
import { TypedSignature } from "../protocol/lib/TypedSignature.sol";
import { P1BalanceMath } from "../protocol/v1/lib/P1BalanceMath.sol";
import { P1Types } from "../protocol/v1/lib/P1Types.sol";


/**
 * @title Test_Lib
 * @author dYdX
 *
 * @notice Exposes library functions for testing.
 */
/* solium-disable-next-line camelcase */
contract Test_Lib is
    ReentrancyGuard
{

    // ============ BaseMath.sol ============

    function base()
        external
        pure
        returns (uint256)
    {
        return BaseMath.base();
    }

    function baseMul(
        uint256 value,
        uint256 baseValue
    )
        external
        pure
        returns (uint256)
    {
        return BaseMath.baseMul(value, baseValue);
    }

    function baseDivMul(
        uint256 value,
        uint256 baseValue
    )
        external
        pure
        returns (uint256)
    {
        return BaseMath.baseDivMul(value, baseValue);
    }

    function baseMulRoundUp(
        uint256 value,
        uint256 baseValue
    )
        external
        pure
        returns (uint256)
    {
        return BaseMath.baseMulRoundUp(value, baseValue);
    }

    function baseDiv(
        uint256 value,
        uint256 baseValue
    )
        external
        pure
        returns (uint256)
    {
        return BaseMath.baseDiv(value, baseValue);
    }

    function baseReciprocal(
        uint256 baseValue
    )
        external
        pure
        returns (uint256)
    {
        return BaseMath.baseReciprocal(baseValue);
    }

    // ============ Math.sol ============

    function getFraction(
        uint256 target,
        uint256 numerator,
        uint256 denominator
    )
        external
        pure
        returns (uint256)
    {
        return Math.getFraction(target, numerator, denominator);
    }

    function getFractionRoundUp(
        uint256 target,
        uint256 numerator,
        uint256 denominator
    )
        external
        pure
        returns (uint256)
    {
        return Math.getFractionRoundUp(target, numerator, denominator);
    }

    function min(
        uint256 a,
        uint256 b
    )
        external
        pure
        returns (uint256)
    {
        return Math.min(a, b);
    }

    function max(
        uint256 a,
        uint256 b
    )
        external
        pure
        returns (uint256)
    {
        return Math.max(a, b);
    }

    // ============ Require.sol ============

    function that(
        bool must,
        string calldata requireReason,
        address addr
    )
        external
        pure
    {
        Require.that(
            must,
            requireReason,
            addr
        );
    }

    // ============ SafeCast.sol ============

    function toUint128(
        uint256 value
    )
        external
        pure
        returns (uint128)
    {
        return SafeCast.toUint128(value);
    }

    function toUint120(
        uint256 value
    )
        external
        pure
        returns (uint120)
    {
        return SafeCast.toUint120(value);
    }

    function toUint32(
        uint256 value
    )
        external
        pure
        returns (uint32)
    {
        return SafeCast.toUint32(value);
    }

    // ============ SignedMath.sol ============

    function add(
        SignedMath.Int calldata sint,
        uint256 value
    )
        external
        pure
        returns (SignedMath.Int memory)
    {
        return SignedMath.add(sint, value);
    }

    function sub(
        SignedMath.Int calldata sint,
        uint256 value
    )
        external
        pure
        returns (SignedMath.Int memory)
    {
        return SignedMath.sub(sint, value);
    }

    function signedAdd(
        SignedMath.Int calldata augend,
        SignedMath.Int calldata addend
    )
        external
        pure
        returns (SignedMath.Int memory)
    {
        return SignedMath.signedAdd(augend, addend);
    }

    function signedSub(
        SignedMath.Int calldata minuend,
        SignedMath.Int calldata subtrahend
    )
        external
        pure
        returns (SignedMath.Int memory)
    {
        return SignedMath.signedSub(minuend, subtrahend);
    }

    // ============ Storage.sol ============

    function load(
        bytes32 slot
    )
        external
        view
        returns (bytes32)
    {
        return Storage.load(slot);
    }

    function store(
        bytes32 slot,
        bytes32 value
    )
        external
    {
        Storage.store(slot, value);
    }

    // ============ TypedSignature.sol ============

    function recover(
        bytes32 hash,
        bytes calldata signatureBytes
    )
        external
        pure
        returns (address)
    {
        TypedSignature.Signature memory signature = abi.decode(
            signatureBytes,
            (TypedSignature.Signature)
        );
        return TypedSignature.recover(hash, signature);
    }

    // ============ P1BalanceMath.sol ============

    function copy(
        P1Types.Balance calldata balance
    )
        external
        pure
        returns (P1Types.Balance memory)
    {
        return P1BalanceMath.copy(balance);
    }

    function addToMargin(
        P1Types.Balance calldata balance,
        uint256 amount
    )
        external
        pure
        returns (P1Types.Balance memory)
    {
        // Copy to memory, modify in place, and return the memory object.
        P1Types.Balance memory _balance = balance;
        P1BalanceMath.addToMargin(_balance, amount);
        return _balance;
    }

    function subFromMargin(
        P1Types.Balance calldata balance,
        uint256 amount
    )
        external
        pure
        returns (P1Types.Balance memory)
    {
        // Copy to memory, modify in place, and return the memory object.
        P1Types.Balance memory _balance = balance;
        P1BalanceMath.subFromMargin(_balance, amount);
        return _balance;
    }

    function addToPosition(
        P1Types.Balance calldata balance,
        uint256 amount
    )
        external
        pure
        returns (P1Types.Balance memory)
    {
        // Copy to memory, modify in place, and return the memory object.
        P1Types.Balance memory _balance = balance;
        P1BalanceMath.addToPosition(_balance, amount);
        return _balance;
    }

    function subFromPosition(
        P1Types.Balance calldata balance,
        uint256 amount
    )
        external
        pure
        returns (P1Types.Balance memory)
    {
        // Copy to memory, modify in place, and return the memory object.
        P1Types.Balance memory _balance = balance;
        P1BalanceMath.subFromPosition(_balance, amount);
        return _balance;
    }

    function getPositiveAndNegativeValue(
        P1Types.Balance calldata balance,
        uint256 price
    )
        external
        pure
        returns (uint256, uint256)
    {
        return P1BalanceMath.getPositiveAndNegativeValue(balance, price);
    }

    function getMargin(
        P1Types.Balance calldata balance
    )
        external
        pure
        returns (SignedMath.Int memory)
    {
        return P1BalanceMath.getMargin(balance);
    }

    function getPosition(
        P1Types.Balance calldata balance
    )
        external
        pure
        returns (SignedMath.Int memory)
    {
        return P1BalanceMath.getPosition(balance);
    }

    function setMargin(
        P1Types.Balance calldata balance,
        SignedMath.Int calldata newMargin
    )
        external
        pure
        returns (P1Types.Balance memory)
    {
        // Copy to memory, modify in place, and return the memory object.
        P1Types.Balance memory _balance = balance;
        P1BalanceMath.setMargin(_balance, newMargin);
        return _balance;
    }

    function setPosition(
        P1Types.Balance calldata balance,
        SignedMath.Int calldata newPosition
    )
        external
        pure
        returns (P1Types.Balance memory)
    {
        // Copy to memory, modify in place, and return the memory object.
        P1Types.Balance memory _balance = balance;
        P1BalanceMath.setPosition(_balance, newPosition);
        return _balance;
    }

    // ============ ReentrancyGuard.sol ============

    function nonReentrant1()
        public
        nonReentrant
        returns (uint256)
    {
        return this.nonReentrant2();
    }

    function nonReentrant2()
        public
        nonReentrant
        returns (uint256)
    {
        return 0;
    }
}
