// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Array {



    /**
     * @dev 将2个有序数组合并有序输出
     * @param arr1 第一个整数数组
     * @param arr2 第二个整数数组
     * @return 合并后的整数数组
     */
    function merge(uint256[] memory arr1, uint256[] memory arr2) public pure returns (uint256[] memory) {
        uint256[] memory merged = new uint256[](arr1.length + arr2.length);
        uint256 i = 0;
        uint256 j = 0;
        uint256 k = 0;
        // 合并两个有序数组，直到其中一个数组遍历完毕
        while (i < arr1.length && j < arr2.length) {
            if (arr1[i] < arr2[j]) {
                merged[k++] = arr1[i++];
            } else {
                merged[k++] = arr2[j++];
            }
        }
        // 合并剩余元素
        while (i < arr1.length) {
            merged[k++] = arr1[i++];
        }
        // 合并剩余元素
        while (j < arr2.length) {
            merged[k++] = arr2[j++];
        }
        return merged;
    }

    //给定一个有序数组 采用二分查找
    /**
     * @dev 采用二分查找在有序数组中查找目标值
     * @param arr 有序整数数组
     * @param target 要查找的目标值
     * @return 目标值在数组中的索引，如果不存在则返回-1
     */
    function binarySearch(uint256[] memory arr, uint256 target) public pure returns (int256) {
        int256 left = 0;
        int256 right = int256(arr.length) - 1;
        while (left <= right) {
            int256 mid = left + (right - left) / 2;
            if (arr[uint256(mid)] == target) {
                return mid;
            } else if (arr[uint256(mid)] < target) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
        return -1;
    }
}