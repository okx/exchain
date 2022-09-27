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
import { P1TraderConstants } from "./P1TraderConstants.sol";
import { BaseMath } from "../../lib/BaseMath.sol";
import { TypedSignature } from "../../lib/TypedSignature.sol";
import { P1Getters } from "../impl/P1Getters.sol";
import { P1Types } from "../lib/P1Types.sol";


/**
 * @title P1InverseOrders
 * @author dYdX
 *
 * @notice Contract allowing trading between accounts using cryptographically signed messages.
 *  The limit and trigger prices are inverted, i.e. the base and quote currencies are swapped.
 *  This is to be used with inverse perpetuals, where the base currency is the margin currency used
 *  by the Perpetual smart contract, and the quote currency is the “position” currency.
 */
contract P1InverseOrders is
    P1TraderConstants
{
    using BaseMath for uint256;
    using SafeMath for uint256;

    // ============ Constants ============

    // EIP191 header for EIP712 prefix
    bytes2 constant private EIP191_HEADER = 0x1901;

    // EIP712 Domain Name value
    string constant private EIP712_DOMAIN_NAME = "P1InverseOrders";

    // EIP712 Domain Version value
    string constant private EIP712_DOMAIN_VERSION = "1.0";

    // Hash of the EIP712 Domain Separator Schema
    /* solium-disable-next-line indentation */
    bytes32 constant private EIP712_DOMAIN_SEPARATOR_SCHEMA_HASH = keccak256(abi.encodePacked(
        "EIP712Domain(",
        "string name,",
        "string version,",
        "uint256 chainId,",
        "address verifyingContract",
        ")"
    ));

    // Hash of the EIP712 LimitOrder struct
    /* solium-disable-next-line indentation */
    bytes32 constant private EIP712_ORDER_STRUCT_SCHEMA_HASH = keccak256(abi.encodePacked(
        "Order(",
        "bytes32 flags,",
        "uint256 amount,",
        "uint256 limitPrice,",
        "uint256 triggerPrice,",
        "uint256 limitFee,",
        "address maker,",
        "address taker,",
        "uint256 expiration",
        ")"
    ));

    // Bitmasks for the flags field
    bytes32 constant FLAG_MASK_NULL = bytes32(uint256(0));
    bytes32 constant FLAG_MASK_IS_BUY = bytes32(uint256(1));
    bytes32 constant FLAG_MASK_IS_DECREASE_ONLY = bytes32(uint256(1 << 1));
    bytes32 constant FLAG_MASK_IS_NEGATIVE_LIMIT_FEE = bytes32(uint256(1 << 2));

    // ============ Enums ============

    enum OrderStatus {
        Open,
        Approved,
        Canceled
    }

    // ============ Structs ============

    struct Order {
        bytes32 flags;
        uint256 amount;
        uint256 limitPrice;
        uint256 triggerPrice;
        uint256 limitFee;
        address maker;
        address taker;
        uint256 expiration;
    }

    struct Fill {
        uint256 amount;
        uint256 price;
        uint256 fee;
        bool isNegativeFee;
    }

    struct TradeData {
        Order order;
        Fill fill;
        TypedSignature.Signature signature;
    }

    struct OrderQueryOutput {
        OrderStatus status;
        uint256 filledAmount;
    }

    // ============ Events ============

    event LogOrderCanceled(
        address indexed maker,
        bytes32 orderHash
    );

    event LogOrderApproved(
        address indexed maker,
        bytes32 orderHash
    );

    event LogOrderFilled(
        bytes32 orderHash,
        bytes32 flags,
        uint256 triggerPrice,
        Fill fill
    );

    // ============ Immutable Storage ============

    // address of the perpetual contract
    address public _PERPETUAL_V1_;

    // Hash of the EIP712 Domain Separator data
    bytes32 public _EIP712_DOMAIN_HASH_;

    // ============ Mutable Storage ============

    // order hash => filled amount (in position amount)
    mapping (bytes32 => uint256) public _FILLED_AMOUNT_;

    // order hash => status
    mapping (bytes32 => OrderStatus) public _STATUS_;

    // ============ Constructor ============

    constructor (
        address perpetualV1,
        uint256 chainId
    )
        public
    {
        _PERPETUAL_V1_ = perpetualV1;

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
     * @notice Allows an account to take an order cryptographically signed by a different account.
     * @dev Emits the LogOrderFilled event.
     *
     * @param  sender  The address that called the trade() function on PerpetualV1.
     * @param  maker   The maker of the order.
     * @param  taker   The taker of the order.
     * @param  price   The current oracle price of the underlying asset.
     * @param  data    A struct of type TradeData.
     * @return         The assets to be traded and traderFlags that indicate that a trade occurred.
     */
    function trade(
        address sender,
        address maker,
        address taker,
        uint256 price,
        bytes calldata data,
        bytes32 /* traderFlags */
    )
        external
        returns (P1Types.TradeResult memory)
    {
        address perpetual = _PERPETUAL_V1_;

        require(
            msg.sender == perpetual,
            "msg.sender must be PerpetualV1"
        );

        if (taker != sender) {
            require(
                P1Getters(perpetual).hasAccountPermissions(taker, sender),
                "Sender does not have permissions for the taker"
            );
        }

        TradeData memory tradeData = abi.decode(data, (TradeData));
        bytes32 orderHash = _getOrderHash(tradeData.order);

        // validations
        _verifyOrderStateAndSignature(
            tradeData,
            orderHash
        );
        _verifyOrderRequest(
            tradeData,
            maker,
            taker,
            perpetual,
            price
        );

        // set _FILLED_AMOUNT_
        uint256 oldFilledAmount = _FILLED_AMOUNT_[orderHash];
        uint256 newFilledAmount = oldFilledAmount.add(tradeData.fill.amount);
        require(
            newFilledAmount <= tradeData.order.amount,
            "Cannot overfill order"
        );
        _FILLED_AMOUNT_[orderHash] = newFilledAmount;

        // Inverse perpetual: The fill amount is denoted in quote currency (position).
        emit LogOrderFilled(
            orderHash,
            tradeData.order.flags,
            tradeData.order.triggerPrice,
            tradeData.fill
        );

        // `isBuyOrder` is from the maker's perspective.
        bool isBuyOrder = _isBuy(tradeData.order);

        // Inverse perpetual:
        // - When isBuy is true, maker is buying base currency (margin) & selling quote (position).
        // - The fill amount is denoted in quote currency (position).
        // - Prices are denoted in quote currency per unit of base currency (position per margin).
        // - The fee is paid (or received) in base currency (margin).
        uint256 feeFactor = (isBuyOrder == tradeData.fill.isNegativeFee)
            ? BaseMath.base().add(tradeData.fill.fee)
            : BaseMath.base().sub(tradeData.fill.fee);
        // Note: Skip BaseMath since we multiply by a base value then divide by a base value.
        uint256 marginAmount = tradeData.fill.amount.mul(feeFactor).div(tradeData.fill.price);

        // Inverse perpetual: Note that `isBuy` in the order is from the maker's perspective and
        // relative to the base currency, whereas `isBuy` in the TradeResult is from the taker's
        // perspective, and relative to the quote currency.
        return P1Types.TradeResult({
            marginAmount: marginAmount,
            positionAmount: tradeData.fill.amount,
            isBuy: isBuyOrder,
            traderFlags: TRADER_FLAG_ORDERS
        });
    }

    /**
     * @notice On-chain approves an order.
     * @dev Emits the LogOrderApproved event.
     *
     * @param  order  The order that will be approved.
     */
    function approveOrder(
        Order calldata order
    )
        external
    {
        require(
            msg.sender == order.maker,
            "Order cannot be approved by non-maker"
        );
        bytes32 orderHash = _getOrderHash(order);
        require(
            _STATUS_[orderHash] != OrderStatus.Canceled,
            "Canceled order cannot be approved"
        );
        _STATUS_[orderHash] = OrderStatus.Approved;
        emit LogOrderApproved(msg.sender, orderHash);
    }

    /**
     * @notice On-chain cancels an order.
     * @dev Emits the LogOrderCanceled event.
     *
     * @param  order  The order that will be permanently canceled.
     */
    function cancelOrder(
        Order calldata order
    )
        external
    {
        require(
            msg.sender == order.maker,
            "Order cannot be canceled by non-maker"
        );
        bytes32 orderHash = _getOrderHash(order);
        _STATUS_[orderHash] = OrderStatus.Canceled;
        emit LogOrderCanceled(msg.sender, orderHash);
    }

    // ============ Getter Functions ============

    /**
     * @notice Gets the status (open/approved/canceled) and filled amount of each order in a list.
     *
     * @param  orderHashes  A list of the hashes of the orders to check.
     * @return              A list of OrderQueryOutput structs containing the status and filled
     *                      amount of each order.
     */
    function getOrdersStatus(
        bytes32[] calldata orderHashes
    )
        external
        view
        returns (OrderQueryOutput[] memory)
    {
        OrderQueryOutput[] memory result = new OrderQueryOutput[](orderHashes.length);
        for (uint256 i = 0; i < orderHashes.length; i++) {
            bytes32 orderHash = orderHashes[i];
            result[i] = OrderQueryOutput({
                status: _STATUS_[orderHash],
                filledAmount: _FILLED_AMOUNT_[orderHash]
            });
        }
        return result;
    }

    // ============ Helper Functions ============

    function _verifyOrderStateAndSignature(
        TradeData memory tradeData,
        bytes32 orderHash
    )
        private
        view
    {
        OrderStatus orderStatus = _STATUS_[orderHash];

        if (orderStatus == OrderStatus.Open) {
            require(
                tradeData.order.maker == TypedSignature.recover(orderHash, tradeData.signature),
                "Order has an invalid signature"
            );
        } else {
            require(
                orderStatus != OrderStatus.Canceled,
                "Order was already canceled"
            );
            assert(orderStatus == OrderStatus.Approved);
        }
    }

    function _verifyOrderRequest(
        TradeData memory tradeData,
        address maker,
        address taker,
        address perpetual,
        uint256 price
    )
        private
        view
    {
        require(
            tradeData.order.maker == maker,
            "Order maker does not match maker"
        );
        require(
            tradeData.order.taker == taker || tradeData.order.taker == address(0),
            "Order taker does not match taker"
        );
        require(
            tradeData.order.expiration >= block.timestamp || tradeData.order.expiration == 0,
            "Order has expired"
        );

        // `isBuyOrder` is from the maker's perspective.
        bool isBuyOrder = _isBuy(tradeData.order);
        bool validPrice = isBuyOrder
            ? tradeData.fill.price <= tradeData.order.limitPrice
            : tradeData.fill.price >= tradeData.order.limitPrice;
        require(
            validPrice,
            "Fill price is invalid"
        );

        bool validFee = _isNegativeLimitFee(tradeData.order)
            ? tradeData.fill.isNegativeFee && tradeData.fill.fee >= tradeData.order.limitFee
            : tradeData.fill.isNegativeFee || tradeData.fill.fee <= tradeData.order.limitFee;
        require(
            validFee,
            "Fill fee is invalid"
        );

        // Inverse perpetual: The trigger price must be compared against the inverted oracle price.
        uint256 invertedOraclePrice = price.baseReciprocal();
        if (tradeData.order.triggerPrice != 0) {
            bool validTriggerPrice = isBuyOrder
                ? tradeData.order.triggerPrice <= invertedOraclePrice
                : tradeData.order.triggerPrice >= invertedOraclePrice;
            require(
                validTriggerPrice,
                "Trigger price has not been reached"
            );
        }

        if (_isDecreaseOnly(tradeData.order)) {
            P1Types.Balance memory balance = P1Getters(perpetual).getAccountBalance(maker);
            // Inverse perpetual: Buying the base currency means selling position and selling the
            // base currency means buying position.
            require(
                isBuyOrder == balance.positionIsPositive
                && tradeData.fill.amount <= balance.position,
                "Fill does not decrease position"
            );
        }
    }

    /**
     * @dev Returns the EIP712 hash of an order.
     */
    function _getOrderHash(
        Order memory order
    )
        private
        view
        returns (bytes32)
    {
        // compute the overall signed struct hash
        /* solium-disable-next-line indentation */
        bytes32 structHash = keccak256(abi.encode(
            EIP712_ORDER_STRUCT_SCHEMA_HASH,
            order
        ));

        // compute eip712 compliant hash
        /* solium-disable-next-line indentation */
        return keccak256(abi.encodePacked(
            EIP191_HEADER,
            _EIP712_DOMAIN_HASH_,
            structHash
        ));
    }

    function _isBuy(
        Order memory order
    )
        private
        pure
        returns (bool)
    {
        return (order.flags & FLAG_MASK_IS_BUY) != FLAG_MASK_NULL;
    }

    function _isDecreaseOnly(
        Order memory order
    )
        private
        pure
        returns (bool)
    {
        return (order.flags & FLAG_MASK_IS_DECREASE_ONLY) != FLAG_MASK_NULL;
    }

    function _isNegativeLimitFee(
        Order memory order
    )
        private
        pure
        returns (bool)
    {
        return (order.flags & FLAG_MASK_IS_NEGATIVE_LIMIT_FEE) != FLAG_MASK_NULL;
    }
}
