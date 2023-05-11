// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "./lib/JsonWriter.sol";
import "./lib/StringHelper.sol";

contract Counter is StringHelper {
    uint128 public count;

    using JsonWriter for JsonWriter.Json;

    event __OKCCallToWasm(string wasmAddr, uint256 value, string wasmMsg);

    function addCounterForWasm(
        string memory _wasmContractAddress,
        string memory delta
    ) public {
        //Assemble JSON data
        JsonWriter.Json memory _wasmMsg;

        _wasmMsg = _wasmMsg.writeStartObject();
        _wasmMsg = _wasmMsg.writeStartObject("add");
        _wasmMsg = _wasmMsg.writeStringProperty("delta", delta);
        _wasmMsg = _wasmMsg.writeEndObject();
        _wasmMsg = _wasmMsg.writeEndObject();

        //The specific event “__OKCCallToWasm” can trigger a wasm transaction
        emit __OKCCallToWasm(
            _wasmContractAddress, //wasm contract address(to)
            0, //The native token you want to send
            stringToHexString(_wasmMsg.value) //JSON => HexString
        );
    }

    function add(uint128 delta) public {
        count = count + delta;
    }
}
