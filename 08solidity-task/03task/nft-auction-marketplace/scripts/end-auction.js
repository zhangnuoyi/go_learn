import { ethers } from "hardhat";
import * as dotenv from "dotenv";

dotenv.config();

async function main() {
  // 获取环境变量或使用命令行参数
  const marketplaceAddress = process.env.MARKETPLACE_ADDRESS || "";
  const auctionId = process.argv[2] ? parseInt(process.argv[2]) : 0;
  
  if (!marketplaceAddress) {
    console.error("请设置 MARKETPLACE_ADDRESS 环境变量");
    process.exit(1);
  }
  
  if (auctionId <= 0) {
    console.error("请提供有效的拍卖ID作为命令行参数");
    console.log("使用方法: node scripts/end-auction.js <auctionId>");
    process.exit(1);
  }
  
  console.log("准备结束拍卖...");
  console.log("市场合约地址:", marketplaceAddress);
  console.log("拍卖ID:", auctionId);
  
  const [caller] = await ethers.getSigners();
  console.log("调用者地址:", caller.address);
  
  // 获取合约实例
  const marketplaceABI = [
    "function endAuction(uint256 auctionId) public",
    "function _auctions(uint256) public view returns (uint256,address,address,uint256,uint256,uint256,address,uint256,address,bool)",
  ];
  
  const marketplace = new ethers.Contract(marketplaceAddress, marketplaceABI, caller);
  
  // 检查拍卖状态
  try {
    const auction = await marketplace._auctions(auctionId);
    const [id, seller, nftContract, tokenId, startTime, endTime, highestBidder, highestBid, paymentToken, ended] = auction;
    
    console.log("拍卖信息:");
    console.log("- ID:", id.toString());
    console.log("- 卖家:", seller);
    console.log("- 结束时间:", new Date(endTime.toNumber() * 1000).toLocaleString());
    console.log("- 当前最高出价:", ethers.utils.formatEther(highestBid), "ETH");
    console.log("- 最高出价者:", highestBidder);
    console.log("- 是否已结束:", ended ? "是" : "否");
    
    if (ended) {
      console.error("错误: 拍卖已经结束");
      process.exit(1);
    }
    
    const currentTime = Date.now() / 1000;
    if (currentTime < endTime.toNumber()) {
      console.warn("警告: 拍卖还未到结束时间，可能只有卖家可以结束");
    }
  } catch (error) {
    console.error("获取拍卖信息失败:", error);
    process.exit(1);
  }
  
  // 结束拍卖
  console.log("结束拍卖...");
  try {
    const tx = await marketplace.endAuction(auctionId);
    console.log("交易哈希:", tx.hash);
    const receipt = await tx.wait();
    console.log("拍卖结束成功!");
    
    // 检查事件
    for (const event of receipt.events || []) {
      if (event.event === "AuctionEnded") {
        console.log("拍卖结束事件:");
        console.log("- 拍卖ID:", event.args?.auctionId.toString());
        console.log("- 中标者:", event.args?.winner);
        console.log("- 中标金额:", ethers.utils.formatEther(event.args?.amount), "ETH");
        break;
      }
    }
  } catch (error) {
    console.error("结束拍卖失败:", error);
    process.exit(1);
  }
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
