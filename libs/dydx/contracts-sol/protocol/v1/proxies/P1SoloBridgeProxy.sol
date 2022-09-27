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

import { SafeMath } from "@openzeppelin/contracts/math/SafeMath.sol";
import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/SafeERC20.sol";
import { P1Proxy } from "./P1Proxy.sol";
import { I_Solo } from "../../../external/dydx/I_Solo.sol";
import { BaseMath } from "../../lib/BaseMath.sol";
import { ReentrancyGuard } from "../../lib/ReentrancyGuard.sol";
import { SignedMath } from "../../lib/SignedMath.sol";
import { TypedSignature } from "../../lib/TypedSignature.sol";
import { I_PerpetualV1 } from "../intf/I_PerpetualV1.sol";
import { P1BalanceMath } from "../lib/P1BalanceMath.sol";
import { P1Types } from "../lib/P1Types.sol";


/**
 * @title P1SoloBridgeProxy
 * @author dYdX
 *
 * @notice Facilitates transfers between the PerpetualV1 and Solo smart contracts.
 */
contract P1SoloBridgeProxy is
    P1Proxy,
    ReentrancyGuard
{
    using BaseMath for uint256;
    using SafeMath for uint256;
    using SignedMath for SignedMath.Int;
    using P1BalanceMath for P1Types.Balance;
    using SafeERC20 for IERC20;

    // ============ Constants ============

    // EIP191 header for EIP712 prefix
    bytes2 constant private EIP191_HEADER = 0x1901;

    // EIP712 Domain Name value
    string constant private EIP712_DOMAIN_NAME = "P1SoloBridgeProxy";

    // EIP712 Domain Version value
    string constant private EIP712_DOMAIN_VERSION = "1.0";

    // EIP712 hash of the Domain Separator Schema
    /* solium-disable-next-line indentation */
    bytes32 constant private EIP712_DOMAIN_SEPARATOR_SCHEMA_HASH = keccak256(abi.encodePacked(
        "EIP712Domain(",
        "string name,",
        "string version,",
        "uint256 chainId,",
        "address verifyingContract",
        ")"
    ));

    // EIP712 hash of the Transfer struct
    /* solium-disable-next-line indentation */
    bytes32 constant private EIP712_TRANSFER_STRUCT_SCHEMA_HASH = keccak256(abi.encodePacked(
        "Transfer(",
        "address account,",
        "address perpetual,",
        "uint256 soloAccountNumber,",
        "uint256 soloMarketId,",
        "uint256 amount,",
        "bytes32 options",
        ")"
    ));

    // Constants for the options field of the Transfer struct.
    bytes32 private constant OPTIONS_MASK_TRANSFER_MODE = bytes32(uint256((1 << 8) - 1));
    bytes32 private constant OPTIONS_MASK_EXPIRATION = bytes32(uint256((1 << 120) - 1));
    uint256 private constant OPTIONS_OFFSET_EXPIRATION = 8;

    // ============ Enums ============

    enum TransferMode {
        SomeToPerpetual,
        SomeToSolo,
        AllToPerpetual
    }

    // ============ Structs ============

    struct Transfer {
        address account;
        address perpetual;
        uint256 soloAccountNumber;
        uint256 soloMarketId;
        uint256 amount;
        bytes32 options; // salt (16 bytes), expiration (15 bytes), transfer mode (1 byte)
    }

    // ============ Events ============

    event LogTransferred(
        address indexed account,
        address perpetual,
        uint256 soloAccountNumber,
        uint256 soloMarketId,
        bool toPerpetual,
        uint256 amount
    );

    event LogSignatureInvalidated(
        address indexed account,
        bytes32 transferHash
    );

    // ============ Immutable Storage ============

    // Address of the Solo margin contract.
    address public _SOLO_MARGIN_;

    // Hash of the EIP712 Domain Separator data
    bytes32 public _EIP712_DOMAIN_HASH_;

    // ============ Mutable Storage ============

    // transfer hash => bool
    mapping (bytes32 => bool) public _SIGNATURE_USED_;

    // ============ Constructor ============

    constructor (
        address soloMargin,
        uint256 chainId
    )
        public
    {
        _SOLO_MARGIN_ = soloMargin;

        /* solium-disable-next-line indentation */
        _EIP712_DOMAIN_HASH_ = keccak256(abi.encode(
            EIP712_DOMAIN_SEPARATOR_SCHEMA_HASH,
            keccak256(bytes(EIP712_DOMAIN_NAME)),
            keccak256(bytes(EIP712_DOMAIN_VERSION)),
            chainId,
            address(this)
        ));
    }

    // ============ External Functions ============

    /**
     * @notice Sets the maximum allowance on the Solo contract for a given market. Must be called
     *  at least once on a given market before deposits can be made.
     * @dev Cannot be run in the constructor due to technical restrictions in Solidity.
     */
    function approveMaximumOnSolo(
        uint256 soloMarketId
    )
        external
    {
        address solo = _SOLO_MARGIN_;
        IERC20 tokenContract = IERC20(I_Solo(solo).getMarketTokenAddress(soloMarketId));

        // safeApprove requires unsetting the allowance first.
        tokenContract.safeApprove(solo, 0);

        // Set the allowance to the highest possible value.
        tokenContract.safeApprove(solo, uint256(-1));
    }

    /**
     * @notice Executes a transfer from Solo to Perpetual or vice vera.
     * @dev Emits the LogTransferred event.
     *
     * @param  transfer   The transfer to execute.
     * @param  signature  (Optional) Signature for the transfer, will be checked if sender does not
     *                    have permission to operate the account.
     */
    function bridgeTransfer(
        Transfer calldata transfer,
        TypedSignature.Signature calldata signature
    )
        external
        nonReentrant
        returns (uint256)
    {
        bytes32 transferHash = _getTransferHash(transfer);
        I_Solo solo = I_Solo(_SOLO_MARGIN_);
        I_PerpetualV1 perpetual = I_PerpetualV1(transfer.perpetual);
        address tokenAddress = perpetual.getTokenContract();
        TransferMode transferMode = _getTransferMode(transfer);
        bool toPerpetual = (
            transferMode == TransferMode.SomeToPerpetual ||
            transferMode == TransferMode.AllToPerpetual
        );

        // Validations.
        _verifySoloMarket(
            solo,
            transfer,
            tokenAddress
        );

        // Verify permissions: Either msg.sender must be a valid operator of the account or the
        // signature must be valid.
        bool isValidOperator = _isValidOperator(
            solo,
            perpetual,
            transfer,
            toPerpetual
        );
        if (!isValidOperator) {
            _verifySignature(
                transfer,
                transferHash,
                signature
            );
            _SIGNATURE_USED_[transferHash] = true;
        }

        // Execute the transfer.
        uint256 amount;
        if (toPerpetual) {

            // Withdraw from Solo.
            uint256 initialBalance = IERC20(tokenAddress).balanceOf(address(this));
            _doSoloOperation(
                solo,
                transfer,
                true,
                transferMode == TransferMode.AllToPerpetual
            );
            uint256 finalBalance = IERC20(tokenAddress).balanceOf(address(this));

            // Deposit to Perpetual.
            amount = finalBalance.sub(initialBalance);
            perpetual.deposit(transfer.account, amount);
        } else { // transferMode == TransferMode.SomeToSolo

            // Withdraw from Perpetual.
            amount = transfer.amount;
            perpetual.withdraw(transfer.account, address(this), amount);

            // Deposit to Solo.
            _doSoloOperation(
                solo,
                transfer,
                false,
                false
            );
        }

        // Log the transfer.
        emit LogTransferred(
            transfer.account,
            transfer.perpetual,
            transfer.soloAccountNumber,
            transfer.soloMarketId,
            toPerpetual,
            amount
        );
    }

    /**
     * @notice Invalidate a signature, given the exact transfer parameters.
     * @dev Emits the LogSignatureInvalidated event.
     *
     * @param  transfer  The parameters for the signature that will be invalidated.
     */
    function invalidateSignature(
        Transfer calldata transfer
    )
        external
    {
        // Check permissions. Short-circuit if sender is the account owner.
        if (msg.sender != transfer.account) {
            I_Solo solo = I_Solo(_SOLO_MARGIN_);
            I_PerpetualV1 perpetual = I_PerpetualV1(transfer.perpetual);
            TransferMode transferMode = _getTransferMode(transfer);
            bool toPerpetual = (
                transferMode == TransferMode.SomeToPerpetual ||
                transferMode == TransferMode.AllToPerpetual
            );
            require(
                _isValidOperator(
                    solo,
                    perpetual,
                    transfer,
                    toPerpetual
                ),
                "Sender does not have permission to invalidate"
            );
        }

        // Mark this signature as used to prevent replay attacks.
        bytes32 transferHash = _getTransferHash(transfer);
        _SIGNATURE_USED_[transferHash] = true;

        // Log the invalidation.
        emit LogSignatureInvalidated(
            transfer.account,
            transferHash
        );
    }

    // ============ Helper Functions ============

    /**
     * @dev Execute a withdrawal or deposit operation on Solo.
     */
    function _doSoloOperation(
        I_Solo solo,
        Transfer memory transfer,
        bool isWithdrawal,
        bool withdrawToZero
    )
        private
    {
        // Create Solo accounts array.
        I_Solo.AccountInfo[] memory soloAccounts = new I_Solo.AccountInfo[](1);
        soloAccounts[0] = I_Solo.AccountInfo({
            owner: transfer.account,
            number: transfer.soloAccountNumber
        });

        // Create Solo actions array.
        I_Solo.AssetAmount memory amount = I_Solo.AssetAmount({
            sign: !isWithdrawal,
            denomination: I_Solo.AssetDenomination.Wei,
            ref: withdrawToZero ? I_Solo.AssetReference.Target : I_Solo.AssetReference.Delta,
            value: withdrawToZero ? 0 : transfer.amount
        });
        I_Solo.ActionArgs[] memory soloActions = new I_Solo.ActionArgs[](1);
        soloActions[0] = I_Solo.ActionArgs({
            actionType: isWithdrawal ? I_Solo.ActionType.Withdraw : I_Solo.ActionType.Deposit,
            accountId: transfer.soloAccountNumber,
            amount: amount,
            primaryMarketId: transfer.soloMarketId,
            secondaryMarketId: 0,
            otherAddress: address(this),
            otherAccountId: 0,
            data: ""
        });

        // Execute the withdrawal or deposit.
        solo.operate(soloAccounts, soloActions);
    }

    /**
     * Verify that the signature is valid for the given hash and not used, invalidated, or expired.
     */
    function _verifySignature(
        Transfer memory transfer,
        bytes32 transferHash,
        TypedSignature.Signature memory signature
    )
        private
        view
    {
        // Verify expiration.
        uint256 expiration = _getExpiration(transfer);
        require(
            expiration >= block.timestamp || expiration == 0,
            "Signature has expired"
        );

        // Check whether the signature was previously used or invalidated.
        require(
            !_SIGNATURE_USED_[transferHash],
            "Signature was already used or invalidated"
        );

        require(
            TypedSignature.recover(transferHash, signature) == transfer.account,
            "Sender does not have account permissions and signature is invalid"
        );
    }

    /**
     * Check whether msg.sender has permission to operate the account.
     */
    function _isValidOperator(
        I_Solo solo,
        I_PerpetualV1 perpetual,
        Transfer memory transfer,
        bool toPerpetual
    )
        private
        view
        returns (bool)
    {
        // Short-circuit if sender is the account owner.
        if (msg.sender == transfer.account) {
            return true;
        }

        bool isSoloOperator =
            solo.getIsGlobalOperator(msg.sender) ||
            solo.getIsLocalOperator(transfer.account, msg.sender);

        // Solo only accepts deposits from valid operators of an account, so require Solo operator
        // permissions for either transfer direction. Only require Perpetual operator permissions to
        // withdraw from Perpetual.
        return
            isSoloOperator &&
            (toPerpetual || perpetual.hasAccountPermissions(transfer.account, msg.sender));
    }

    /**
     * Verify token addresses.
     */
    function _verifySoloMarket(
        I_Solo solo,
        Transfer memory transfer,
        address tokenAddress
    )
        private
        view
    {
        // Verify that the Solo market asset matches the Perpetual margin asset.
        require(
            solo.getMarketTokenAddress(transfer.soloMarketId) == tokenAddress,
            "Solo and Perpetual assets are not the same"
        );
    }

    /**
     * Returns the EIP712 hash of a transfer.
     */
    function _getTransferHash(
        Transfer memory transfer
    )
        private
        view
        returns (bytes32)
    {
        // Compute the overall signed struct hash
        /* solium-disable-next-line indentation */
        bytes32 structHash = keccak256(abi.encode(
            EIP712_TRANSFER_STRUCT_SCHEMA_HASH,
            transfer
        ));

        // Compute EIP712 compliant hash
        /* solium-disable-next-line indentation */
        return keccak256(abi.encodePacked(
            EIP191_HEADER,
            _EIP712_DOMAIN_HASH_,
            structHash
        ));
    }

    function _getTransferMode(
        Transfer memory transfer
    )
        private
        pure
        returns (TransferMode)
    {
        return TransferMode(uint256(transfer.options & OPTIONS_MASK_TRANSFER_MODE));
    }

    function _getExpiration(
        Transfer memory transfer
    )
        private
        pure
        returns (uint256)
    {
        return uint256((transfer.options >> OPTIONS_OFFSET_EXPIRATION) & OPTIONS_MASK_EXPIRATION);
    }
}
