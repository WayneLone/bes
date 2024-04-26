package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/waynelone/bes/internal/initial"
)

type InitOptions struct {
	path string
}

var initOptions InitOptions

func init() {
	flags := cmdInit.Flags()

	flags.StringVarP(&initOptions.path, "path", "d", "./", "文件生成路径")

	cmdBes.AddCommand(cmdInit)
}

var cmdInit = &cobra.Command{
	Use:     "init",
	Short:   "初始化，生成代码示例项目",
	Version: cmdVersion.Version,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initProj(); err != nil {
			fmt.Fprintf(os.Stderr, "initial error: %v\n", err)
		}
	},
}

func initProj() error {
	initial.Init(initOptions.path)
	return nil
}
