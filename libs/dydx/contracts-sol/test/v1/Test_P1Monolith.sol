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

import { Test_P1Funder } from "./Test_P1Funder.sol";
import { Test_P1Oracle } from "./Test_P1Oracle.sol";
import { Test_P1Trader } from "./Test_P1Trader.sol";


/**
 * @title Test_P1Monolith
 * @author dYdX
 *
 * @notice A second contract for testing the funder, oracle, and trader.
 */
/* solium-disable-next-line camelcase, no-empty-blocks */
contract Test_P1Monolith is
    Test_P1Funder,
    Test_P1Oracle,
    Test_P1Trader
{}
