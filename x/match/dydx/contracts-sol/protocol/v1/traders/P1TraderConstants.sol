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
 * @title P1TraderConstants
 * @author dYdX
 *
 * @notice Constants for traderFlags set by contracts implementing the I_P1Trader interface.
 */
contract P1TraderConstants {
    bytes32 constant internal TRADER_FLAG_ORDERS = bytes32(uint256(1));
    bytes32 constant internal TRADER_FLAG_LIQUIDATION = bytes32(uint256(2));
    bytes32 constant internal TRADER_FLAG_DELEVERAGING = bytes32(uint256(4));
}
