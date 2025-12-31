const { expect } = require("chai");
const { ethers, upgrades, network } = require("hardhat");

describe("NFTMarketplace", function () {
  let NFTMarketplace, marketplace, AuctionNFT, nftContract, owner, seller, bidder, addr3;
  let tokenId = 1;

  beforeEach(async function () {
    [owner, seller, bidder, addr3] = await ethers.getSigners();
    
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
    
    // 铸造NFT给卖家
    await nftContract.mint(seller.address, "https://example.com/nft/1");
    
    // 卖家授权市场合约管理NFT
    await nftContract.connect(seller).approve(marketplace.address, tokenId);
  });

  it("应该正确初始化合约", async function () {
    expect(await marketplace.owner()).to.equal(owner.address);
    expect(await marketplace.feePercent()).to.equal(10);
  });

  it("应该允许创建拍卖", async function () {
    const startingPrice = ethers.utils.parseEther("1");
    const endTime = Math.floor(Date.now() / 1000) + 3600;
    
    await expect(
      marketplace.connect(seller).createAuction(
        nftContract.address,
        tokenId,
        0, // ETH
        startingPrice,
        endTime
      )
    ).to.emit(marketplace, "AuctionCreated")
     .withArgs(1, seller.address, nftContract.address, tokenId);
    
    const auction = await marketplace.getAuction(1);
    expect(auction.seller).to.equal(seller.address);
    expect(auction.nftContract).to.equal(nftContract.address);
    expect(auction.tokenId).to.equal(tokenId);
    expect(auction.token).to.equal(0); // ETH
    expect(auction.startingPrice).to.equal(startingPrice);
    expect(auction.currentBid).to.equal(startingPrice);
    expect(auction.endTime).to.equal(endTime);
    expect(auction.isActive).to.be.true;
  });

  it("应该允许出价", async function () {
    // 创建拍卖
    const startingPrice = ethers.utils.parseEther("1");
    const endTime = Math.floor(Date.now() / 1000) + 3600;
    await marketplace.connect(seller).createAuction(
      nftContract.address,
      tokenId,
      0, // ETH
      startingPrice,
      endTime
    );
    
    // 出价
    const bidAmount = ethers.utils.parseEther("1.5");
    await expect(
      marketplace.connect(bidder).placeBid(1, { value: bidAmount })
    ).to.emit(marketplace, "BidPlaced")
     .withArgs(1, bidder.address, bidAmount);
    
    const auction = await marketplace.getAuction(1);
    expect(auction.currentBid).to.equal(bidAmount);
    expect(auction.highestBidder).to.equal(bidder.address);
  });

  it("应该拒绝低于当前最高出价的出价", async function () {
    // 创建拍卖
    const startingPrice = ethers.utils.parseEther("1");
    const endTime = Math.floor(Date.now() / 1000) + 3600;
    await marketplace.connect(seller).createAuction(
      nftContract.address,
      tokenId,
      0, // ETH
      startingPrice,
      endTime
    );
    
    // 第一个出价
    await marketplace.connect(bidder).placeBid(1, { value: ethers.utils.parseEther("1.5") });
    
    // 低出价应该被拒绝
    await expect(
      marketplace.connect(addr3).placeBid(1, { value: ethers.utils.parseEther("1.2") })
    ).to.be.revertedWith("Bid amount must be higher than current bid");
  });

  it("应该允许结束拍卖并转移NFT", async function () {
    // 创建拍卖
    const startingPrice = ethers.utils.parseEther("1");
    const endTime = Math.floor(Date.now() / 1000) + 1; // 1秒后结束
    await marketplace.connect(seller).createAuction(
      nftContract.address,
      tokenId,
      0, // ETH
      startingPrice,
      endTime
    );
    
    // 出价
    await marketplace.connect(bidder).placeBid(1, { value: ethers.utils.parseEther("1.5") });
    
    // 等待拍卖结束
    await network.provider.send("evm_increaseTime", [2]);
    await network.provider.send("evm_mine");
    
    // 结束拍卖
    await expect(
      marketplace.connect(seller).endAuction(1)
    ).to.emit(marketplace, "AuctionEnded")
     .withArgs(1, bidder.address, ethers.utils.parseEther("1.5"));
    
    // 验证NFT转移
    expect(await nftContract.ownerOf(tokenId)).to.equal(bidder.address);
  });
});
