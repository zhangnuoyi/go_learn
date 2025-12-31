const { ethers, upgrades } = require("hardhat");
const { expect } = require("chai");

describe("NFT Auction Marketplace Full Flow Test", async function () {
  let auctionNFT;
  let marketplace;
  let owner;
  let bidder1;
  let bidder2;
  let bidder3;
  let platformFeePercentage = 10;
  let tokenId = 1;
  let auctionId;

  before(async function () {
    // 获取签名者
    [owner, bidder1, bidder2, bidder3] = await ethers.getSigners();

    // 部署AuctionNFT合约
    console.log("部署AuctionNFT合约...");
    const AuctionNFT = await ethers.getContractFactory("AuctionNFT");
    auctionNFT = await upgrades.deployProxy(AuctionNFT, ["Auction NFT", "ANFT"], { initializer: "initialize" });
    await auctionNFT.deployed();
    console.log("AuctionNFT合约地址:", auctionNFT.address);

    // 部署NFTMarketplace合约
    console.log("部署NFTMarketplace合约...");
    const NFTMarketplace = await ethers.getContractFactory("NFTMarketplace");
    marketplace = await upgrades.deployProxy(NFTMarketplace, [platformFeePercentage], { initializer: "initialize" });
    await marketplace.deployed();
    console.log("NFTMarketplace合约地址:", marketplace.address);
  });

  it("1. 铸造NFT并授权给市场合约", async function () {
    console.log("\n1. 铸造NFT...");
    await auctionNFT.connect(owner).mint(owner.address, "https://ipfs.io/ipfs/QmWATWQ7fVPP2EFGu71UkfnqhYXDYH566qy47CnJDgvs8u");
    console.log(`NFT ${tokenId} 铸造成功，归属于 ${owner.address}`);

    console.log("2. 授权NFT给市场合约...");
    await auctionNFT.connect(owner).approve(marketplace.address, tokenId);
    console.log(`NFT ${tokenId} 已授权给市场合约 ${marketplace.address}`);
  });

  it("2. 创建拍卖", async function () {
    const startingPrice = ethers.utils.parseEther("0.1");
    const duration = 3600; // 1小时

    console.log("\n3. 创建拍卖...");
    const tx = await marketplace.connect(owner).createAuction(
      auctionNFT.address,
      tokenId,
      startingPrice,
      duration
    );

    const receipt = await tx.wait();
    // 获取auctionId
    auctionId = receipt.events.find(event => event.event === "AuctionCreated").args.auctionId;
    console.log(`拍卖创建成功，拍卖ID: ${auctionId}`);

    const auction = await marketplace.getAuction(auctionId);
    expect(auction.nftContract).to.equal(auctionNFT.address);
    expect(auction.tokenId).to.equal(tokenId);
    expect(auction.seller).to.equal(owner.address);
    expect(auction.startingPrice).to.equal(startingPrice);
    expect(auction.currentBid).to.equal(startingPrice);
    console.log(`拍卖详情验证成功: 起始价格 ${ethers.utils.formatEther(startingPrice)} ETH`);
  });

  it("3. 参与者出价", async function () {
    const bidAmount1 = ethers.utils.parseEther("0.2");
    const bidAmount2 = ethers.utils.parseEther("0.5");
    const bidAmount3 = ethers.utils.parseEther("1.0");

    console.log("\n4. 参与者1出价...");
    await marketplace.connect(bidder1).placeBid(auctionId, { value: bidAmount1 });
    console.log(`参与者1 (${bidder1.address}) 出价 ${ethers.utils.formatEther(bidAmount1)} ETH`);

    console.log("5. 参与者2出价...");
    await marketplace.connect(bidder2).placeBid(auctionId, { value: bidAmount2 });
    console.log(`参与者2 (${bidder2.address}) 出价 ${ethers.utils.formatEther(bidAmount2)} ETH`);

    console.log("6. 参与者3出价...");
    await marketplace.connect(bidder3).placeBid(auctionId, { value: bidAmount3 });
    console.log(`参与者3 (${bidder3.address}) 出价 ${ethers.utils.formatEther(bidAmount3)} ETH`);

    // 验证最高出价
    const auction = await marketplace.getAuction(auctionId);
    expect(auction.currentBid).to.equal(bidAmount3);
    expect(auction.currentBidder).to.equal(bidder3.address);
    console.log(`最高出价验证成功: ${ethers.utils.formatEther(bidAmount3)} ETH`);
  });

  it("4. 结束拍卖并结算", async function () {
    // 等待拍卖结束
    const auction = await marketplace.getAuction(auctionId);
    const endTime = auction.endTime.toNumber();
    const currentTime = Math.floor(Date.now() / 1000);
    const timeLeft = endTime - currentTime;

    if (timeLeft > 0) {
      console.log(`\n7. 等待拍卖结束... 剩余 ${timeLeft} 秒`);
      await new Promise(resolve => setTimeout(resolve, timeLeft * 1000));
    }

    // 结束拍卖
    console.log("\n8. 结束拍卖...");
    await marketplace.connect(owner).endAuction(auctionId);
    console.log(`拍卖 ${auctionId} 已结束`);

    // 验证拍卖结果
    const finalAuction = await marketplace.getAuction(auctionId);
    expect(finalAuction.status).to.equal(2); // 2 = ENDED

    // 计算平台费用和卖家获得的金额
    const finalBid = finalAuction.currentBid;
    const platformFee = finalBid.mul(platformFeePercentage).div(100);
    const sellerAmount = finalBid.sub(platformFee);

    // 验证NFT所有权转移
    const newOwner = await auctionNFT.ownerOf(tokenId);
    expect(newOwner).to.equal(bidder3.address);
    console.log(`\n拍卖结果验证:`);
    console.log(`- NFT ${tokenId} 所有权已转移至中标者 ${newOwner}`);
    console.log(`- 最终成交价: ${ethers.utils.formatEther(finalBid)} ETH`);
    console.log(`- 平台费用 (${platformFeePercentage}%): ${ethers.utils.formatEther(platformFee)} ETH`);
    console.log(`- 卖家获得金额: ${ethers.utils.formatEther(sellerAmount)} ETH`);
    console.log(`\n✅ NFT拍卖市场完整流程测试成功!`);
  });
});

// 执行测试
async function main() {
  console.log("开始NFT拍卖市场完整流程测试...");
  
  try {
    // 获取签名者
    const [owner, bidder1, bidder2, bidder3] = await ethers.getSigners();
    
    // 部署AuctionNFT合约
    console.log("部署AuctionNFT合约...");
    const AuctionNFT = await ethers.getContractFactory("AuctionNFT");
    const auctionNFT = await upgrades.deployProxy(AuctionNFT, ["Auction NFT", "ANFT"], { initializer: "initialize" });
    await auctionNFT.deployed();
    console.log("AuctionNFT合约地址:", auctionNFT.address);
    
    // 部署NFTMarketplace合约
    console.log("部署NFTMarketplace合约...");
    const NFTMarketplace = await ethers.getContractFactory("NFTMarketplace");
    const marketplace = await upgrades.deployProxy(NFTMarketplace, [10], { initializer: "initialize" });
    await marketplace.deployed();
    console.log("NFTMarketplace合约地址:", marketplace.address);
    
    // 铸造NFT
    console.log("铸造NFT...");
    await auctionNFT.connect(owner).mint(owner.address, "https://ipfs.io/ipfs/QmWATWQ7fVPP2EFGu71UkfnqhYXDYH566qy47CnJDgvs8u");
    
    // 授权NFT给市场
    console.log("授权NFT给市场合约...");
    await auctionNFT.connect(owner).approve(marketplace.address, 1);
    
    // 创建拍卖
    console.log("创建拍卖...");
    const tx = await marketplace.connect(owner).createAuction(
      auctionNFT.address,
      1,
      ethers.utils.parseEther("0.1"),
      60 // 60秒拍卖时间
    );
    
    const receipt = await tx.wait();
    const auctionId = receipt.events.find(event => event.event === "AuctionCreated").args.auctionId;
    console.log(`拍卖创建成功，拍卖ID: ${auctionId}`);
    
    // 参与者出价
    console.log("参与者1出价...");
    await marketplace.connect(bidder1).placeBid(auctionId, { value: ethers.utils.parseEther("0.2") });
    
    console.log("参与者2出价...");
    await marketplace.connect(bidder2).placeBid(auctionId, { value: ethers.utils.parseEther("0.5") });
    
    // 等待拍卖结束
    console.log("等待拍卖结束...");
    await new Promise(resolve => setTimeout(resolve, 60000)); // 等待60秒
    
    // 结束拍卖
    console.log("结束拍卖...");
    await marketplace.connect(owner).endAuction(auctionId);
    
    console.log("\n✅ NFT拍卖市场完整流程测试成功!");
  } catch (error) {
    console.error("❌ 测试过程中出现错误:", error);
    process.exit(1);
  }
}

// 如果直接运行此脚本，则执行main函数
if (require.main === module) {
  main()
    .then(() => process.exit(0))
    .catch(error => {
      console.error(error);
      process.exit(1);
    });
}

module.exports = { main };