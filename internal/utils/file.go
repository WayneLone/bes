package utils

import (
	"io"
	"io/fs"
	"os"
)

func IsDir(path string) bool {
	fileStat, _ := os.Stat(path)
	return fileStat.IsDir()
}

func EnsureDir(dir string) {
	err := os.MkdirAll(dir, 0755)
	Check(err)
}

func ExistsDir(dir string) bool {
	fs, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return fs.IsDir()
}

func CopyFile(src, dst string) {
	dat, err := os.ReadFile(src)
	Check(err)
	err = os.WriteFile(dst, dat, 0644)
	Check(err)
}

func CopyFileFromFS(sfs fs.FS, srcPath string, dstPath string) {
	data := readFileFromFS(sfs, srcPath)
	err := os.WriteFile(dstPath, data, 0643)
	Check(err)
}

func ReadFile(path string) []byte {
	data, err := os.ReadFile(path)
	Check(err)
	return data
}

func MustReadFile(path string) string {
	bytes, err := os.ReadFile(path)
	Check(err)
	return string(bytes)
}

func MustReadFileFromFS(tfs fs.FS, path string) string {
	return string(readFileFromFS(tfs, path))
}

func readFileFromFS(efs fs.FS, path string) []byte {
	f, err := efs.Open(path)
	Check(err)
	data, err := io.ReadAll(f)
	Check(err)
	return data
}
