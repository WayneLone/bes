package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynelone/bes/internal/generator"
)

type BuildOptions struct {
	path       string
	configPath string
}

var buildOptions BuildOptions

func init() {
	flags := cmdBuild.Flags()

	flags.StringVarP(&buildOptions.path, "path", "d", "./public", "指定静态文件输出目录")
	flags.StringVarP(&buildOptions.configPath, "config", "c", "./config.toml", "指定配置文件路径")

	cmdBes.AddCommand(cmdBuild)
}

var cmdBuild = &cobra.Command{
	Use:     "build",
	Short:   "生成代码示例静态文件",
	Version: cmdVersion.Version,
	Run: func(cmd *cobra.Command, args []string) {
		generator.Build(buildOptions.path, buildOptions.configPath)
	},
}
