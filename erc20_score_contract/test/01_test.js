const {  ethers,deployments,upgrades } = require("hardhat")
const { expect } = require("chai")
const fs = require("fs");
const path = require("path");
const dotenv = require('dotenv')

let cycle_number = 60;

describe("test MyERC20 ", async function() {


    it("should deploy  MyERC20 ", async function() {

        const { deployer,user1,user2,user3,user4,user5,user6,user7,user8,user9 } = await getNamedAccounts();
        const deployerSigner = await ethers.getSigner(deployer);
        const signer1 = await ethers.getSigner(user1);
        const signer2 = await ethers.getSigner(user2);
        const signer3 = await ethers.getSigner(user3);
        const signer4 = await ethers.getSigner(user4);
        const signer5 = await ethers.getSigner(user5);
        const signer6 = await ethers.getSigner(user6);
        const signer7 = await ethers.getSigner(user7);
        const signer8 = await ethers.getSigner(user8);
        const signer9 = await ethers.getSigner(user9);
        const userArrray = new Array(deployer,user1,user2,user3,user4,user5,user6,user7,user8,user9);
        const signerArrray = new Array(deployerSigner,signer1,signer2,signer3,signer4,signer5,signer6,signer7,signer8,signer9);

        console.log("userArrray:",userArrray)

        
      
        console.log("key----------",process.env['INFURA_API_KEY']) 

        //部署MemeToken合约
        await deployments.fixture(["deployMyERC20"])

        const erc20Data = await deployments.get("MyERC20")
        const erc20Contract = await ethers.getContractAt("MyERC20",erc20Data.address,deployerSigner)
        console.log("[test01]MyERC20合约地址：",erc20Data.address)

        console.log("[test01]开始************给deployer,user1 -- user9铸造10000币****************")
        for(i=0;i<10;i++){

            const mintTx = await erc20Contract.mint(userArrray[i],10000);
            const mintReceipt = await mintTx.wait()
            console.log("[test01]user"+i+"  balance is :",await erc20Contract.balanceOf(userArrray[i]))

        }
        console.log("[test01]完成************给deployer,user1 -- user9铸造10000币****************")

        
        console.log("[test01]开启:随机transfer模式 和 burn模式................每秒随机生成一次tranfer和burn")
        while(cycle_number>0){
            cycle_number = cycle_number-1;

            //生成随机数
            let min = 0;
            let max = 9;
            fromNum = Math.floor(Math.random() * (max - min + 1)) + min;
            toNum = Math.floor(Math.random() * (max - min + 1)) + min;
            burnNum = Math.floor(Math.random() * (max - min + 1)) + min;
            value = (Math.floor(Math.random() * (max - min + 1)) + min+1)*10;

            console.log("[test01][before]value:"+value)

            //tranfer
            if(fromNum != toNum){
                
                console.log("[test01][before]user:"+userArrray[fromNum]+"  balance is :",await erc20Contract.balanceOf(userArrray[fromNum]))
                console.log("[test01][before]user:"+userArrray[toNum]+"  balance is :",await erc20Contract.balanceOf(userArrray[toNum]))

                const mintTx = await erc20Contract.connect(signerArrray[fromNum]).transfer(userArrray[toNum],value);
                const mintReceipt = await mintTx.wait()
                console.log("[test01]event--- [tranfer]------from:",userArrray[fromNum],".......to:",userArrray[toNum],".........value:",value);
                
                console.log("[test01][after]user:"+userArrray[fromNum]+"  balance is :",await erc20Contract.balanceOf(userArrray[fromNum]))
                console.log("[test01][after]user:"+userArrray[toNum]+"  balance is :",await erc20Contract.balanceOf(userArrray[toNum]))

            }
            
            //burn
            if(fromNum == 0 && fromNum != burnNum){

                console.log("[test01][before]user:"+userArrray[burnNum]+"  balance is :",await erc20Contract.balanceOf(userArrray[burnNum]))
                
                const mintTx = await erc20Contract.connect(signerArrray[fromNum]).burn(userArrray[burnNum],value);
                const mintReceipt = await mintTx.wait()
                
                console.log("[test01]event--- [burn]------from:",userArrray[fromNum],".......to:",userArrray[toNum],".........value:",value);
                console.log("[test01][after]user:"+userArrray[burnNum]+"  balance is :",await erc20Contract.balanceOf(userArrray[burnNum]))

            }

            //休眠1s
            await new Promise(resolve => setTimeout(resolve, 500));


        }






    });



})