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

type baseFuncData struct {
	FuncDecl  *ast.FuncDecl
	Signature *types.Signature
	Results   []*varField
}

type varField struct {
	Name     string
	Type     types.Type
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
	baseFunc := &baseFuncData{
		FuncDecl: funcDecl,
	}

	err = baseFunc.setSignature(pass)
	if err != nil {
		return nil, err
	}

	err = baseFunc.setReturns(pass)
	if err != nil {
		return nil, err
	}
	of := &outputField{}
	of.TestFuncName = genTestFuncName(funcDecl.Name.String())
	of.ExpectedStruct = genExpectedStruct(baseFunc)
	of.TestCasesDef = genTestCasesDef(baseFunc)
	of.ExecBaseFunc, err = genExecBaseCode(baseFunc)
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

func genExecBaseCode(bf *baseFuncData) (string, error) {
	var result string
	funcName := bf.FuncDecl.Name.String()

	resultNames := make([]string, 0)
	for _, re := range bf.Results {
		resultNames = append(resultNames, re.Name)
	}
	if len(resultNames) != 0 {
		result += strings.Join(resultNames, ",") + ":="
	}
	result += fmt.Sprintf("%s()", funcName)
	return result, nil
}

func (bf *baseFuncData) setSignature(pass *analysis.Pass) error {
	obj, _ := pass.TypesInfo.ObjectOf(bf.FuncDecl.Name).(*types.Func)
	if obj == nil {
		return fmt.Errorf("faild to find object for  FucDecl")
	}
	sig, _ := obj.Type().(*types.Signature)
	if sig == nil {
		return fmt.Errorf("faild to assign types.Func to types.Signature")
	}

	bf.Signature = sig

	return nil
}

func tupleToVarFields(tuple *types.Tuple, prefix string) ([]*varField, error) {
	varFields := make([]*varField, 0)
	nameMap := map[string]int{}
	errType := types.Universe.Lookup("error").Type()

	for i := 0; i < tuple.Len(); i++ {
		v := tuple.At(i)
		name := v.Name()
		if name == "" {
			name = v.Type().String()
		}

		name = prefix + name

		if count, exist := nameMap[name]; exist {
			name += strconv.Itoa(count + 1)
			nameMap[name]++
		} else {
			nameMap[name] = 1
		}

		isError := types.AssignableTo(v.Type(), errType)

		value := &varField{
			Name:     name,
			Type:     v.Type(),
			TypeName: v.Type().String(),
			IsError:  isError,
		}

		varFields = append(varFields, value)
	}

	return varFields, nil
}

func (bf *baseFuncData) setReturns(pass *analysis.Pass) error {
	results, err := tupleToVarFields(bf.Signature.Results(), "got")
	if err != nil {
		return err
	}

	bf.Results = results

	return nil
}

func genExpectedStruct(bf *baseFuncData) string {
	if len(bf.Results) == 0 {
		return ""
	}
	result := "type expected struct {\n"
	for _, re := range bf.Results {
		result += fmt.Sprintf("%s %s\n", re.Name, re.TypeName)
	}
	result += "}"

	return result
}

func genTestCasesDef(bf *baseFuncData) string {
	result := "tests := []struct{"

	var expected string
	if len(bf.Results) != 0 {
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
