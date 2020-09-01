package gentest

import (
	"fmt"
	"go/ast"
	"io"
	"os"
	"strings"
	"text/template"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const doc = "gentest is ..."

var writer io.Writer
var fileName string // -file flag
var offset int      // -offset flag

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "gentest",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

type outputField struct {
	TestFuncName string
}

func init() {
	writer = os.Stdout
	Analyzer.Flags.IntVar(&offset, "offset", offset, "offset")
}

func fprint(a ...interface{}) {
	fmt.Fprint(writer, a...)
}

func run(pass *analysis.Pass) (interface{}, error) {
	funcDecl, err := findTargetFunc(pass)
	if err != nil {
		return nil, err
	}

	of := &outputField{}
	of.TestFuncName = genTestFuncName(funcDecl.Name.String())
	outputTestCode(of)

	_, _ = findTargetFunc(pass)

	return nil, nil
}

func findTargetFunc(pass *analysis.Pass) (*ast.FuncDecl, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, _ := decl.(*ast.FuncDecl)
			if funcDecl == nil {
				continue
			}

			blockStmt := funcDecl.Body
			lbrecePosition := pass.Fset.Position(blockStmt.Lbrace)
			rbrecePosition := pass.Fset.Position(blockStmt.Rbrace)
			if lbrecePosition.Offset <= offset && offset <= rbrecePosition.Offset {
				return funcDecl, nil
			}
		}
	}
	return nil, fmt.Errorf("not found function with offset %d", offset)
}

func genTestFuncName(funcName string) string {
	startWithUpper := strings.ToUpper(string(funcName[0]))
	if len(funcName) > 1 {
		startWithUpper += funcName[1:]
	}
	return "Test" + startWithUpper
}

func outputTestCode(of *outputField) {
	testCodeTemplate := `
	{{define "base"}}
	func {{.TestFuncName}}(){t *testing.T}
	{{end}}
	`

	field := map[string]string{
		"TestFuncName": of.TestFuncName,
	}

	t, _ := template.New("base").Parse(testCodeTemplate)
	t.Execute(writer, field)
}
