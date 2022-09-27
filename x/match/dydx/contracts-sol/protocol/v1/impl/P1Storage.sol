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

import { Adminable } from "../../lib/Adminable.sol";
import { ReentrancyGuard } from "../../lib/ReentrancyGuard.sol";
import { P1Types } from "../lib/P1Types.sol";


/**
 * @title P1Storage
 * @author dYdX
 *
 * @notice Storage contract. Contains or inherits from all contracts that have ordered storage.
 */
contract P1Storage is
    Adminable,
    ReentrancyGuard
{
    mapping(address => P1Types.Balance) internal _BALANCES_;
    mapping(address => P1Types.Index) internal _LOCAL_INDEXES_;

    mapping(address => bool) internal _GLOBAL_OPERATORS_;
    mapping(address => mapping(address => bool)) internal _LOCAL_OPERATORS_;

    address internal _TOKEN_;
    address internal _ORACLE_;
    address internal _FUNDER_;

    P1Types.Index internal _GLOBAL_INDEX_;
    uint256 internal _MIN_COLLATERAL_;

    bool internal _FINAL_SETTLEMENT_ENABLED_;
    uint256 internal _FINAL_SETTLEMENT_PRICE_;
}
