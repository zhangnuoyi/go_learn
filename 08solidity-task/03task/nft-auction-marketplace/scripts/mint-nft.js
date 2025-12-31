const { ethers } = require("hardhat");

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log(`使用账户 ${deployer.address} 进行操作`);
  
  // 获取合约地址
  const nftAddress = process.env.NFT_ADDRESS;
  if (!nftAddress) {
    console.error("请设置 NFT_ADDRESS 环境变量");
    process.exit(1);
  }
  
  // 连接到合约
  const auctionNFT = await ethers.getContractAt("AuctionNFT", nftAddress);
  
  // 设置铸造参数
  const recipient = process.argv[2] || deployer.address;
  const tokenURI = process.argv[3] || `https://example.com/nft/${Date.now()}`;
  
  console.log(`\n铸造NFT到地址: ${recipient}`);
  console.log(`TokenURI: ${tokenURI}`);
  
  // 铸造NFT
  const tx = await auctionNFT.mint(recipient, tokenURI);
  console.log(`交易已发送，等待确认...`);
  
  const receipt = await tx.wait();
  console.log(`交易已确认！`);
  console.log(`交易哈希: ${receipt.transactionHash}`);
  
  // 获取铸造的tokenId
  const events = receipt.events.filter(event => event.event === "Transfer");
  if (events.length > 0) {
    const tokenId = events[0].args.tokenId;
    console.log(`\n✅ NFT铸造成功！`);
    console.log(`TokenID: ${tokenId.toString()}`);
    console.log(`合约地址: ${nftAddress}`);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("铸造NFT失败:", error);
    process.exit(1);
  });
