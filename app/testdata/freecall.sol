// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.18;
// pragma solidity >=0.7.0 <0.9.0;


contract FreeCall {
    address public constant moduleAddress =
    address(0x1033796B018B2bf0Fc9CB88c0793b2F275eDB624);

    event __OKCCallToWasm(string wasmAddr, uint256 value, string data);

    function callByWasm(string memory callerWasmAddr,string memory  data) public payable returns (string memory response) {
        string memory temp1 = strConcat("callByWasm return: ",callerWasmAddr);
        string memory temp2 = strConcat(temp1," ---data: ");
        string memory temp3 = strConcat(temp2,data);
        return temp3;
    }


    function callToWasm(string memory wasmAddr, uint256 value, string memory data) public returns (bool success){
        emit __OKCCallToWasm(wasmAddr,value,data);
        return true;
    }


    function strConcat(string memory _a, string  memory _b) internal returns (string memory){
        bytes memory _ba = bytes(_a);
        bytes memory _bb = bytes(_b);
        string memory ret = new string(_ba.length + _bb.length);
        bytes memory bret = bytes(ret);
        uint k = 0;
        for (uint i = 0; i < _ba.length; i++) {
            bret[k++] = _ba[i];
        }
        for (uint i = 0; i < _bb.length; i++) {
            bret[k++] = _bb[i];
        }
        return string(ret);
    }
}
