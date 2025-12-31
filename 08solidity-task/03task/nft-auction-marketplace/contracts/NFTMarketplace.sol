// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "./interfaces/AggregatorV3Interface.sol";

contract NFTMarketplace is Initializable, OwnableUpgradeable, UUPSUpgradeable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // 拍卖结构体
    struct Auction {
        uint256 id;
        address seller;
        address nftContract;
        uint256 tokenId;
        uint256 startTime;
        uint256 endTime;
        uint256 reservePrice;
        address highestBidder;
        uint256 highestBid;
        address paymentToken; // 0x0 表示 ETH
        bool ended;
    }

    // 状态变量
    uint256 public _auctionCounter;
    mapping(uint256 => Auction) public _auctions;
    mapping(address => uint256) public _pendingReturns;
    mapping(address => address) public _priceFeeds;

    // 事件定义
    event AuctionCreated(uint256 indexed auctionId, address indexed seller, address indexed nftContract, uint256 tokenId);
    event BidPlaced(uint256 indexed auctionId, address indexed bidder, uint256 amount);
    event AuctionEnded(uint256 indexed auctionId, address indexed winner, uint256 amount);
    event PriceFeedSet(address indexed token, address indexed aggregator);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize() public initializer {
        __Ownable_init();
        __UUPSUpgradeable_init();
        _auctionCounter = 0;
    }

    // 创建拍卖
    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 reservePrice,
        uint256 duration,
        address paymentToken
    ) public returns (uint256) {
        require(nftContract != address(0), "Invalid NFT contract");
        require(reservePrice > 0, "Reserve price must be greater than 0");
        require(duration > 0, "Duration must be greater than 0");
        
        // 验证NFT所有权并转移
        IERC721(nftContract).safeTransferFrom(msg.sender, address(this), tokenId);
        
        _auctionCounter++;
        uint256 auctionId = _auctionCounter;
        
        _auctions[auctionId] = Auction({
            id: auctionId,
            seller: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            startTime: block.timestamp,
            endTime: block.timestamp + duration,
            reservePrice: reservePrice,
            highestBidder: address(0),
            highestBid: 0,
            paymentToken: paymentToken,
            ended: false
        });
        
        emit AuctionCreated(auctionId, msg.sender, nftContract, tokenId);
        return auctionId;
    }

    // 出价
    function placeBid(uint256 auctionId) public payable nonReentrant {
        Auction storage auction = _auctions[auctionId];
        require(auction.id != 0, "Auction does not exist");
        require(!auction.ended, "Auction already ended");
        require(block.timestamp < auction.endTime, "Auction has ended");
        
        uint256 bidAmount;
        if (auction.paymentToken == address(0)) {
            // ETH 出价
            require(msg.value > 0, "ETH amount must be greater than 0");
            bidAmount = msg.value;
        } else {
            // ERC20 出价
            require(msg.value == 0, "ETH not accepted for this auction");
            require(msg.data.length > 0, "Missing bid amount");
            
            // 解析出价金额
            assembly {
                bidAmount := calldataload(4)
            }
            
            require(bidAmount > 0, "Bid amount must be greater than 0");
            IERC20(auction.paymentToken).safeTransferFrom(msg.sender, address(this), bidAmount);
        }
        
        // 验证出价是否足够高
        require(bidAmount > auction.highestBid, "Bid amount too low");
        require(bidAmount >= auction.reservePrice, "Bid below reserve price");
        
        // 退还之前最高出价者的资金
        if (auction.highestBidder != address(0)) {
            if (auction.paymentToken == address(0)) {
                _pendingReturns[auction.highestBidder] += auction.highestBid;
            } else {
                IERC20(auction.paymentToken).safeTransfer(auction.highestBidder, auction.highestBid);
            }
        }
        
        // 更新最高出价信息
        auction.highestBidder = msg.sender;
        auction.highestBid = bidAmount;
        
        emit BidPlaced(auctionId, msg.sender, bidAmount);
    }

    // 结束拍卖
    function endAuction(uint256 auctionId) public nonReentrant {
        Auction storage auction = _auctions[auctionId];
        require(auction.id != 0, "Auction does not exist");
        require(!auction.ended, "Auction already ended");
        require(block.timestamp >= auction.endTime, "Auction has not ended yet");
        
        auction.ended = true;
        
        // 处理结算
        if (auction.highestBidder != address(0)) {
            // 将NFT转移给最高出价者
            IERC721(auction.nftContract).safeTransferFrom(address(this), auction.highestBidder, auction.tokenId);
            
            // 将资金转给卖家
            if (auction.paymentToken == address(0)) {
                payable(auction.seller).transfer(auction.highestBid);
            } else {
                IERC20(auction.paymentToken).safeTransfer(auction.seller, auction.highestBid);
            }
            
            emit AuctionEnded(auctionId, auction.highestBidder, auction.highestBid);
        } else {
            // 没有出价，退还NFT给卖家
            IERC721(auction.nftContract).safeTransferFrom(address(this), auction.seller, auction.tokenId);
            emit AuctionEnded(auctionId, address(0), 0);
        }
    }

    // 提取待退还的ETH
    function withdraw() public nonReentrant returns (uint256) {
        uint256 amount = _pendingReturns[msg.sender];
        if (amount > 0) {
            _pendingReturns[msg.sender] = 0;
            payable(msg.sender).transfer(amount);
        }
        return amount;
    }

    // 获取拍卖信息
    function getAuction(uint256 auctionId) public view returns (Auction memory) {
        return _auctions[auctionId];
    }

    // 设置价格feed
    function setPriceFeed(address token, address aggregator) public onlyOwner {
        _priceFeeds[token] = aggregator;
        emit PriceFeedSet(token, aggregator);
    }

    // 获取最新价格
    function getLatestPrice(address aggregator) internal view returns (int) {
        require(aggregator != address(0), "Invalid aggregator address");
        AggregatorV3Interface priceFeed = AggregatorV3Interface(aggregator);
        
        (,
            int256 price,
            ,
            uint256 updatedAt,
        ) = priceFeed.latestRoundData();
        
        require(price > 0, "Invalid price");
        require(updatedAt > block.timestamp - 1 hours, "Stale price data");
        
        return price;
    }

    // 转换为USD
    function convertToUSD(uint256 amount, address token) public view returns (uint256) {
        if (token == address(0)) {
            // ETH 转换
            address ethPriceFeed = _priceFeeds[address(0)];
            require(ethPriceFeed != address(0), "ETH price feed not set");
            
            int256 ethPrice = getLatestPrice(ethPriceFeed);
            uint8 decimals = AggregatorV3Interface(ethPriceFeed).decimals();
            
            // 计算 USD 金额 (amount * ethPrice) / 10^decimals
            return (amount * uint256(ethPrice)) / (10 ** uint256(decimals));
        } else {
            // ERC20 转换
            address tokenPriceFeed = _priceFeeds[token];
            require(tokenPriceFeed != address(0), "Token price feed not set");
            
            int256 tokenPrice = getLatestPrice(tokenPriceFeed);
            uint8 decimals = AggregatorV3Interface(tokenPriceFeed).decimals();
            
            // 计算 USD 金额
            return (amount * uint256(tokenPrice)) / (10 ** uint256(decimals));
        }
    }

    // 比较出价
    function compareBids(uint256 bid1, address token1, uint256 bid2, address token2) internal view returns (bool) {
        if (token1 == token2) {
            return bid1 > bid2;
        }
        // 转换为USD后比较
        uint256 bid1USD = convertToUSD(bid1, token1);
        uint256 bid2USD = convertToUSD(bid2, token2);
        return bid1USD > bid2USD;
    }

    // 授权升级
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    // 接收ETH
    receive() external payable {}
}
