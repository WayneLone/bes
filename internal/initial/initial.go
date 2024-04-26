package initial

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/waynelone/bes/assets"
	"github.com/waynelone/bes/internal/utils"
)

func Init(dstDir string) {
	sfs, err := fs.Sub(assets.ScaffoldFS, "scaffold")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read scaffold error: %v\n", err)
		return
	}
	err = fs.WalkDir(sfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, path)
		if d.IsDir() {
			os.MkdirAll(dstPath, 0755)
		} else {
			utils.CopyFileFromFS(sfs, path, dstPath)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "generate code example files error: %v\n", err)
	}
}
