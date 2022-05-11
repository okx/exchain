// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./ERC20.sol";


contract ModuleERC20 is ERC20  {
    address constant module_address = 0xc63cf6c8E1f3DF41085E9d8Af49584dae1432b4f;

    event __OkcSendToIbc(address sender, string recipient, uint256 amount);

    constructor(string memory denom_, uint8 decimals_) ERC20(denom_, denom_, decimals_) {}

    function native_denom() public view returns (string memory) {
        return symbol();
    }

    function mint_by_okc_module(address addr, uint amount) public {
        require(msg.sender == module_address);
        _mint(addr, amount);
    }

    function burn_by_okc_module(address addr, uint amount) public {
        require(msg.sender == module_address);
        _burn(addr, amount);
    }

    // send an "amount" of the contract token to recipient through IBC
    function send_to_ibc(string memory recipient, uint amount) public {
        _burn(msg.sender, amount);
        emit __OkcSendToIbc(msg.sender, recipient, amount);
    }
}
