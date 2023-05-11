// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract ContractA {

    address precomplieContarct = 0x0000000000000000000000000000000000000100;
    uint256 public number;
    event pushLog(string data);

    function callWasm(string memory wasmAddr, string memory msgData,bool requireASuccess) public payable returns (bytes memory response){
        number = number + 1;
        (bool success, bytes memory data) = precomplieContarct.call{value: msg.value} (
            abi.encodeWithSignature("callToWasm(string,string)", wasmAddr,msgData)
        );
        if (requireASuccess) {
            require(success);
            string memory res = abi.decode(data,(string));
            emit pushLog(res);
        }
        number = number + 1;
        return data;
    }

    function queryWasm(string memory msgData,bool requireASuccess) public payable returns (bytes memory response){
        number = number + 1;
        (bool success, bytes memory data) = precomplieContarct.call{value: msg.value} (
            abi.encodeWithSignature("queryToWasm(string)",msgData)
        );
        if (requireASuccess) {
            require(success);
            string memory res = abi.decode(data,(string));
            emit pushLog(res);
        }
        number = number + 1;
        return data;
    }

    function callToWasm(string memory wasmAddr, string memory data) public payable returns (string memory response) {
        return "";
    }

    function queryToWasm(string memory data) public view returns (string memory response) {
        return "";
    }
}

contract ContractB {
    uint256 public number;

    function callWasm(address contractA ,string memory wasmAddr, string memory msgData, bool requireASuccess,bool requireBSuccess) public payable returns (bytes memory response){
        number = number + 1;
        (bool success, bytes memory data) = contractA.call{value: msg.value} (
            abi.encodeWithSignature("callWasm(string,string,bool)", wasmAddr,msgData,requireASuccess)
        );
        number = number + 1;
        if (requireBSuccess) {
            require(success);
        }
        return data;
    }

    function queryWasm(address contractA , string memory msgData, bool requireASuccess,bool requireBSuccess) public payable returns (bytes memory response){
        number = number + 1;
        (bool success, bytes memory data) = contractA.call{value: msg.value} (
            abi.encodeWithSignature("queryWasm(string,bool)",msgData,requireASuccess)
        );
        number = number + 1;
        if (requireBSuccess) {
            require(success);
        }
        return data;
    }
}
