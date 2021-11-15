package clientgen

import (
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"text/template"

	"github.com/swaggest/usecase"
)

//go:embed *.tpl
var templates embed.FS

type Data struct {
	Def Def
	Url string
}

const mainFile = "main.go"
const mainTemplate = "main.go.tpl"

func SourceCode(targetFolder, url string, interactors []usecase.Interactor) error {
	mustMkDir(targetFolder)

	var def = ApiDefinition(interactors)
	var t = mustTemplate(mainTemplate)
	var mainFilePath = path.Join(targetFolder, mainFile)
	var mainFile = mustCreate(mainFilePath)

	defer mustClose(mainFile)

	mustExecute(t, mainFile, Data{
		Def: def,
		Url: url,
	})

	mustFmt(mainFilePath)

	return nil
}

func mustTemplate(name string) *template.Template {
	return template.Must(template.New(name).Funcs(map[string]interface{}{"GoType": func(schemaType string) string {
		return map[string]string{
			"integer": "int",
			"string":  "string",
			"number":  "float64",
		}[schemaType]
	}}).ParseFS(templates, name))
}

func mustMkDir(path string) {
	var err = os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
}
func mustCreate(mainFilePath string) *os.File {
	var f, err = os.OpenFile(mainFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0655)
	if err != nil {
		panic(err)
	}

	return f
}

func mustExecute(t *template.Template, targetFile io.Writer, data interface{}) {
	var err = t.Execute(targetFile, data)
	if err != nil {
		panic(err)
	}
}
func mustFmt(mainFilePath string) {
	var output, err = exec.Command("gofmt", "-w", mainFilePath).CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		panic(err)
	}
}

func mustClose(closer io.Closer) {
	var err = closer.Close()
	if err != nil {
		panic(fmt.Sprintf("error closing: %v", err))
	}
}
