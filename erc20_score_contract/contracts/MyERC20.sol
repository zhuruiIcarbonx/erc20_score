// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";


contract MyERC20 is ERC20,Ownable {

    constructor(string memory name_, string memory symbol_) ERC20(name_,symbol_) Ownable(msg.sender){
        
    }

    event Mint(address indexed to,uint256 value);
    
    event Burn(address indexed from,uint256 value);
    
    // event Transfer(address indexed from, address indexed to, uint256 value);

    function mint(address account, uint256 value) public onlyOwner{

       require(account != address(0)," account should not 0!");
       require(value >0," value should greater than 0!");
       _mint(account,value);
      emit Mint(account,value);

    }

    function burn(address account, uint256 value) public onlyOwner{

       require(account != address(0)," account should not 0!");
       require(value >0," value should greater than 0!");
       _burn(account,value);
      emit Burn(account,value);

    }

   //  function transfer(address from,address to, uint256 value) public{

   //     require(from != address(0)," from account should not 0!");
   //     require(to != address(0)," to account should not 0!");
   //     require(value >0," value should greater than 0!");
   //     _transfer(from,to,value);

   //  }

     



}