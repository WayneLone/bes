package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	cmdBes.PersistentFlags()
}

var cmdBes = &cobra.Command{
	Use:     "bes",
	Short:   "bes 是用来生成代码示例文档，基于 go by example 项目的脚手架",
	Version: Version,
	Run: func(c *cobra.Command, args []string) {
		if err := c.Help(); err != nil {
			fmt.Fprintf(os.Stderr, "occurs errors: %v\n", err)
		}
	},
}

func Execute() {
	if err := cmdBes.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
