const { ethers, upgrades } = require("hardhat");

async function main() {
  const proxyAddress = process.env.MARKETPLACE_ADDRESS || "";
  
  if (!proxyAddress) {
    console.error("请提供NFTMarketplace代理合约地址");
    process.exit(1);
  }
  
  console.log(`Upgrading NFTMarketplace at ${proxyAddress}...`);
  
  const NFTMarketplaceV2 = await ethers.getContractFactory("NFTMarketplace");
  const upgraded = await upgrades.upgradeProxy(proxyAddress, NFTMarketplaceV2);
  
  console.log("NFTMarketplace upgraded successfully");
  console.log("New implementation address:", await upgrades.erc1967.getImplementationAddress(upgraded.address));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
