const { ethers, upgrades } = require("hardhat");

async function main() {
  console.log("Deploying AuctionNFT contract...");
  
  const AuctionNFT = await ethers.getContractFactory("AuctionNFT");
  const proxy = await upgrades.deployProxy(AuctionNFT, ["AuctionNFT", "ANFT"], {
    initializer: "initialize",
    kind: "uups"
  });
  
  await proxy.deployed();
  console.log("AuctionNFT deployed to:", proxy.address);
  console.log("Implementation address:", await upgrades.erc1967.getImplementationAddress(proxy.address));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
