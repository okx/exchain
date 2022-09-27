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

import { P1FundingOracle } from "./P1FundingOracle.sol";


/**
 * @title P1InverseFundingOracle
 * @author dYdX
 *
 * @notice P1FundingOracle that uses the inverted rate (i.e. flips base and quote currencies)
 *  when getting the funding amount.
 */
contract P1InverseFundingOracle is
    P1FundingOracle
{
    // ============ Constructor ============

    constructor(
        address fundingRateProvider
    )
        P1FundingOracle(fundingRateProvider)
        public
    {
    }

    // ============ External Functions ============

    /**
     * @notice Calculates the signed funding amount that has accumulated over a period of time.
     *
     * @param  timeDelta  Number of seconds over which to calculate the accumulated funding amount.
     * @return            True if the funding rate is positive, and false otherwise.
     * @return            The funding amount as a unitless rate, represented as a fixed-point number
     *                    with 18 decimals.
     */
    function getFunding(
        uint256 timeDelta
    )
        public
        view
        returns (bool, uint256)
    {
        (bool isPositive, uint256 fundingAmount) = super.getFunding(timeDelta);
        return (!isPositive, fundingAmount);
    }
}
