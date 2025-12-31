const { ethers, upgrades } = require("hardhat");

async function main() {
  console.log("Deploying NFTMarketplace contract...");
  
  // Sepolia测试网的Chainlink价格预言机地址（示例）
  const priceFeeds = {
    // ETH/USD on Sepolia
    "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE": "0x694AA1769357215DE4FAC081bf1f309aDC325306",
    // 可以添加更多代币的价格预言机
  };
  
  const NFTMarketplace = await ethers.getContractFactory("NFTMarketplace");
  const proxy = await upgrades.deployProxy(NFTMarketplace, [10], {
    initializer: "initialize",
    kind: "uups"
  });
  
  await proxy.deployed();
  console.log("NFTMarketplace deployed to:", proxy.address);
  console.log("Implementation address:", await upgrades.erc1967.getImplementationAddress(proxy.address));
  
  // 配置价格预言机
  for (const [token, feed] of Object.entries(priceFeeds)) {
    const tx = await proxy.setPriceFeed(token, feed);
    await tx.wait();
    console.log(`Set price feed for ${token} to ${feed}`);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
