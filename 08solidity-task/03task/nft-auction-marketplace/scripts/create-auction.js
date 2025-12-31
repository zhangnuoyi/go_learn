import { ethers } from "hardhat";
import * as dotenv from "dotenv";

dotenv.config();

async function main() {
  // 获取环境变量或使用默认值
  const marketplaceAddress = process.env.MARKETPLACE_ADDRESS || "";
  const nftContractAddress = process.env.NFT_CONTRACT_ADDRESS || "";
  const tokenId = process.env.TOKEN_ID ? parseInt(process.env.TOKEN_ID) : 1;
  const reservePrice = process.env.RESERVE_PRICE || ethers.utils.parseEther("0.1").toString();
  const duration = process.env.DURATION || 86400; // 24小时默认
  
  if (!marketplaceAddress || !nftContractAddress) {
    console.error("请设置 MARKETPLACE_ADDRESS 和 NFT_CONTRACT_ADDRESS 环境变量");
    process.exit(1);
  }
  
  console.log("准备创建拍卖...");
  console.log("市场合约地址:", marketplaceAddress);
  console.log("NFT合约地址:", nftContractAddress);
  console.log("Token ID:", tokenId);
  console.log("保留价格:", ethers.utils.formatEther(reservePrice), "ETH");
  console.log("拍卖时长:", duration, "秒");
  
  const [seller] = await ethers.getSigners();
  console.log("卖家地址:", seller.address);
  
  // 获取合约实例
  const marketplaceABI = [
    "function createAuction(address nftContract, uint256 tokenId, uint256 reservePrice, uint256 duration, address paymentToken) public returns (uint256)",
  ];
  
  const nftABI = [
    "function approve(address to, uint256 tokenId) public",
    "function ownerOf(uint256 tokenId) public view returns (address)",
  ];
  
  const marketplace = new ethers.Contract(marketplaceAddress, marketplaceABI, seller);
  const nftContract = new ethers.Contract(nftContractAddress, nftABI, seller);
  
  // 检查NFT所有权
  const owner = await nftContract.ownerOf(tokenId);
  if (owner.toLowerCase() !== seller.address.toLowerCase()) {
    console.error("错误: 您不是这个NFT的所有者");
    process.exit(1);
  }
  
  // 授权市场合约操作NFT
  console.log("授权市场合约操作NFT...");
  const approveTx = await nftContract.approve(marketplaceAddress, tokenId);
  await approveTx.wait();
  console.log("授权成功");
  
  // 创建拍卖
  console.log("创建拍卖...");
  const tx = await marketplace.createAuction(
    nftContractAddress,
    tokenId,
    reservePrice,
    duration,
    ethers.constants.AddressZero // 使用ETH作为支付代币
  );
  
  console.log("交易哈希:", tx.hash);
  const receipt = await tx.wait();
  console.log("拍卖创建成功!");
  
  // 尝试从事件中提取拍卖ID
  for (const event of receipt.events || []) {
    if (event.event === "AuctionCreated") {
      console.log("拍卖ID:", event.args?.auctionId.toString());
      break;
    }
  }
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
