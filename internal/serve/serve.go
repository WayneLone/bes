package serve

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/waynelone/bes/internal/utils"
)

func Serve(port uint, path string) {
	if !utils.ExistsDir(path) {
		fmt.Fprintln(os.Stderr, "代码示例静态文件目录不存在")
		return
	}
	fmt.Printf("Serving Code Example at http://127.0.0.1:%d\n", port)
	http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), http.FileServer(http.Dir(path)))
}
