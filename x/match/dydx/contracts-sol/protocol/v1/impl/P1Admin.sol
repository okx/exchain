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

import { P1FinalSettlement } from "./P1FinalSettlement.sol";
import { P1Storage } from "./P1Storage.sol";
import { BaseMath } from "../../lib/BaseMath.sol";
import { I_P1Funder } from "../intf/I_P1Funder.sol";
import { I_P1Oracle } from "../intf/I_P1Oracle.sol";
import { P1Types } from "../lib/P1Types.sol";


/**
 * @title P1Admin
 * @author dYdX
 *
 * @notice Contract allowing the Admin address to set certain parameters.
 */
contract P1Admin is
    P1Storage,
    P1FinalSettlement
{
    // ============ Events ============

    event LogSetGlobalOperator(
        address operator,
        bool approved
    );

    event LogSetOracle(
        address oracle
    );

    event LogSetFunder(
        address funder
    );

    event LogSetMinCollateral(
        uint256 minCollateral
    );

    event LogFinalSettlementEnabled(
        uint256 settlementPrice
    );

    // ============ Functions ============

    /**
     * @notice Add or remove a Global Operator address.
     * @dev Must be called by the PerpetualV1 admin. Emits the LogSetGlobalOperator event.
     *
     * @param  operator  The address for which to enable or disable global operator privileges.
     * @param  approved  True if approved, false if disapproved.
     */
    function setGlobalOperator(
        address operator,
        bool approved
    )
        external
        onlyAdmin
        nonReentrant
    {
        _GLOBAL_OPERATORS_[operator] = approved;
        emit LogSetGlobalOperator(operator, approved);
    }

    /**
     * @notice Sets a new price oracle contract.
     * @dev Must be called by the PerpetualV1 admin. Emits the LogSetOracle event.
     *
     * @param  oracle  The address of the new price oracle contract.
     */
    function setOracle(
        address oracle
    )
        external
        onlyAdmin
        nonReentrant
    {
        require(
            I_P1Oracle(oracle).getPrice() != 0,
            "New oracle cannot return a zero price"
        );
        _ORACLE_ = oracle;
        emit LogSetOracle(oracle);
    }

    /**
     * @notice Sets a new funder contract.
     * @dev Must be called by the PerpetualV1 admin. Emits the LogSetFunder event.
     *
     * @param  funder  The address of the new funder contract.
     */
    function setFunder(
        address funder
    )
        external
        onlyAdmin
        nonReentrant
    {
        // call getFunding to ensure that no reverts occur
        I_P1Funder(funder).getFunding(0);

        _FUNDER_ = funder;
        emit LogSetFunder(funder);
    }

    /**
     * @notice Sets a new value for the minimum collateralization percentage.
     * @dev Must be called by the PerpetualV1 admin. Emits the LogSetMinCollateral event.
     *
     * @param  minCollateral  The new value of the minimum initial collateralization percentage,
     *                        as a fixed-point number with 18 decimals.
     */
    function setMinCollateral(
        uint256 minCollateral
    )
        external
        onlyAdmin
        nonReentrant
    {
        require(
            minCollateral >= BaseMath.base(),
            "The collateral requirement cannot be under 100%"
        );
        _MIN_COLLATERAL_ = minCollateral;
        emit LogSetMinCollateral(minCollateral);
    }

    /**
     * @notice Enables final settlement if the oracle price is between the provided bounds.
     * @dev Must be called by the PerpetualV1 admin. The current result of the price oracle
     *  must be between the two bounds supplied. Emits the LogFinalSettlementEnabled event.
     *
     * @param  priceLowerBound  The lower-bound (inclusive) of the acceptable price range.
     * @param  priceUpperBound  The upper-bound (inclusive) of the acceptable price range.
     */
    function enableFinalSettlement(
        uint256 priceLowerBound,
        uint256 priceUpperBound
    )
        external
        onlyAdmin
        noFinalSettlement
        nonReentrant
    {
        // Update the Global Index and grab the Price.
        P1Types.Context memory context = _loadContext();

        // Check price bounds.
        require(
            context.price >= priceLowerBound,
            "Oracle price is less than the provided lower bound"
        );
        require(
            context.price <= priceUpperBound,
            "Oracle price is greater than the provided upper bound"
        );

        // Save storage variables.
        _FINAL_SETTLEMENT_PRICE_ = context.price;
        _FINAL_SETTLEMENT_ENABLED_ = true;

        emit LogFinalSettlementEnabled(_FINAL_SETTLEMENT_PRICE_);
    }
}
