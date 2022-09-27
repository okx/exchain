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

import { ERC20 } from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import { ERC20Detailed } from "@openzeppelin/contracts/token/ERC20/ERC20Detailed.sol";


/**
 * @title Test_Token2
 * @author dYdX
 *
 * @notice A second ERC-20 token for testing.
 */
/* solium-disable-next-line camelcase */
contract Test_Token2 is
    ERC20,
    ERC20Detailed("Test Token 2", "TEST2", 6)
{
    function mint(address account, uint256 amount) external {
        _mint(account, amount);
    }
}
