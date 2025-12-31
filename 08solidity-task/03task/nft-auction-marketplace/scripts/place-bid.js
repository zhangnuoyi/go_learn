import { ethers } from "hardhat";
import * as dotenv from "dotenv";

dotenv.config();

async function main() {
  const [bidder] = await ethers.getSigners();
  console.log(`使用账户 ${bidder.address} 进行出价`);
  
  // 获取合约地址
  const marketplaceAddress = process.env.MARKETPLACE_ADDRESS;
  if (!marketplaceAddress) {
    console.error("请设置 MARKETPLACE_ADDRESS 环境变量");
    process.exit(1);
  }
  
  // 获取拍卖ID和出价金额
  const auctionId = process.argv[2];
  const bidAmount = process.argv[3];
  
  if (!auctionId || !bidAmount) {
    console.error("用法: npx hardhat run scripts/place-bid.js --network [网络名称] [拍卖ID] [出价金额(ETH)]");
    process.exit(1);
  }
  
  // 连接到合约
  const marketplace = await ethers.getContractAt("NFTMarketplace", marketplaceAddress);
  
  // 检查拍卖是否存在
  try {
    const auction = await marketplace.getAuction(auctionId);
    if (!auction.isActive) {
      console.error(`错误: 拍卖 ${auctionId} 不活跃或不存在`);
      process.exit(1);
    }
    
    console.log(`\n拍卖信息:`);
    console.log(`- 拍卖ID: ${auctionId}`);
    console.log(`- 当前最高出价: ${ethers.utils.formatEther(auction.currentBid)} ETH`);
    console.log(`- 最高出价者: ${auction.highestBidder}`);
    
    // 转换出价金额
    const bidAmountWei = ethers.utils.parseEther(bidAmount);
    
    // 检查出价金额
    if (bidAmountWei.le(auction.currentBid)) {
      console.error(`错误: 出价金额必须高于当前最高出价 ${ethers.utils.formatEther(auction.currentBid)} ETH`);
      process.exit(1);
    }
    
    // 执行出价
    console.log(`\n正在出价 ${bidAmount} ETH...`);
    const tx = await marketplace.placeBid(auctionId, { value: bidAmountWei });
    
    console.log(`交易已发送，等待确认...`);
    const receipt = await tx.wait();
    
    console.log(`\n✅ 出价成功！`);
    console.log(`交易哈希: ${receipt.transactionHash}`);
    console.log(`拍卖ID: ${auctionId}`);
    console.log(`出价金额: ${bidAmount} ETH`);
    
  } catch (error) {
    console.error(`出价失败: ${error.message}`);
    process.exit(1);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });