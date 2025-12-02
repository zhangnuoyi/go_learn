package main

import (
	"fmt"
	"sort"
)

func main() {
	//回文数测试
	fmt.Println("回文数测试  start")
	nums := []int{4, 1, 2, 1, 2}
	println(singleNumber1(nums))
	println(singleNumber2(nums))
	fmt.Println("回文数测试  end")

	fmt.Println("括号匹配测试  start")
	s := "()[]{}"
	fmt.Println(isValid(s))
	fmt.Println("括号匹配测试  end")

	// 大整数加1测试
	fmt.Println("大整数加1测试  start")
	digits := []int{1, 2, 3}
	fmt.Println(plusOne(digits))
	fmt.Println("大整数加1测试  end")
	// 删除有序数组中的重复项测试
	fmt.Println("删除有序数组中的重复项测试  start")
	nums = []int{1, 1, 2}
	fmt.Println(removeDuplicates(nums))
	fmt.Println("删除有序数组中的重复项测试  end")
	// 合并区间测试
	fmt.Println("合并区间测试  start")
	intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	fmt.Println(merge(intervals))
	fmt.Println("合并区间测试  end")
}

// 给你一个 非空 整数数组 nums ，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
// 排序后遍历，每次跳2个，判断是否相等，不相等则返回当前值
func singleNumber1(nums []int) int {
	// 排序
	sort.Ints(nums)
	// 遍历
	for i := 0; i < len(nums); i += 2 {
		if i == len(nums)-1 {
			return nums[i]
		}
		if nums[i] != nums[i+1] {
			return nums[i]
		}
	}
	return 0

}

// 使用map 实现
func singleNumber2(nums []int) int {
	// 遍历
	mapNum := make(map[int]int)
	for _, v := range nums {
		mapNum[v]++
	}
	// 遍历map
	for k, v := range mapNum {
		if v == 1 {
			return k
		}
	}
	return 0
}

// 题目：给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效
// 实现思路：
// 1. 定义一个map，key为左括号，value为右括号
// 2. 定义一个栈，遍历字符串，遇到左括号入栈，遇到右括号出栈，判断是否匹配
// 3. 最后判断栈是否为空
// 有效字符串需满足：
// 左括号必须用相同类型的右括号闭合。
// 左括号必须以正确的顺序闭合。
// 注意空字符串可被认为是有效字符串。
func isValid(s string) bool {
	// 定义map
	mapStr := map[byte]byte{
		'(': ')',
		'[': ']',
		'{': '}',
	}
	// 定义栈
	stack := make([]byte, 0)
	// 遍历字符串
	for i := 0; i < len(s); i++ {
		// 如果是左括号，入栈
		if _, ok := mapStr[s[i]]; ok {
			stack = append(stack, s[i])
		} else {
			// 如果是右括号，判断栈是否为空
			if len(stack) == 0 {
				return false
			}
			// 出栈
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			// 判断是否匹配
			if mapStr[top] != s[i] {
				return false
			}
		}
	}
	// 判断栈是否为空
	return len(stack) == 0
}

// 给定一个表示 大整数 的整数数组 digits，其中 digits[i] 是整数的第 i 位数字。这些数字按从左到右，从最高位到最低位排列。这个大整数不包含任何前导 0。
// 实现思路：
// 1. 从最低位开始，每次加1，判断是否超过9，超过9则进位，否则直接返回
// 2. 如果最高位也进位了，则在数组最前面插入1
func plusOne(digits []int) []int {
	// 从最低位开始
	for i := len(digits) - 1; i >= 0; i-- {
		// 加1
		digits[i]++
		// 判断是否超过9
		if digits[i] <= 9 {
			return digits
		}
		// 进位
		digits[i] = 0
	}
	// 如果最高位也进位了，则在数组最前面插入1
	digits = append([]int{1}, digits...)
	return digits
}

// 删除有序数组中的重复项
// 给你一个 非严格递增排列 的数组 nums ，请你 原地 删除重复出现的元素，使每个元素 只出现一次 ，返回删除后数组的新长度。元素的 相对顺序 应该保持 一致 。然后返回 nums 中唯一元素的个数。
// 实现思路：
// 1. 定义一个指针，指向当前唯一元素的位置
// 2. 遍历数组，遇到唯一元素，指针加1，将唯一元素赋值给指针指向的位置
// 3. 最后返回指针的值加1
func removeDuplicates(nums []int) int {
	// 定义指针
	pointer := 0
	// 遍历数组
	for i := 1; i < len(nums); i++ {
		// 如果遇到唯一元素
		if nums[i] != nums[pointer] {
			// 指针加1
			pointer++
			// 将唯一元素赋值给指针指向的位置
			nums[pointer] = nums[i]
		}
	}
	// 返回指针的值加1
	return pointer + 1
}

// 合并区间
// 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。请你合并所有重叠的区间，并返回 一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间
func merge(intervals [][]int) [][]int {
	// 定义结果数组
	res := make([][]int, 0)
	// 排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	// 遍历数组
	for i := 0; i < len(intervals); i++ {
		// 如果结果数组为空，或者当前区间的左端点大于结果数组中最后一个区间的右端点
		if len(res) == 0 || intervals[i][0] > res[len(res)-1][1] {
			// 直接加入结果数组
			res = append(res, intervals[i])
		} else {
			// 否则，更新结果数组中最后一个区间的右端点
			res[len(res)-1][1] = max(res[len(res)-1][1], intervals[i][1])
		}
	}
	// 返回结果数组
	return res
}
