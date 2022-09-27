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

import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/SafeERC20.sol";
import { I_ExchangeWrapper } from "../../../external/I_ExchangeWrapper.sol";
import { I_PerpetualV1 } from "../intf/I_PerpetualV1.sol";


/**
 * @title P1CurrencyConverterProxy
 * @author dYdX
 *
 * @notice Proxy contract which executes a trade via an ExchangeWrapper before making a deposit
 *  or after making a withdrawal.
 */
contract P1CurrencyConverterProxy {
    using SafeERC20 for IERC20;

    // ============ Events ============

    event LogConvertedDeposit(
        address indexed account,
        address source,
        address perpetual,
        address exchangeWrapper,
        address tokenFrom,
        address tokenTo,
        uint256 tokenFromAmount,
        uint256 tokenToAmount
    );

    event LogConvertedWithdrawal(
        address indexed account,
        address destination,
        address perpetual,
        address exchangeWrapper,
        address tokenFrom,
        address tokenTo,
        uint256 tokenFromAmount,
        uint256 tokenToAmount
    );

    // ============ State-Changing Functions ============

    /**
     * @notice Sets the maximum allowance on the Perpetual contract. Must be called at least once
     *  on a given Perpetual before deposits can be made.
     * @dev Cannot be run in the constructor due to technical restrictions in Solidity.
     */
    function approveMaximumOnPerpetual(
        address perpetual
    )
        external
    {
        IERC20 tokenContract = IERC20(I_PerpetualV1(perpetual).getTokenContract());

        // safeApprove requires unsetting the allowance first.
        tokenContract.safeApprove(perpetual, 0);

        // Set the allowance to the highest possible value.
        tokenContract.safeApprove(perpetual, uint256(-1));
    }

    /**
     * @notice Make a margin deposit to a Perpetual, after converting funds to the margin currency.
     *  Funds will be withdrawn from the sender and deposited into the specified account.
     * @dev Emits LogConvertedDeposit event.
     *
     * @param  account          The account for which to credit the deposit.
     * @param  perpetual        The PerpetualV1 contract to deposit to.
     * @param  exchangeWrapper  The ExchangeWrapper contract to trade with.
     * @param  tokenFrom        The token to convert from.
     * @param  tokenFromAmount  The amount of `tokenFrom` tokens to deposit.
     * @param  data             Trade parameters for the ExchangeWrapper.
     */
    function deposit(
        address account,
        address perpetual,
        address exchangeWrapper,
        address tokenFrom,
        uint256 tokenFromAmount,
        bytes calldata data
    )
        external
        returns (uint256)
    {
        I_PerpetualV1 perpetualContract = I_PerpetualV1(perpetual);
        address tokenTo = perpetualContract.getTokenContract();
        address self = address(this);

        // Send fromToken to the ExchangeWrapper.
        //
        // TODO: Take possible ERC20 fee into account.
        IERC20(tokenFrom).safeTransferFrom(
            msg.sender,
            exchangeWrapper,
            tokenFromAmount
        );

        // Convert fromToken to toToken on the ExchangeWrapper.
        I_ExchangeWrapper exchangeWrapperContract = I_ExchangeWrapper(exchangeWrapper);
        uint256 tokenToAmount = exchangeWrapperContract.exchange(
            msg.sender,
            self,
            tokenTo,
            tokenFrom,
            tokenFromAmount,
            data
        );

        // Receive toToken from the ExchangeWrapper.
        IERC20(tokenTo).safeTransferFrom(
            exchangeWrapper,
            self,
            tokenToAmount
        );

        // Deposit toToken to the Perpetual.
        perpetualContract.deposit(
            account,
            tokenToAmount
        );

        // Log the result.
        emit LogConvertedDeposit(
            account,
            msg.sender,
            perpetual,
            exchangeWrapper,
            tokenFrom,
            tokenTo,
            tokenFromAmount,
            tokenToAmount
        );

        return tokenToAmount;
    }

    /**
     * @notice Withdraw margin from a Perpetual, then convert the funds to another currency. Funds
     *  will be withdrawn from the specified account and transfered to the specified destination.
     * @dev Emits LogConvertedWithdrawal event.
     *
     * @param  account          The account to withdraw from.
     * @param  destination      The address to send the withdrawn funds to.
     * @param  perpetual        The PerpetualV1 contract to withdraw from to.
     * @param  exchangeWrapper  The ExchangeWrapper contract to trade with.
     * @param  tokenTo          The token to convert to.
     * @param  tokenFromAmount  The amount of `tokenFrom` tokens to withdraw.
     * @param  data             Trade parameters for the ExchangeWrapper.
     */
    function withdraw(
        address account,
        address destination,
        address perpetual,
        address exchangeWrapper,
        address tokenTo,
        uint256 tokenFromAmount,
        bytes calldata data
    )
        external
        returns (uint256)
    {
        I_PerpetualV1 perpetualContract = I_PerpetualV1(perpetual);
        address tokenFrom = perpetualContract.getTokenContract();
        address self = address(this);

        // Verify that the sender has permission to withdraw from the account.
        require(
            account == msg.sender || perpetualContract.hasAccountPermissions(account, msg.sender),
            "msg.sender cannot operate the account"
        );

        // Withdraw fromToken from the Perpetual.
        perpetualContract.withdraw(
            account,
            exchangeWrapper,
            tokenFromAmount
        );

        // Convert fromToken to toToken on the ExchangeWrapper.
        I_ExchangeWrapper exchangeWrapperContract = I_ExchangeWrapper(exchangeWrapper);
        uint256 tokenToAmount = exchangeWrapperContract.exchange(
            msg.sender,
            self,
            tokenTo,
            tokenFrom,
            tokenFromAmount,
            data
        );

        // Transfer toToken from the ExchangeWrapper to the destination address.
        IERC20(tokenTo).safeTransferFrom(
            exchangeWrapper,
            destination,
            tokenToAmount
        );

        // Log the result.
        emit LogConvertedWithdrawal(
            account,
            destination,
            perpetual,
            exchangeWrapper,
            tokenFrom,
            tokenTo,
            tokenFromAmount,
            tokenToAmount
        );

        return tokenToAmount;
    }
}
