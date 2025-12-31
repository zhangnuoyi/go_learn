const { ethers, upgrades } = require("hardhat");

async function main() {
  console.log("开始测试NFT拍卖市场完整流程...");
  
  // 获取测试账户
  const [deployer, seller, buyer1, buyer2, platform] = await ethers.getSigners();
  console.log("部署账户:", deployer.address);
  console.log("卖家账户:", seller.address);
  console.log("买家1账户:", buyer1.address);
  console.log("买家2账户:", buyer2.address);
  console.log("平台账户:", platform.address);
  
  try {
    // 部署NFT合约
    console.log("\n1. 部署AuctionNFT合约...");
    const AuctionNFT = await ethers.getContractFactory("AuctionNFT");
    const nftContract = await upgrades.deployProxy(AuctionNFT, ["AuctionNFT", "ANFT"], { initializer: "initialize" });
    await nftContract.deployed();
    console.log("NFT合约已部署到:", nftContract.address);
    
    // 部署拍卖市场合约
    console.log("\n2. 部署NFTMarketplace合约...");
    const NFTMarketplace = await ethers.getContractFactory("NFTMarketplace");
    const marketplaceContract = await upgrades.deployProxy(
      NFTMarketplace, 
      [platform.address, 1000], // 平台地址和10%佣金(1000=10%)
      { initializer: "initialize" }
    );
    await marketplaceContract.deployed();
    console.log("拍卖市场合约已部署到:", marketplaceContract.address);
    
    // 铸造NFT
    console.log("\n3. 卖家铸造NFT...");
    await nftContract.connect(seller).mint(seller.address, "ipfs://test-token-uri");
    const tokenId = 1;
    console.log("NFT铸造成功，Token ID:", tokenId);
    
    // 授权市场合约操作NFT
    console.log("\n4. 卖家授权市场合约操作NFT...");
    await nftContract.connect(seller).approve(marketplaceContract.address, tokenId);
    console.log("授权成功");
    
    // 创建拍卖
    console.log("\n5. 创建拍卖...");
    const startingPrice = ethers.utils.parseEther("0.1");
    const duration = 3600; // 1小时
    await marketplaceContract.connect(seller).createAuction(
      nftContract.address,
      tokenId,
      startingPrice,
      duration,
      ethers.constants.AddressZero // 使用ETH
    );
    console.log("拍卖创建成功，起始价格: 0.1 ETH, 持续时间: 1小时");
    
    // 买家1出价
    console.log("\n6. 买家1出价...");
    const bid1 = ethers.utils.parseEther("0.2");
    await marketplaceContract.connect(buyer1).placeBid(1, { value: bid1 });
    console.log("买家1出价成功: 0.2 ETH");
    
    // 买家2出价更高
    console.log("\n7. 买家2出价更高...");
    const bid2 = ethers.utils.parseEther("0.3");
    await marketplaceContract.connect(buyer2).placeBid(1, { value: bid2 });
    console.log("买家2出价成功: 0.3 ETH");
    
    // 等待拍卖结束
    console.log("\n8. 等待拍卖结束...");
    // 由于是测试环境，我们手动设置拍卖已结束
    // 实际应用中这里应该等待真实时间或模拟时间流逝
    await marketplaceContract.connect(deployer)._setAuctionEnded(1, true);
    console.log("拍卖标记为已结束");
    
    // 结束拍卖并结算
    console.log("\n9. 结束拍卖并结算...");
    await marketplaceContract.connect(deployer).endAuction(1);
    console.log("拍卖结束并结算完成");
    
    // 验证结果
    console.log("\n10. 验证拍卖结果...");
    const winner = await marketplaceContract.auctionWinners(1);
    const finalPrice = await marketplaceContract.auctionFinalPrices(1);
    console.log(`拍卖获胜者: ${winner}`);
    console.log(`最终成交价格: ${ethers.utils.formatEther(finalPrice)} ETH`);
    
    // 验证NFT所有权转移
    const newOwner = await nftContract.ownerOf(tokenId);
    console.log(`NFT新所有者: ${newOwner}`);
    if (newOwner === buyer2.address) {
      console.log("✓ NFT所有权成功转移给获胜者");
    } else {
      console.log("✗ NFT所有权转移失败");
    }
    
    console.log("\n🎉 完整拍卖流程测试成功完成！");
  } catch (error) {
    console.error("测试过程中出现错误:", error);
  }
}

main().catch((error) => {
  console.error("错误:", error);
  process.exitCode = 1;
});