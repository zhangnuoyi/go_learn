import { ethers } from "hardhat";

async function main() {
    console.log("部署MyNFT合约...");
    
    const [deployer] = await ethers.getSigners();
    console.log("部署账户:", deployer.address);
    console.log("部署账户余额:", (await deployer.getBalance()).toString());
    
    // 部署MyNFT合约
    const MyNFT = await ethers.getContractFactory("MyNFT");
    const myNFT = await MyNFT.deploy();
    
    await myNFT.deployed();
    
    console.log("MyNFT合约已部署到:", myNFT.address);
    
    // 输出验证命令
    console.log(`\n验证命令:`);
    console.log(`npx hardhat verify --network <network> ${myNFT.address}`);
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
