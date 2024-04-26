package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const Version = "1.0.0"

func init() {
	cmdBes.AddCommand(cmdVersion)
}

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "输出 bes 当前版本",
	Run: func(cmd *cobra.Command, args []string) {
		if err := printVersion(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "print version: %v\n", err)
		}
	},
}

func printVersion(_ *cobra.Command, _ []string) error {
	fmt.Printf("Bes version:\t%s\n", Version)
	return nil
}
