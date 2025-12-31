const { expect } = require("chai");
const { ethers, upgrades, network } = require("hardhat");

describe("NFT拍卖市场集成测试", function () {
  let NFTMarketplace, marketplace, AuctionNFT, nftContract;
  let owner, seller, bidder1, bidder2, bidder3;
  let tokenId = 1;

  beforeEach(async function () {
    [owner, seller, bidder1, bidder2, bidder3] = await ethers.getSigners();
    
    // 部署NFT合约
    AuctionNFT = await ethers.getContractFactory("AuctionNFT");
    nftContract = await upgrades.deployProxy(AuctionNFT, ["AuctionNFT", "ANFT"], {
      initializer: "initialize",
      kind: "uups"
    });
    await nftContract.deployed();
    
    // 部署市场合约
    NFTMarketplace = await ethers.getContractFactory("NFTMarketplace");
    marketplace = await upgrades.deployProxy(NFTMarketplace, [10], {
      initializer: "initialize",
      kind: "uups"
    });
    await marketplace.deployed();
  });

  it("应该完整测试整个拍卖流程", async function () {
    console.log("\n1. 铸造NFT");
    // 铸造NFT给卖家
    await nftContract.mint(seller.address, "https://example.com/nft/1");
    expect(await nftContract.ownerOf(tokenId)).to.equal(seller.address);
    console.log(`   ✅ NFT ${tokenId} 已铸造给卖家 ${seller.address}`);
    
    console.log("\n2. 授权NFT给市场合约");
    // 卖家授权市场合约管理NFT
    await nftContract.connect(seller).approve(marketplace.address, tokenId);
    expect(await nftContract.getApproved(tokenId)).to.equal(marketplace.address);
    console.log(`   ✅ NFT ${tokenId} 已授权给市场合约 ${marketplace.address}`);
    
    console.log("\n3. 创建拍卖");
    // 创建拍卖
    const startingPrice = ethers.utils.parseEther("1");
    const endTime = Math.floor(Date.now() / 1000) + 10; // 10秒后结束
    await marketplace.connect(seller).createAuction(
      nftContract.address,
      tokenId,
      0, // ETH
      startingPrice,
      endTime
    );
    
    const auction = await marketplace.getAuction(1);
    expect(auction.isActive).to.be.true;
    expect(auction.seller).to.equal(seller.address);
    console.log(`   ✅ 拍卖已创建，起拍价: ${ethers.utils.formatEther(startingPrice)} ETH`);
    
    console.log("\n4. 多个用户出价");
    // 多个用户出价
    await marketplace.connect(bidder1).placeBid(1, { value: ethers.utils.parseEther("1.2") });
    console.log(`   ✅ 竞价者1出价 1.2 ETH`);
    
    await marketplace.connect(bidder2).placeBid(1, { value: ethers.utils.parseEther("1.5") });
    console.log(`   ✅ 竞价者2出价 1.5 ETH`);
    
    await marketplace.connect(bidder3).placeBid(1, { value: ethers.utils.parseEther("2.0") });
    console.log(`   ✅ 竞价者3出价 2.0 ETH (最高出价)`);
    
    const updatedAuction = await marketplace.getAuction(1);
    expect(updatedAuction.highestBidder).to.equal(bidder3.address);
    expect(updatedAuction.currentBid).to.equal(ethers.utils.parseEther("2.0"));
    
    console.log("\n5. 等待拍卖结束");
    // 等待拍卖结束
    await network.provider.send("evm_increaseTime", [15]);
    await network.provider.send("evm_mine");
    console.log("   ✅ 拍卖已结束");
    
    console.log("\n6. 结束拍卖并结算");
    // 记录卖家初始余额
    const sellerInitialBalance = await ethers.provider.getBalance(seller.address);
    
    // 结束拍卖
    await marketplace.connect(seller).endAuction(1);
    
    // 验证NFT转移给最高出价者
    expect(await nftContract.ownerOf(tokenId)).to.equal(bidder3.address);
    console.log(`   ✅ NFT已转移给最高出价者 ${bidder3.address}`);
    
    // 验证卖家收到付款（扣除10%的佣金）
    const sellerFinalBalance = await ethers.provider.getBalance(seller.address);
    const expectedPayment = ethers.utils.parseEther("1.8"); // 2.0 ETH - 10% 佣金
    const actualPayment = sellerFinalBalance.sub(sellerInitialBalance);
    
    // 允许一定的误差（gas费等）
    expect(actualPayment.gt(expectedPayment.div(100).mul(95))).to.be.true;
    console.log(`   ✅ 卖家收到约 ${ethers.utils.formatEther(actualPayment)} ETH (含10%佣金)`);
    
    console.log("\n✅ 完整拍卖流程测试成功！");
  });
});
