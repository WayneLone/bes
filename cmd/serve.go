package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynelone/bes/internal/serve"
)

type ServeOptions struct {
	port uint
	path string
}

var serveOptions ServeOptions

func init() {
	flags := cmdServe.Flags()

	flags.StringVarP(&serveOptions.path, "path", "d", "./public", "代码示例静态文件生成目录")
	flags.UintVarP(&serveOptions.port, "port", "p", 7559, "Web 服务器端口")

	cmdBes.AddCommand(cmdServe)
}

var cmdServe = &cobra.Command{
	Use:     "serve",
	Short:   "预览代码示例",
	Version: cmdVersion.Version,
	Run: func(cmd *cobra.Command, args []string) {
		serve.Serve(serveOptions.port, serveOptions.path)
	},
}
