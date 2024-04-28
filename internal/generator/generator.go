package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/waynelone/bes/assets"
	"github.com/waynelone/bes/internal/utils"

	"github.com/russross/blackfriday/v2"
)

// siteDir is the target directory into which the HTML gets generated. Its
// default is set here but can be changed by an argument passed into the
// program.
var siteDir = "./public"

func Build(dstDir string, configPath string) {
	siteDir = dstDir
	parseConfig(configPath)
	utils.EnsureDir(siteDir)

	utils.CopyFileFromFS(assets.TemplateFS, "templates/site.css", siteDir+"/site.css")
	utils.CopyFileFromFS(assets.TemplateFS, "templates/site.js", siteDir+"/site.js")
	utils.CopyFileFromFS(assets.TemplateFS, "templates/favicon.ico", siteDir+"/favicon.ico")
	utils.CopyFileFromFS(assets.TemplateFS, "templates/play.png", siteDir+"/play.png")
	utils.CopyFileFromFS(assets.TemplateFS, "templates/clipboard.png", siteDir+"/clipboard.png")
	examples := parseExamples()
	pageInfo := &IndexPageInfo{
		ConfigInfo: &configInfo,
		Examples:   examples,
	}
	renderIndex(pageInfo)
	renderExamples(examples)
	render404()
}

func verbose() bool {
	return true
}

func markdown(src string) string {
	return string(blackfriday.Run([]byte(src)))
}

func readLines(path string) []string {
	src := utils.MustReadFile(path)
	return strings.Split(src, "\n")
}

func mustGlob(glob string) []string {
	paths, err := filepath.Glob(glob)
	utils.Check(err)
	target := 0
	return utils.MoveSliceElemTo(paths, target, func(p string) bool {
		return strings.HasSuffix(p, configInfo.FileSuffix)
	})
}

func whichLexer(path string) string {
	if strings.HasSuffix(path, configInfo.FileSuffix) {
		return configInfo.Lexer
	} else if strings.HasSuffix(path, ".sh") {
		return "console"
	}
	panic("No lexer for " + path)
}

var docsPat = regexp.MustCompile(`^(\s*(\/\/|#)\s|\s*\/\/$)`)
var dashPat = regexp.MustCompile(`\-+`)

type IndexPageInfo struct {
	*ConfigInfo
	Examples []*Example
}

// Seg is a segment of an example
type Seg struct {
	Docs, DocsRendered              string
	Code, CodeRendered, CodeForJs   string
	CodeEmpty, CodeLeading, CodeRun bool
}

// Example is info extracted from an example file
type Example struct {
	ID          string
	Name        string
	GoCode      string
	Segs        [][]*Seg
	PrevExample *Example
	NextExample *Example
	*ConfigInfo
}

func parseSegs(sourcePath string) ([]*Seg, string) {
	var (
		lines  []string
		source []string
	)
	// Convert tabs to spaces for uniform rendering.
	for _, line := range readLines(sourcePath) {
		lines = append(lines, strings.Replace(line, "\t", "    ", -1))
		source = append(source, line)
	}
	segs := []*Seg{}
	lastSeen := ""
	for _, line := range lines {
		if line == "" {
			lastSeen = ""
			continue
		}
		matchDocs := docsPat.MatchString(line)
		matchCode := !matchDocs
		newDocs := (lastSeen == "") || ((lastSeen != "docs") && (segs[len(segs)-1].Docs != ""))
		newCode := (lastSeen == "") || ((lastSeen != "code") && (segs[len(segs)-1].Code != ""))
		if matchDocs {
			trimmed := docsPat.ReplaceAllString(line, "")
			if newDocs {
				newSeg := Seg{Docs: trimmed, Code: ""}
				segs = append(segs, &newSeg)
			} else {
				segs[len(segs)-1].Docs = segs[len(segs)-1].Docs + "\n" + trimmed
			}
			lastSeen = "docs"
		} else if matchCode {
			if newCode {
				newSeg := Seg{Docs: "", Code: line}
				segs = append(segs, &newSeg)
			} else {
				lastSeg := segs[len(segs)-1]
				if len(lastSeg.Code) == 0 {
					lastSeg.Code = line
				} else {
					lastSeg.Code = lastSeg.Code + "\n" + line
				}
			}
			lastSeen = "code"
		}
	}
	cr := false
	for i, seg := range segs {
		seg.CodeEmpty = (seg.Code == "")
		seg.CodeLeading = (i < (len(segs) - 1))
		seg.CodeRun = strings.HasSuffix(sourcePath, configInfo.FileSuffix) && !seg.CodeEmpty && !cr
		if seg.CodeRun {
			cr = true
		}
	}
	return segs, strings.Join(source, "\n")
}

func chromaFormat(code, filePath string) string {
	lexer := lexers.Get(filePath)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	if strings.HasSuffix(filePath, ".sh") {
		lexer = SimpleShellOutputLexer
	}

	lexer = chroma.Coalesce(lexer)

	style := styles.Get("swapoff")
	if style == nil {
		style = styles.Fallback
	}
	formatter := html.New(html.WithClasses(true))
	iterator, err := lexer.Tokenise(nil, string(code))
	utils.Check(err)
	buf := new(bytes.Buffer)
	err = formatter.Format(buf, style, iterator)
	utils.Check(err)
	return buf.String()
}

func parseAndRenderSegs(sourcePath string) ([]*Seg, string) {
	segs, filecontent := parseSegs(sourcePath)
	lexer := whichLexer(sourcePath)
	for _, seg := range segs {
		if seg.Docs != "" {
			seg.DocsRendered = markdown(seg.Docs)
		}
		if seg.Code != "" {
			seg.CodeRendered = chromaFormat(seg.Code, sourcePath)

			// adding the content to the js code for copying to the clipboard
			if strings.HasSuffix(sourcePath, configInfo.FileSuffix) {
				seg.CodeForJs = strings.Trim(seg.Code, "\n") + "\n"
			}
		}
	}
	// we are only interested in the 'go' code to pass to play.golang.org
	if lexer != configInfo.Lexer {
		filecontent = ""
	}
	return segs, filecontent
}

func parseExamples() []*Example {
	var exampleNames []string
	for _, line := range readLines("examples.txt") {
		if line != "" && !strings.HasPrefix(line, "#") {
			exampleNames = append(exampleNames, line)
		}
	}
	examples := make([]*Example, 0)
	for i, exampleName := range exampleNames {
		if verbose() {
			fmt.Printf("Processing %s [%d/%d]\n", exampleName, i+1, len(exampleNames))
		}
		example := Example{Name: exampleName}
		exampleID := strings.ToLower(exampleName)
		exampleID = strings.Replace(exampleID, " ", "-", -1)
		exampleID = strings.Replace(exampleID, "/", "-", -1)
		exampleID = strings.Replace(exampleID, "'", "", -1)
		exampleID = dashPat.ReplaceAllString(exampleID, "-")
		example.ID = exampleID
		example.Segs = make([][]*Seg, 0)
		example.ConfigInfo = &configInfo
		sourcePaths := mustGlob("examples/" + exampleID + "/*")
		for _, sourcePath := range sourcePaths {
			if !utils.IsDir(sourcePath) {
				sourceSegs, filecontents := parseAndRenderSegs(sourcePath)
				if filecontents != "" {
					example.GoCode = filecontents
				}
				example.Segs = append(example.Segs, sourceSegs)
			}
		}
		examples = append(examples, &example)
	}
	for i, example := range examples {
		if i > 0 {
			example.PrevExample = examples[i-1]
		}
		if i < (len(examples) - 1) {
			example.NextExample = examples[i+1]
		}
	}
	return examples
}

func renderIndex(pageInfo *IndexPageInfo) {
	if verbose() {
		fmt.Println("Rendering index")
	}
	indexTmpl := template.New("index")
	template.Must(indexTmpl.Parse(utils.MustReadFileFromFS(assets.TemplateFS, "templates/footer.tmpl")))
	template.Must(indexTmpl.Parse(utils.MustReadFileFromFS(assets.TemplateFS, "templates/index.tmpl")))
	indexF, err := os.Create(siteDir + "/index.html")
	utils.Check(err)
	defer indexF.Close()
	utils.Check(indexTmpl.Execute(indexF, pageInfo))
}

func renderExamples(examples []*Example) {
	if verbose() {
		fmt.Println("Rendering examples")
	}
	exampleTmpl := template.New("example")
	template.Must(exampleTmpl.Parse(utils.MustReadFileFromFS(assets.TemplateFS, "templates/footer.tmpl")))
	template.Must(exampleTmpl.Parse(utils.MustReadFileFromFS(assets.TemplateFS, "templates/example.tmpl")))
	for _, example := range examples {
		exampleF, err := os.Create(siteDir + "/" + example.ID + ".html")
		utils.Check(err)
		defer exampleF.Close()
		utils.Check(exampleTmpl.Execute(exampleF, example))
	}
}

func render404() {
	if verbose() {
		fmt.Println("Rendering 404")
	}
	tmpl := template.New("404")
	template.Must(tmpl.Parse(utils.MustReadFileFromFS(assets.TemplateFS, "templates/footer.tmpl")))
	template.Must(tmpl.Parse(utils.MustReadFileFromFS(assets.TemplateFS, "templates/404.tmpl")))
	file, err := os.Create(siteDir + "/404.html")
	utils.Check(err)
	defer file.Close()
	utils.Check(tmpl.Execute(file, &configInfo))
}

var SimpleShellOutputLexer = chroma.MustNewLexer(
	&chroma.Config{
		Name:      "Shell Output",
		Aliases:   []string{"console"},
		Filenames: []string{"*.sh"},
		MimeTypes: []string{},
	},
	func() chroma.Rules {
		return chroma.Rules{
			"root": {
				// $ or > triggers the start of prompt formatting
				{`^\$`, chroma.GenericPrompt, chroma.Push("prompt")},
				{`^>`, chroma.GenericPrompt, chroma.Push("prompt")},

				// empty lines are just text
				{`^$\n`, chroma.Text, nil},

				// otherwise its all output
				{`[^\n]+$\n?`, chroma.GenericOutput, nil},
			},
			"prompt": {
				// when we find newline, do output formatting rules
				{`\n`, chroma.Text, chroma.Push("output")},
				// otherwise its all text
				{`[^\n]+$`, chroma.Text, nil},
			},
			"output": {
				// sometimes there isn't output so we go right back to prompt
				{`^\$`, chroma.GenericPrompt, chroma.Pop(1)},
				{`^>`, chroma.GenericPrompt, chroma.Pop(1)},
				// otherwise its all output
				{`[^\n]+$\n?`, chroma.GenericOutput, nil},
			},
		}
	},
)
