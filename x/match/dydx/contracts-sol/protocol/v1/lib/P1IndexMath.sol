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

import { P1Types } from "./P1Types.sol";


/**
 * @title P1IndexMath
 * @author dYdX
 *
 * @dev Library for manipulating P1Types.Index structs.
 */
library P1IndexMath {

    // ============ Constants ============

    uint256 private constant FLAG_IS_POSITIVE = 1 << (8 * 16);

    // ============ Functions ============

    /**
     * @dev Returns a compressed bytes32 representation of the index for logging.
     */
    function toBytes32(
        P1Types.Index memory index
    )
        internal
        pure
        returns (bytes32)
    {
        uint256 result =
            index.value
            | (index.isPositive ? FLAG_IS_POSITIVE : 0)
            | (uint256(index.timestamp) << 136);
        return bytes32(result);
    }
}
