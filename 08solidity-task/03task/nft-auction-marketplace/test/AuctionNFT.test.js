const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("AuctionNFT", function () {
  let AuctionNFT, auctionNFT, owner, addr1;

  beforeEach(async function () {
    [owner, addr1] = await ethers.getSigners();
    AuctionNFT = await ethers.getContractFactory("AuctionNFT");
    auctionNFT = await upgrades.deployProxy(AuctionNFT, ["AuctionNFT", "ANFT"], {
      initializer: "initialize",
      kind: "uups"
    });
    await auctionNFT.deployed();
  });

  it("应该正确初始化合约", async function () {
    expect(await auctionNFT.name()).to.equal("AuctionNFT");
    expect(await auctionNFT.symbol()).to.equal("ANFT");
    expect(await auctionNFT.owner()).to.equal(owner.address);
  });

  it("应该允许铸造NFT", async function () {
    const tokenURI = "https://example.com/nft/1";
    await auctionNFT.mint(owner.address, tokenURI);
    
    expect(await auctionNFT.tokenURI(0)).to.equal(tokenURI);
    expect(await auctionNFT.ownerOf(0)).to.equal(owner.address);
    expect(await auctionNFT.totalSupply()).to.equal(1);
  });

  it("非所有者不应该能够铸造NFT", async function () {
    const tokenURI = "https://example.com/nft/1";
    await expect(
      auctionNFT.connect(addr1).mint(addr1.address, tokenURI)
    ).to.be.revertedWith("Ownable: caller is not the owner");
  });

  it("应该正确更新合约", async function () {
    // 模拟升级合约
    const AuctionNFTV2 = await ethers.getContractFactory("AuctionNFT");
    const upgraded = await upgrades.upgradeProxy(auctionNFT.address, AuctionNFTV2);
    
    // 验证升级后仍然可以正常工作
    expect(await upgraded.name()).to.equal("AuctionNFT");
    expect(await upgraded.symbol()).to.equal("ANFT");
  });
});
