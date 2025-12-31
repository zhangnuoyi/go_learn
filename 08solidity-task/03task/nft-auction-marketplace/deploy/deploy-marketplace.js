import { ethers, upgrades } from "hardhat";

async function main() {
    console.log("部署NFTMarketplace可升级合约...");
    
    const [deployer] = await ethers.getSigners();
    console.log("部署账户:", deployer.address);
    console.log("部署账户余额:", (await deployer.getBalance()).toString());
    
    // 部署NFTMarketplace可升级合约
    const NFTMarketplace = await ethers.getContractFactory("NFTMarketplace");
    console.log("准备部署代理合约...");
    
    const marketplace = await upgrades.deployProxy(NFTMarketplace, [], { 
        initializer: "initialize",
        kind: "uups"
    });
    
    await marketplace.deployed();
    
    console.log("NFTMarketplace代理合约已部署到:", marketplace.address);
    
    // 获取实现合约地址
    const implementationAddress = await upgrades.erc1967.getImplementationAddress(marketplace.address);
    console.log("NFTMarketplace实现合约地址:", implementationAddress);
    
    // 获取代理管理员地址
    const adminAddress = await upgrades.erc1967.getAdminAddress(marketplace.address);
    console.log("代理管理员地址:", adminAddress);
    
    console.log(`\n验证命令:`);
    console.log(`npx hardhat verify --network <network> ${implementationAddress}`);
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
