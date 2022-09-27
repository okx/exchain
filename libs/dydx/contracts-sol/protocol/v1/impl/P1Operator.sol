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

import { P1Storage } from "./P1Storage.sol";


/**
 * @title P1Operator
 * @author dYdX
 *
 * @notice Contract for setting local operators for an account.
 */
contract P1Operator is
    P1Storage
{
    // ============ Events ============

    event LogSetLocalOperator(
        address indexed sender,
        address operator,
        bool approved
    );

    // ============ Functions ============

    /**
     * @notice Grants or revokes permission for another account to perform certain actions on behalf
     *  of the sender.
     * @dev Emits the LogSetLocalOperator event.
     *
     * @param  operator  The account that is approved or disapproved.
     * @param  approved  True for approval, false for disapproval.
     */
    function setLocalOperator(
        address operator,
        bool approved
    )
        external
    {
        _LOCAL_OPERATORS_[msg.sender][operator] = approved;
        emit LogSetLocalOperator(msg.sender, operator, approved);
    }
}
