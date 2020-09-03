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
var offset int // -offset flag

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
	InputStruct    string
	ExpectedStruct string
	TestCasesDef   string
	ExecBaseFunc   string
	Asserts        string
}

type baseFuncData struct {
	FuncDecl         *ast.FuncDecl
	Signature        *types.Signature
	Params           []*varField
	Results          []*varField
	ResultErrorCount int
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

	baseFunc.setReturns()
	baseFunc.setParams()

	of := &outputField{}
	of.TestFuncName = genTestFuncName(funcDecl.Name.String())
	of.InputStruct = genInputStruct(baseFunc)
	of.ExpectedStruct = genExpectedStruct(baseFunc)
	of.TestCasesDef = genTestCasesDef(baseFunc)
	of.ExecBaseFunc, err = genExecBaseCode(baseFunc)
	of.Asserts = genAsserts(baseFunc)
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

			startOffset := pass.Fset.Position(funcDecl.Name.Pos()).Offset
			endOffset := pass.Fset.Position(funcDecl.Name.End()).Offset

			if startOffset <= offset && offset <= endOffset {
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

	paramNames := make([]string, 0)
	for _, param := range bf.Params {
		input := fmt.Sprintf("test.Input.%s", param.Name)
		paramNames = append(paramNames, input)
	}

	result += fmt.Sprintf("%s(%s)", funcName, strings.Join(paramNames, ","))
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

func tupleToVarFields(tuple *types.Tuple, prefix string) ([]*varField, int) {
	varFields := make([]*varField, 0)
	nameMap := map[string]int{}
	errType := types.Universe.Lookup("error").Type()
	errCount := 0

	for i := 0; i < tuple.Len(); i++ {
		v := tuple.At(i)
		typeString := v.Type().String()
		name := v.Name()

		if name == "" {
			switch vType := v.Type().(type) {
			case *types.Array:
				if elem, _ := vType.Elem().(*types.Basic); elem != nil {
					name = elem.Name() + "List"
				} else {
					name = "list"
				}
			case *types.Slice:
				if elem, _ := vType.Elem().(*types.Basic); elem != nil {
					name = elem.Name() + "List"
				} else {
					name = "list"
				}
			case *types.Map:
				name = "mp"
			case *types.Pointer:
				if elem, _ := vType.Elem().(*types.Basic); elem != nil {
					name = "p" + elem.Name()
				} else {
					name = "p"
				}
			case *types.Signature:
				name = "fn"
			case *types.Chan:
				name = "ch"
			default:
				name = typeString
			}
		}

		name = prefix + name

		if count, exist := nameMap[name]; exist {
			name += strconv.Itoa(count + 1)
			nameMap[name]++
		} else {
			nameMap[name] = 1
		}

		isError := types.AssignableTo(v.Type(), errType)
		if isError {
			errCount++
		}

		value := &varField{
			Name:     name,
			Type:     v.Type(),
			TypeName: typeString,
			IsError:  isError,
		}

		varFields = append(varFields, value)
	}

	return varFields, errCount
}

func (bf *baseFuncData) setReturns() {
	bf.Results, bf.ResultErrorCount = tupleToVarFields(bf.Signature.Results(), "got")
}

func (bf *baseFuncData) setParams() {
	bf.Params, _ = tupleToVarFields(bf.Signature.Params(), "")
}

func genInputStruct(bf *baseFuncData) string {
	if len(bf.Params) == 0 {
		return ""
	}
	result := "type input struct {\n"
	for _, param := range bf.Params {
		result += fmt.Sprintf("%s %s\n", param.Name, param.TypeName)
	}
	result += "}"
	return result
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
	elements := make([]string, 0)
	label := "Label string"
	elements = append(elements, label)

	if len(bf.Results) != 0 {
		expected := "Expected expected"
		elements = append(elements, expected)
	}

	if bf.ResultErrorCount != 0 {
		isErr := "IsError bool"
		elements = append(elements, isErr)
	}

	return "tests := []struct{" + strings.Join(elements, "\n") + "}{}"
}
func genAsserts(be *baseFuncData) string {
	checkErrs := make([]string, 0)
	equals := make([]string, 0)

	for _, v := range be.Results {
		if v.IsError {
			checkErr := fmt.Sprintf(`
			if test.Expected.IsError {
				assert.Error(t, %s)
				return
			} else {
				assert.NoError(t, %s)
			}
			`, v.Name, v.Name)
			checkErrs = append(checkErrs, checkErr)
		} else {
			eq := fmt.Sprintf("assert.Equal(t, test.Expected.%s, %s)", v.Name, v.Name)
			equals = append(equals, eq)
		}
	}

	return strings.Join(checkErrs, "") + "\n" + strings.Join(equals, "\n")
}

func outputTestCode(of *outputField) error {
	testCodeTemplate := `
func {{.TestFuncName}}(t *testing.T){
	{{.InputStruct}}
	{{.ExpectedStruct}}
	{{.TestCasesDef}}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			{{.ExecBaseFunc}}
			{{.Asserts}}
		})
	}
}`

	testCodeTemplate = "{{define \"base\"}}" + testCodeTemplate + "{{end}}"

	field := map[string]string{
		"TestFuncName":   of.TestFuncName,
		"InputStruct":    of.InputStruct,
		"ExpectedStruct": of.ExpectedStruct,
		"TestCasesDef":   of.TestCasesDef,
		"ExecBaseFunc":   of.ExecBaseFunc,
		"Asserts":        of.Asserts,
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
