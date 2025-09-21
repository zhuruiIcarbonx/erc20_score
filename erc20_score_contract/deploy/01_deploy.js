const { deployments, upgrades } = require("hardhat")

const fs = require("fs")//filesystem
const path = require("path")



module.exports = async ({getNamedAccounts, deployments}) => {

    //用户deployer、user1、user2、user3、user4信息
    const {save} = deployments;
    const {deployer, user1, user2, user3, user4} = await getNamedAccounts();
    const deployerSigner = await ethers.getSigner(deployer);


    
    //部署合约
    console.log("[01]开始部署MyERC20合约...")
    const factory = await ethers.getContractFactory("MyERC20", deployerSigner)
    const erc20 = await factory.deploy("MyERC20","MyERC20 token")
    await erc20.waitForDeployment();
    const erc20Address = erc20.target
    console.log("[01]MyERC20合约地址：", erc20Address)

       //保存合约地址
    await save("MyERC20", {
        abi:factory.interface.format("json"),
        address: erc20Address,
        // args:[],
        // log:true,    
    })

    
    const storePath = path.resolve(__dirname,"./.cache/MyERC20.json")
    fs.writeFileSync(
        storePath,
        JSON.stringify({
            address: erc20Address,
            abi: factory.interface.format("json")
        })
    )


}

module.exports.tags = ['deployMyERC20'];