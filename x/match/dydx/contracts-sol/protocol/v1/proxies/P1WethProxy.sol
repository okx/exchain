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

import { WETH9 } from "canonical-weth/contracts/WETH9.sol";
import { P1Proxy } from "./P1Proxy.sol";
import { ReentrancyGuard } from "../../lib/ReentrancyGuard.sol";
import { I_PerpetualV1 } from "../intf/I_PerpetualV1.sol";


/**
 * @title P1WethProxy
 * @author dYdX
 *
 * @notice A proxy for depositing and withdrawing ETH to/from a Perpetual contract that uses WETH as
 *  its margin token. The ETH will be wrapper and unwrapped by the proxy.
 */
contract P1WethProxy is
    P1Proxy,
    ReentrancyGuard
{
    // ============ Storage ============

    WETH9 public _WETH_;

    // ============ Constructor ============

    constructor (
        address payable weth
    )
        public
    {
        _WETH_ = WETH9(weth);
    }

    // ============ External Functions ============

    /**
     * Fallback function. Disallows ether to be sent to this contract without data except when
     * unwrapping WETH.
     */
    function ()
        external
        payable
    {
        require(
            msg.sender == address(_WETH_),
            "Cannot receive ETH"
        );
    }

    /**
     * @notice Deposit ETH into a Perpetual, by first wrapping it as WETH. Any ETH paid to this
     *  function will be converted and deposited.
     *
     * @param  perpetual  Address of the Perpetual contract to deposit to.
     * @param  account    The account on the Perpetual for which to credit the deposit.
     */
    function depositEth(
        address perpetual,
        address account
    )
        external
        payable
        nonReentrant
    {
        WETH9 weth = _WETH_;
        address marginToken = I_PerpetualV1(perpetual).getTokenContract();
        require(
            marginToken == address(weth),
            "The perpetual does not use WETH for margin deposits"
        );

        // Wrap ETH.
        weth.deposit.value(msg.value)();

        // Deposit all WETH into the perpetual.
        uint256 amount = weth.balanceOf(address(this));
        I_PerpetualV1(perpetual).deposit(account, amount);
    }

    /**
     * @notice Withdraw ETH from a Perpetual, by first withdrawing and unwrapping WETH.
     *
     * @param  perpetual    Address of the Perpetual contract to withdraw from.
     * @param  account      The account on the Perpetual to withdraw from.
     * @param  destination  The address to send the withdrawn ETH to.
     * @param  amount       The amount of ETH/WETH to withdraw.
     */
    function withdrawEth(
        address perpetual,
        address account,
        address payable destination,
        uint256 amount
    )
        external
        nonReentrant
    {
        WETH9 weth = _WETH_;
        address marginToken = I_PerpetualV1(perpetual).getTokenContract();
        require(
            marginToken == address(weth),
            "The perpetual does not use WETH for margin deposits"
        );

        require(
            // Short-circuit if sender is the account owner.
            msg.sender == account ||
                I_PerpetualV1(perpetual).hasAccountPermissions(account, msg.sender),
            "Sender does not have withdraw permissions for the account"
        );

        // Withdraw WETH from the perpetual.
        I_PerpetualV1(perpetual).withdraw(account, address(this), amount);

        // Unwrap all WETH and send it as ETH to the provided destination.
        uint256 balance = weth.balanceOf(address(this));
        weth.withdraw(balance);
        destination.transfer(balance);
    }
}
