// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title Voting - 投票合约
 * @dev 实现了一个简单的投票合约，允许用户对投票选项进行投票。  
    ✅ 创建一个名为Voting的合约，包含以下功能：
    一个mapping来存储候选人的得票数
    一个vote函数，允许用户投票给某个候选人
    一个getVotes函数，返回某个候选人的得票数
    一个resetVotes函数，重置所有候选人的得票数
 */
contract Voting {
    struct Vote {
        uint256 yesVotes;
        uint256 noVotes;
    }

    //一个mapping来存储候选人的得票数
     mapping(VoteOption => uint256) public votes;
   // 一个vote函数，允许用户投票给某个候选人
    function vote(VoteOption option) public {
        votes[option]++;
    }
    
    //一个getVotes函数，返回某个候选人的得票数
    function getVotes(VoteOption option) public view returns (uint256) {
        return votes[option];
    }
   // 一个resetVotes函数，重置所有候选人的得票数
    function resetVotes() public {
        votes[VoteOption.Yes] = 0;
        votes[VoteOption.No] = 0;
    }
    // 投票选项
    enum VoteOption { Yes, No }
}
