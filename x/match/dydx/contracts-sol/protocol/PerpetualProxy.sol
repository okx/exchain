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

/* solium-disable-next-line */
import { AdminUpgradeabilityProxy } from "@openzeppelin/upgrades/contracts/upgradeability/AdminUpgradeabilityProxy.sol";


/**
 * @title PerpetualProxy
 * @author dYdX
 *
 * @notice Proxy contract that forwards calls to the main Perpetual contract.
 */
contract PerpetualProxy is
    AdminUpgradeabilityProxy
{
    /**
     * @dev The constructor of the proxy that sets the admin and logic.
     *
     * @param  logic  The address of the contract that implements the underlying logic.
     * @param  admin  The address of the admin of the proxy.
     * @param  data   Any data to send immediately to the implementation contract.
     */
    constructor(
        address logic,
        address admin,
        bytes memory data
    )
        public
        AdminUpgradeabilityProxy(
            logic,
            admin,
            data
        )
    {}

    /**
     * @dev Overrides the default functionality that prevents the admin from reaching the
     *  implementation contract.
     */
    function _willFallback()
        internal
    { /* solium-disable-line no-empty-blocks */ }
}
