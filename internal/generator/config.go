package generator

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
	"github.com/waynelone/bes/internal/utils"
)

type ConfigInfo struct {
	Title       string
	Description string
	Authors     []Author
	SourceUrl   string
	LicenseUrl  string
	Lexer       string
	FileSuffix  string
}

type Author struct {
	Name string
	Url  string
}

var configInfo ConfigInfo

func parseConfig(configPath string) {
	data := utils.ReadFile(configPath)
	err := toml.Unmarshal(data, &configInfo)
	utils.Check(err)
	fmt.Println("Parse config.tmol successful!")
}
