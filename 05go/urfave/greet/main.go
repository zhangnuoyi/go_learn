package main

import (
    "fmt"
    "log"
    "os"
    "context"
	"strconv"

    "github.com/urfave/cli/v3"
)

func main() {
    cmd := &cli.Command{
        Name: "greet",
        EnableShellCompletion: true,
        Commands: []*cli.Command{
            {
                Name:    "add",
                Aliases: []string{"a"},
				
                Usage:   "多个数相加",
                Action: func(ctx context.Context, cmd *cli.Command) error {
                    fmt.Println("added task: ", cmd.Args().First())
					var sum float64=0.0
					for _, arg := range cmd.Args().Slice() {
						fmt.Println("arg: ", arg)
						num, err := strconv.ParseFloat(arg, 64)
						if err != nil {
							return err
						}
						sum += num
					}
                    fmt.Printf("结果: %.2f\n", sum)
                    return nil
                },
            },
            {
                Name:    "complete",
                Aliases: []string{"c"},
                Usage:   "complete a task on the list",
                Action: func(ctx context.Context, cmd *cli.Command) error {
					for _, arg := range cmd.Args().Slice() {
						fmt.Println("completed task: ", arg)
					}

                    fmt.Println("completed task: ", cmd.Args().First())
                    return nil
                },
            },
            {
                Name:    "template",
                Aliases: []string{"t"},
                Usage:   "options for task templates",
                Commands: []*cli.Command{
                    {
                        Name:  "add",
                        Usage: "add a new template",
                        Action: func(ctx context.Context, cmd *cli.Command) error {
                            fmt.Println("new task template: ", cmd.Args().First())
                            return nil
                        },
                    },
                    {
                        Name:  "remove",
                        Usage: "remove an existing template",
                        Action: func(ctx context.Context, cmd *cli.Command) error {
                            fmt.Println("removed task template: ", cmd.Args().First())
                            return nil
                        },
                    },
                },
            },
        },
    }

    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}