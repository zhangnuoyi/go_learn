// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "@openzeppelin/contracts/governance/Governor.sol";
import "@openzeppelin/contracts/governance/extensions/GovernorSettings.sol";
import "@openzeppelin/contracts/governance/extensions/GovernorCountingSimple.sol";
import "@openzeppelin/contracts/governance/extensions/GovernorVotes.sol";
import "@openzeppelin/contracts/governance/extensions/GovernorVotesQuorumFraction.sol";
import "@openzeppelin/contracts/governance/extensions/GovernorTimelockControl.sol";
import "@openzeppelin/contracts/governance/extensions/GovernorPreventLateQuorum.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
/**
 * @title NovaGovernor
 * @dev NovaToken 治理合约，实现基于代币投票权的去中心化治理
 * @notice 支持提案、投票、执行和时间锁控制
 */
 contract NovaGovernor  is 
    GovernorSettings, // 治理设置
    GovernorCountingSimple, // 简单投票计数
    GovernorVotes, // 基于代币投票权
    GovernorVotesQuorumFraction, // 基于投票数占比的 quorum
    GovernorTimelockControl, // 时间锁控制
    GovernorPreventLateQuorum, // 防止后期 quorum
    Ownable // 合约所有者
{

     /// @dev 治理参数常量
    uint48 public constant VOTING_DELAY = 1 days; // 投票延迟 1 天
    uint48 public constant VOTING_PERIOD = 1 weeks; // 投票周期 1 周\
    uint48 public constant LATE_QUORUM_DELAY = 1 days; // 后期 quorum 延迟 1 天
     uint256 public constant PROPOSAL_THRESHOLD = 1000000 * 10**18; // 100万治理代币
    uint256 public constant QUORUM_FRACTION = 4; // 投票数占比 4%

     /// @dev 提案类型
    enum ProposalType {
        UPGRADE, // 升级
        PARAMETER_CHANGE, // 参数变更
        TREASURY_ALLOCATION, // 资产分配
        OTHER // 其他
    }
    /// @dev 提案元数据
    struct ProposalMetadata{
         ProposalType proposalType; // 提案类型
        string description; // 提案描述
        uint256 createdAt; // 提案创建时间
    }
 /// @dev 提案元数据映射
    mapping(uint256 => ProposalMetadata) public proposalMetadata;


   /// @dev 提案创建事件
    /// @param proposalId 提案 ID
    /// @param proposalType 提案类型
    /// @param description 提案描述
    /// @param createdAt 提案创建时间
   event ProposalCreatedWithMetadata(
        uint256 indexed proposalId,
        ProposalType proposalType,
        string description,
        uint256 createdAt
    );

    /**
     * @dev 构造函数
     * @param _token 治理代币地址
     * @param _timelock 时间锁控制器地址
     * @param _owner 合约所有者地址
     */
    constructor(
        IVotes _token,
        TimelockController _timelock
        address _owner
    )
        Governor("NovaToken Governor");
        GovernorSettings(VOTING_DELAY, VOTING_PERIOD, PROPOSAL_THRESHOLD);
        GovernorVotes(_token);
        GovernorVotesQuorumFraction(QUORUM_FRACTION)
        GovernorTimelockControl(_timelock);
        GovernorPreventLateQuorum(LATE_QUORUM_DELAY)
        Ownable(_owner)
    {}

        /**
     * @dev 提交提案
     * @param targets 目标地址数组
     * @param values ETH 数值数组
     * @param calldatas 调用数据数组
     * @param description 描述
     * @param proposalType 提案类型
     * @return proposalId 提案ID
     */

     function ProposalCreatedWithMetadata (
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory calldatas,
        string memory description,
        ProposalType proposalType
     )public returns( uint256 proposalId){
        // 
        proposalId = propose(targets,values,calldatas,description);
        proposalMetadata[proposalId] = ProposalMetadata({
            proposalType: proposalType,
            description: description,
            createdAt: block.timestamp
        });
        emit ProposalCreatedWithMetadata(proposalId, proposalType, description, block.timestamp);
     }

         /**
     * @dev 获取提案元数据
     * @param proposalId 提案ID
     * @return ProposalMetadata 提案元数据
     */
     function getProposalMetadata(uint256 proposalId)
        public
        view
        returns (ProposalMetadata memory)
    {
        return proposalMetadata[proposalId];
    }

      /**
     * @dev 获取提案截止时间
     * @param proposalId 提案ID
     * @return uint256 截止时间戳
     */
     function proposalDeadline(uint256 proposalId)  
        public
        view
        returns (uint256)
    {
        return super.proposalDeadline(proposalId);
    }
        /**
     * @dev 检查提案是否需要排队
     * @param proposalId 提案ID
     * @return bool 是否需要排队
     */
     function proposalNeedsQueuing(uint256 proposalId)
        public
        view
        override(Governor,GovernorTimelockControl)
        returns (bool)
    {
        return super.proposalNeedsQueuing(proposalId);
    }
       /**
     * @dev 获取提案阈值
     * @return uint256 提案阈值
     */

     function proposalThreshold() public view 
     override(Governor,GovernorSettings) 
     returns (uint256)
    {
        return super.proposalThreshold();
    }
    /**
     * @dev 获取提案状态
     * @param proposalId 提案ID
     * @return ProposalState 提案状态
     */

     function state(uint256 proposalId) public view 
     override(Governor,GovernorTimelockControl)
     returns (ProposalState)
    {
        return super.state(proposalId);
    }
 /**
     * @dev 更新投票统计
     * @param proposalId 提案ID
     */
     function _tallyUpdated(uint256 proposalId)
     internal
     override(Governor,GovernorPreventLateQuorum)
     {
        super._tallyUpdated(proposalId);
     }
     /**
     * @dev 排队提案操作
     * @param proposalId 提案ID
     * @param targets 目标地址数组
     * @param values ETH 数值数组
     * @param calldatas 调用数据数组
     * @param descriptionHash 描述哈希
     * @return uint48 预计执行时间
     */
     function _queueOperations(
        uint256 proposalId,
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory calldatas,
        string memory descriptionHash
        )
     public
     override(Governor,GovernorTimelockControl)
     returns (uint48)
    {
        return super._queueOperations(proposalId, targets, values, calldatas, descriptionHash);
    }


    /**
     * @dev 执行提案操作
     * @param proposalId 提案ID
     * @param targets 目标地址数组
     * @param values ETH 数值数组
     * @param calldatas 调用数据数组
     * @param descriptionHash 描述哈希
     */
     function _executeOperations(
        uint256 proposalId,
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory calldatas,
        string memory descriptionHash
        )
     internal
     override(Governor,GovernorTimelockControl)
     {
        super._executeOperations(proposalId, targets, values, calldatas, descriptionHash);
     }

         /**
     * @dev 取消提案内部函数
     * @param targets 目标地址数组
     * @param values ETH 数值数组
     * @param calldatas 调用数据数组
     * @param descriptionHash 描述哈希
     * @return uint256 提案ID
     */
     function _cancel(
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory calldatas,
        string memory descriptionHash
     )
     internal override(Governor,GovernorTimelockControl)
     returns (uint256)
    {
        return super._cancel(targets, values, calldatas, descriptionHash);
    }
    /**
     * @dev 获取执行者地址
     * @return address 执行者地址
     */
     function _executor()
     internal 
     override(Governor,GovernorTimelockControl)
     returns (address)
    {
        return super._executor();
    }

}


    
