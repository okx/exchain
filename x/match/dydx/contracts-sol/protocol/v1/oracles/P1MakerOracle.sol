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

import { Ownable } from "@openzeppelin/contracts/ownership/Ownable.sol";
import { I_MakerOracle } from "../../../external/maker/I_MakerOracle.sol";
import { BaseMath } from "../../lib/BaseMath.sol";
import { I_P1Oracle } from "../intf/I_P1Oracle.sol";


/**
 * @title P1MakerOracle
 * @author dYdX
 *
 * @notice P1Oracle that reads the price from a Maker V2 Oracle.
 */
contract P1MakerOracle is
    Ownable,
    I_P1Oracle
{
    using BaseMath for uint256;

    // ============ Events ============

    event LogRouteSet(
        address indexed sender,
        address oracle
    );

    event LogAdjustmentSet(
        address indexed oracle,
        uint256 adjustment
    );

    // ============ Storage ============

    // @dev Maps from the sender to the oracle address to use.
    mapping(address => address) public _ROUTER_;

    // @dev The amount to adjust the price by. Is as a fixed-point number with 18 decimal places.
    mapping(address => uint256) public _ADJUSTMENTS_;

    // ============ Public Functions ============

    /**
     * @notice Returns the price of the underlying asset relative to the margin token.
     *
     * @return The price as a fixed-point number with 18 decimals.
     */
    function getPrice()
        external
        view
        returns (uint256)
    {
        // get the oracle address to read from
        address oracle = _ROUTER_[msg.sender];

        // revert if no oracle found
        require(
            oracle != address(0),
            "Sender not authorized to get price"
        );

        // get adjustment or default to 1
        uint256 adjustment = _ADJUSTMENTS_[oracle];
        if (adjustment == 0) {
            adjustment = BaseMath.base();
        }

        // get the adjusted price
        uint256 rawPrice = uint256(I_MakerOracle(oracle).read());
        uint256 result = rawPrice.baseMul(adjustment);

        // revert if invalid price
        require(
            result != 0,
            "Oracle would return zero price"
        );

        return result;
    }

    // ============ Admin Functions ============

    /**
     * @dev Allows the owner to set a route for a particular sender.
     *
     * @param  sender The sender to set the route for.
     * @param  oracle The oracle to route the sender to.
     */
    function setRoute(
        address sender,
        address oracle
    )
        external
        onlyOwner
    {
        _ROUTER_[sender] = oracle;
        emit LogRouteSet(sender, oracle);
    }

    /**
     * @dev Allows the owner to set an adjustment to an oracle source.
     *
     * @param  oracle     The oracle to apply the adjustment to.
     * @param  adjustment The adjustment to set when reading from the oracle.
     */
    function setAdjustment(
        address oracle,
        uint256 adjustment
    )
        external
        onlyOwner
    {
        _ADJUSTMENTS_[oracle] = adjustment;
        emit LogAdjustmentSet(oracle, adjustment);
    }
}
