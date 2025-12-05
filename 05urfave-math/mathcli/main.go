package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/urfave/cli/v3"
)

/*
*
* 主函数
1.使用 go install 到go/bin目录下
2. 运行 mathcli 1+1 即可计算 1+1 的结果
* @param args 命令行参数
*/
func main() {
	cmd := &cli.Command{
		Name:  "math",
		Usage: "简单的命令行计算器",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// 检查是否提供了表达式
			args := cmd.Args()
			if args.Len() == 0 {
				return fmt.Errorf("请提供一个表达式，例如：mathcli 1+1")
			}

			expression := args.First()
			result, err := calculate(expression)
			if err != nil {
				return err
			}

			fmt.Println("计算结果:", expression, "=", result)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

/**
 * 计算表达式的结果
 * @param expression 表达式字符串，例如："1+1"
 * @return 计算结果
 * @return 错误信息（如果有）
 */
func calculate(expression string) (float64, error) {
	// Regex to match mathematical expressions with +, -, *, /, %
	//使用正则表达式进行解析
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([+\-*/%])\s*(\d+(?:\.\d+)?)$`)
	matches := re.FindStringSubmatch(expression)

	if len(matches) != 4 {
		return 0, fmt.Errorf("无效的表达式格式，请使用：number operator number")
	}

	aStr := matches[1]
	operator := matches[2]
	bStr := matches[3]

	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		return 0, fmt.Errorf("第一个操作数无效：%v", err)
	}

	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		return 0, fmt.Errorf("第二个操作数无效：%v", err)
	}

	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("除数不能为零")
		}
		return a / b, nil
	case "%":
		// For modulo, convert to integers if both numbers are whole numbers
		if a == float64(int64(a)) && b == float64(int64(b)) {
			return float64(int64(a) % int64(b)), nil
		}
		// For floating point modulo, use the standard library
		return float64(int64(a) % int64(b)), nil
	default:
		return 0, fmt.Errorf("不支持的运算符：%s", operator)
	}
}
