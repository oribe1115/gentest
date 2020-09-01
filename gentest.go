package gentest

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"strings"
	"text/template"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "gentest is ..."

var writer io.Writer
var funcName string // -func flag

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
	Analyzer.Flags.StringVar(&funcName, "func", "funcName", "fuction name for generateing test code")
}

func fprint(a ...interface{}) {
	fmt.Fprint(writer, a...)
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
	})

	of := &outputField{}
	of.TestFuncName = genTestFuncName(funcName)
	outputTestCode(of)

	_, _ = findTargetFunc(pass)

	return nil, nil
}

func findTargetFunc(pass *analysis.Pass) (*ast.FuncDecl, error) {
	for _, file := range pass.Files {
		fmt.Printf("filename: %s\n", file.Name.String())
		token.NewFileSet()

		for _, decl := range file.Decls {
			funcDecl, _ := decl.(*ast.FuncDecl)
			if funcDecl == nil {
				continue
			}

			fmt.Println(funcDecl.Name.String())

			blockStmt := funcDecl.Body
			lbrecePosition := pass.Fset.Position(blockStmt.Lbrace)
			fmt.Printf("lbrace: %d, line: %d\n", lbrecePosition.Offset, lbrecePosition.Line)
			rbrecePosition := pass.Fset.Position(blockStmt.Rbrace)
			fmt.Printf("rbrace: %d, line: %d\n", rbrecePosition.Offset, rbrecePosition.Line)

		}
	}
	return nil, nil
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
