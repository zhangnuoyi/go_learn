const { ethers } = require("hardhat");

async function main() {
  // 获取合约地址
  const marketplaceAddress = process.env.MARKETPLACE_ADDRESS;
  if (!marketplaceAddress) {
    console.error("请设置 MARKETPLACE_ADDRESS 环境变量");
    process.exit(1);
  }
  
  // 连接到合约
  const marketplace = await ethers.getContractAt("NFTMarketplace", marketplaceAddress);
  
  // 获取拍卖总数
  const auctionCount = await marketplace.auctionCount();
  console.log(`\n共有 ${auctionCount.toString()} 个拍卖`);
  console.log("-" .repeat(100));
  
  // 获取活跃的拍卖
  const activeAuctions = [];
  for (let i = 1; i <= auctionCount; i++) {
    try {
      const auction = await marketplace.getAuction(i);
      if (auction.isActive) {
        activeAuctions.push({ id: i, ...auction });
      }
    } catch (error) {
      console.log(`获取拍卖 ${i} 时出错: ${error.message}`);
    }
  }
  
  if (activeAuctions.length === 0) {
    console.log("当前没有活跃的拍卖");
    return;
  }
  
  console.log(`\n活跃拍卖 (${activeAuctions.length}):`);
  console.log("-" .repeat(100));
  console.log(`| ${"拍卖ID".padEnd(8)} | ${"卖家".padEnd(20)} | ${"NFT合约".padEnd(20)} | ${"TokenID".padEnd(10)} | ${"当前出价(ETH)".padEnd(15)} |`);
  console.log("-" .repeat(100));
  
  for (const item of activeAuctions) {
    const seller = item.seller.substring(0, 15) + "...";
    const nftContract = item.nftContract.substring(0, 15) + "...";
    const currentBid = ethers.utils.formatEther(item.currentBid);
    
    console.log(`| ${item.id.toString().padEnd(8)} | ${seller.padEnd(20)} | ${nftContract.padEnd(20)} | ${item.tokenId.toString().padEnd(10)} | ${currentBid.padEnd(15)} |`);
  }
  
  console.log("-" .repeat(100));
  console.log(`\n提示: 使用 npx hardhat run scripts/place-bid.js --network [网络名称] 进行出价`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("获取拍卖列表失败:", error);
    process.exit(1);
  });
