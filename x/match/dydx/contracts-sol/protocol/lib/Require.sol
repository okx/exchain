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


/**
 * @title Require
 * @author dYdX
 *
 * @dev Stringifies parameters to pretty-print revert messages.
 */
library Require {

    // ============ Constants ============

    uint256 constant ASCII_ZERO = 0x30; // '0'
    uint256 constant ASCII_RELATIVE_ZERO = 0x57; // 'a' - 10
    uint256 constant FOUR_BIT_MASK = 0xf;
    bytes23 constant ZERO_ADDRESS =
    0x3a20307830303030303030302e2e2e3030303030303030; // ": 0x00000000...00000000"

    // ============ Library Functions ============

    /**
     * @dev If the must condition is not true, reverts using a string combination of the reason and
     *  the address.
     */
    function that(
        bool must,
        string memory reason,
        address addr
    )
        internal
        pure
    {
        if (!must) {
            revert(string(abi.encodePacked(reason, stringify(addr))));
        }
    }

    // ============ Helper Functions ============

    /**
     * @dev Returns a bytes array that is an ASCII string representation of the input address.
     *  Returns " 0x", the first 4 bytes of the address in lowercase hex, "...", then the last 4
     *  bytes of the address in lowercase hex.
     */
    function stringify(
        address input
    )
        private
        pure
        returns (bytes memory)
    {
        // begin with ": 0x00000000...00000000"
        bytes memory result = abi.encodePacked(ZERO_ADDRESS);

        // initialize values
        uint256 z = uint256(input);
        uint256 shift1 = 8 * 20 - 4;
        uint256 shift2 = 8 * 4 - 4;

        // populate both sections in parallel
        for (uint256 i = 4; i < 12; i++) {
            result[i] = char(z >> shift1); // set char in first section
            result[i + 11] = char(z >> shift2); // set char in second section
            shift1 -= 4;
            shift2 -= 4;
        }

        return result;
    }

    /**
     * @dev Returns the ASCII hex character representing the last four bits of the input (0-9a-f).
     */
    function char(
        uint256 input
    )
        private
        pure
        returns (byte)
    {
        uint256 b = input & FOUR_BIT_MASK;
        return byte(uint8(b + ((b < 10) ? ASCII_ZERO : ASCII_RELATIVE_ZERO)));
    }
}
