pragma solidity ^0.6.8;

import "./token.sol";

contract ModuleERC20 is DSToken  {
    address constant module_address = 0x603871c2ddd41c26Ee77495E2E31e6De7f9957e0;
    string denom;

    event __OecSendToIbc(address sender, string recipient, uint256 amount);

    constructor(string memory denom_, uint8 decimals_) DSToken(denom_) public {
        decimals = decimals_;
        denom = denom_;
    }

    // unsafe_burn burn tokens without user's approval and authentication, used internally
    function unsafe_burn(address addr, uint amount) private {
        // Deduct user's balance without approval
        require(balanceOf[addr] >= amount, "ds-token-insufficient-balance");
        balanceOf[addr] = sub(balanceOf[addr], amount);
        totalSupply = sub(totalSupply, amount);
        emit Burn(addr, amount);
    }

    function native_denom() public view returns (string memory) {
        return denom;
    }

    function mint_by_oec_module(address addr, uint amount) public {
        require(msg.sender == module_address);
        mint(addr, amount);
    }

    function burn_by_oec_module(address addr, uint amount) public {
        require(msg.sender == module_address);
        unsafe_burn(addr, amount);
    }

    // send an "amount" of the contract token to recipient through IBC
    function send_to_ibc(string memory recipient, uint amount) public {
        unsafe_burn(msg.sender, amount);
        emit __OecSendToIbc(msg.sender, recipient, amount);
    }
}
