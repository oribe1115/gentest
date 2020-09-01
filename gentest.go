package gentest

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"io"
	"os"
	"strconv"
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
	TestFuncName   string
	ExpectedStruct string
	TestCasesDef   string
	ExecBaseFunc   string
}

type returnValue struct {
	Name     string
	Type     ast.Expr
	TypeName string
	IsError  bool
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

	returns, err := getReturnValues(pass, funcDecl)
	if err != nil {
		return nil, err
	}
	of := &outputField{}
	of.TestFuncName = genTestFuncName(funcDecl.Name.String())
	of.ExpectedStruct = genExpectedStruct(returns)
	of.TestCasesDef = genTestCasesDef(returns)
	of.ExecBaseFunc, err = genExecBaseCode(funcDecl, returns)
	if err != nil {
		return nil, err
	}
	outputTestCode(of)

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

func genExecBaseCode(baseFuncDecl *ast.FuncDecl, returns []*returnValue) (string, error) {
	var result string
	funcName := baseFuncDecl.Name.String()

	retrunNames := make([]string, 0)
	for _, re := range returns {
		retrunNames = append(retrunNames, re.Name)
	}
	if len(retrunNames) != 0 {
		result += strings.Join(retrunNames, ",") + ":="
	}
	result += fmt.Sprintf("%s()", funcName)
	return result, nil
}

// あとで整理する
func getReturnValues(pass *analysis.Pass, baseFuncDecl *ast.FuncDecl) ([]*returnValue, error) {
	values := make([]*returnValue, 0)
	results := baseFuncDecl.Type.Results
	if results == nil {
		return values, nil
	}

	nameMap := map[string]int{}
	errType := types.Universe.Lookup("error").Type()

	for _, field := range results.List {
		var typeName string
		var isError bool
		valueType := field.Type
		switch valueType := field.Type.(type) {
		case *ast.Ident:
			tv, _ := pass.TypesInfo.Types[valueType]
			isError = types.Identical(tv.Type, errType)
			typeName = tv.Type.String()
		}

		if len(field.Names) == 0 {
			name := "got" + typeName

			if count, exist := nameMap[name]; exist {
				name += strconv.Itoa(count + 1)
				nameMap[name]++
			} else {
				nameMap[name] = 1
			}

			value := &returnValue{
				Name:     name,
				Type:     valueType,
				TypeName: typeName,
				IsError:  isError,
			}

			values = append(values, value)
		} else {
			for _, ident := range field.Names {
				name := ident.Name
				if count, exist := nameMap[name]; exist {
					name += strconv.Itoa(count + 1)
					nameMap[name]++
				} else {
					nameMap[name] = 1
				}
				value := &returnValue{
					Name:     name,
					Type:     valueType,
					TypeName: typeName,
					IsError:  isError,
				}

				values = append(values, value)
			}
		}
	}

	return values, nil
}

func genExpectedStruct(returns []*returnValue) string {
	if len(returns) == 0 {
		return ""
	}
	result := "type expected struct {\n"
	for _, re := range returns {
		result += fmt.Sprintf("%s %s\n", re.Name, re.TypeName)
	}
	result += "}"

	return result
}

func genTestCasesDef(returns []*returnValue) string {
	result := "tests := []struct{"

	var expected string
	if len(returns) != 0 {
		expected = "Expected expected"
	}

	result += strings.Join([]string{expected}, ",")

	result += "}{}"
	return result
}

func outputTestCode(of *outputField) error {
	testCodeTemplate := `
func {{.TestFuncName}}(t *testing.T){
	{{.ExpectedStruct}}
	{{.TestCasesDef}}
	for _, test := range tests {
		t.Run("LABEL", func(t *testing.T) {
			{{.ExecBaseFunc}}
		})
	}
}`

	testCodeTemplate = "{{define \"base\"}}" + testCodeTemplate + "{{end}}"

	field := map[string]string{
		"TestFuncName":   of.TestFuncName,
		"ExpectedStruct": of.ExpectedStruct,
		"TestCasesDef":   of.TestCasesDef,
		"ExecBaseFunc":   of.ExecBaseFunc,
	}

	t, err := template.New("base").Parse(testCodeTemplate)
	if err != nil {
		return err
	}
	buffer := &bytes.Buffer{}
	err = t.Execute(buffer, field)
	if err != nil {
		return err
	}
	result, err := format.Source([]byte(buffer.Bytes()))
	if err != nil {
		return err
	}
	fprint(string(result))
	return nil
}
