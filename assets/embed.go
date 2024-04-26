package assets

import "embed"

//go:embed all:scaffold
var ScaffoldFS embed.FS

//go:embed all:templates
var TemplateFS embed.FS
