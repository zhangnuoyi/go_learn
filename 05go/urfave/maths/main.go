package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/urfave/cli/v3"
)

// 将字符串参数转换为浮点数
func parseArgs(args []string) ([]float64, error) {
	var numbers []float64
	fmt.Println("解析参数:", args)
	for _, arg := range args {
		num, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return nil, fmt.Errorf("无效的数字参数: %s", arg)
		}
		numbers = append(numbers, num)
	}
	return numbers, nil
}

// 表达式解析器和计算器
func calculateExpression(expr string) (float64, error) {
	expr = strings.ReplaceAll(expr, " ", "") // 移除所有空格
	if expr == "" {
		return 0, fmt.Errorf("表达式不能为空")
	}

	// 使用递归下降解析器来处理表达式
	// 支持的操作符: +, -, *, / (优先级: * / 高于 + -)

	i := 0

	// 声明所有解析函数
	var (
		parseNumber     func() float64
		parseFactor     func() (float64, error)
		parseTerm       func() (float64, error)
		parseExpression func() (float64, error)
	)

	// 解析数字
	parseNumber = func() float64 {
		start := i
		for i < len(expr) && (unicode.IsDigit(rune(expr[i])) || expr[i] == '.') {
			i++
		}
		num, _ := strconv.ParseFloat(expr[start:i], 64)
		return num
	}

	// 解析因子 (数字或括号内的表达式)
	parseFactor = func() (float64, error) {
		if i >= len(expr) {
			return 0, fmt.Errorf("表达式格式错误")
		}

		if expr[i] == '(' {
			i++ // 跳过 '('
			result, err := parseExpression()
			if err != nil {
				return 0, err
			}
			if i >= len(expr) || expr[i] != ')' {
				return 0, fmt.Errorf("缺少右括号")
			}
			i++ // 跳过 ')'
			return result, nil
		}

		// 处理负号
		if expr[i] == '-' {
			i++
			factor, err := parseFactor()
			if err != nil {
				return 0, err
			}
			return -factor, nil
		}

		if expr[i] == '+' {
			i++
			return parseFactor()
		}

		return parseNumber(), nil
	}

	// 解析项 (* 和 /)
	parseTerm = func() (float64, error) {
		result, err := parseFactor()
		if err != nil {
			return 0, err
		}

		for i < len(expr) {
			op := expr[i]
			if op != '*' && op != '/' {
				break
			}
			i++

			factor, err := parseFactor()
			if err != nil {
				return 0, err
			}

			switch op {
			case '*':
				result *= factor
			case '/':
				if factor == 0 {
					return 0, fmt.Errorf("除数不能为零")
				}
				result /= factor
			}
		}

		return result, nil
	}

	// 解析表达式 (+ 和 -)
	parseExpression = func() (float64, error) {
		result, err := parseTerm()
		if err != nil {
			return 0, err
		}

		for i < len(expr) {
			op := expr[i]
			if op != '+' && op != '-' {
				break
			}
			i++

			term, err := parseTerm()
			if err != nil {
				return 0, err
			}

			switch op {
			case '+':
				result += term
			case '-':
				result -= term
			}
		}

		// 移除对剩余字符的检查，因为括号后面的操作符应该由上层函数处理
		return result, nil
	}

	return parseExpression()
}

func main() {
	result, err := calculateExpression("(2 + 3) * 4 - 5")
	if err != nil {
		log.Fatalf("计算表达式失败: %v", err)
	}
	fmt.Printf("结果: %.2f\n", result)

	cmd := &cli.Command{
		Name:  "math",
		Usage: "命令行计算器",
		Commands: []*cli.Command{
			{
				Name:      "eval",
				Aliases:   []string{"e"},
				Usage:     "计算数学表达式",
				ArgsUsage: "<表达式>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.NArg() == 0 {
						return fmt.Errorf("请提供一个数学表达式")
					}

					expr := strings.Join(cmd.Args().Slice(), " ")
					result, err := calculateExpression(expr)
					if err != nil {
						return fmt.Errorf("计算表达式失败: %w", err)
					}

					fmt.Printf("表达式: %s\n", expr)
					fmt.Printf("结果: %.2f\n", result)
					return nil
				},
			},
			{
				Name:      "add",
				Aliases:   []string{"+"},
				Usage:     "相加两个或多个数字",
				ArgsUsage: "<数字1> <数字2> [数字...]",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					numbers, err := parseArgs(cmd.Args().Slice())
					if err != nil {
						return err
					}

					if len(numbers) < 2 {
						return fmt.Errorf("至少需要两个数字进行相加")
					}

					result := 0.0
					for _, num := range numbers {
						result += num
					}

					fmt.Printf("结果: %.2f\n", result)
					return nil
				},
			},
			{
				Name:      "subtract",
				Usage:     "用第一个数字减去后面的数字",
				Aliases:   []string{"-"},
				ArgsUsage: "<数字1> <数字2> [数字...]",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					numbers, err := parseArgs(cmd.Args().Slice())
					if err != nil {
						return err
					}

					if len(numbers) < 2 {
						return fmt.Errorf("至少需要两个数字进行相减")
					}

					result := numbers[0]
					for i := 1; i < len(numbers); i++ {
						result -= numbers[i]
					}

					fmt.Printf("结果: %.2f\n", result)
					return nil
				},
			},
			{
				Name:      "multiply",
				Aliases:   []string{"*"},
				Usage:     "相乘两个或多个数字",
				ArgsUsage: "<数字1> <数字2> [数字...]",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					numbers, err := parseArgs(cmd.Args().Slice())
					if err != nil {
						return err
					}

					if len(numbers) < 2 {
						return fmt.Errorf("至少需要两个数字进行相乘")
					}

					result := 1.0
					for _, num := range numbers {
						result *= num
					}

					fmt.Printf("结果: %.2f\n", result)
					return nil
				},
			},
			{
				Name:      "divide",
				Aliases:   []string{"/"},
				Usage:     "用第一个数字除以后面的数字",
				ArgsUsage: "<数字1> <数字2> [数字...]",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					numbers, err := parseArgs(cmd.Args().Slice())
					if err != nil {
						return err
					}

					if len(numbers) < 2 {
						return fmt.Errorf("至少需要两个数字进行相除")
					}

					result := numbers[0]
					for i := 1; i < len(numbers); i++ {
						if numbers[i] == 0 {
							return fmt.Errorf("除数不能为零")
						}
						result /= numbers[i]
					}

					fmt.Printf("结果: %.2f\n", result)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
