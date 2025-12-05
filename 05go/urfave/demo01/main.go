package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := (&cli.Command{
		Name:  "三国志",
		Usage: "三国志游戏",
		Action: func(ctx context.Context, c *cli.Command) error {
			fmt.Println("三国志游戏开始")
			return nil
		},
	})

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
	}
}
