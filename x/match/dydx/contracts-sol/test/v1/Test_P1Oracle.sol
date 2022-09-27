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

import { I_P1Oracle } from "../../protocol/v1/intf/I_P1Oracle.sol";


/**
 * @title Test_P1Oracle
 * @author dYdX
 *
 * @notice I_P1Oracle implementation for testing.
 */
/* solium-disable-next-line camelcase */
contract Test_P1Oracle is
    I_P1Oracle
{
    uint256 public _PRICE_ = 0;

    function getPrice()
        external
        view
        returns (uint256)
    {
        return _PRICE_;
    }

    function setPrice(
        uint256 newPrice
    )
        external
    {
        _PRICE_ = newPrice;
    }
}
