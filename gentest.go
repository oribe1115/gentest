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
var offset int           // -offset flag
var offsetComment string // for test
var parallelMode bool    // -parallel flag

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
	UseDef         string
	InputStruct    string
	ExpectedStruct string
	TestCasesDef   string
	ExecBaseFunc   string
	Asserts        string
	Parallel       string
	Cleanup        string
}

type baseFuncData struct {
	FuncDecl         *ast.FuncDecl
	Signature        *types.Signature
	Params           []*varField
	Results          []*varField
	ResultErrorCount int
	Recv             *varField
	IsRecvChenged    bool
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
	Analyzer.Flags.BoolVar(&parallelMode, "parallel", false, "parallel")
}

func fprint(a ...interface{}) {
	fmt.Fprint(writer, a...)
}

func run(pass *analysis.Pass) (interface{}, error) {
	var err error
	if offsetComment != "" {
		offset, err = getOffsetByComment(pass, offsetComment)
		if err != nil {
			return nil, err
		}
	}

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
	baseFunc.setRecv()
	baseFunc.setRecvChenged(pass)

	of := &outputField{}
	of.genTestFuncName(baseFunc)
	of.genInputStruct(baseFunc)
	of.genExpectedStruct(baseFunc)
	of.genTestCasesDef(baseFunc)
	of.genExecBaseCode(baseFunc)
	of.genAsserts(baseFunc)
	of.genParalles(baseFunc)

	err = outputTestCode(of)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func getOffsetByComment(pass *analysis.Pass, targetComment string) (int, error) {
	targetComment += "\n"
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, _ := decl.(*ast.FuncDecl)
			if funcDecl == nil || funcDecl.Doc == nil {
				continue
			}

			if funcDecl.Doc.Text() == targetComment {
				return pass.Fset.Position(funcDecl.Name.Pos()).Offset, nil
			}
		}
	}

	return 0, fmt.Errorf("faild to get offset by comment: %s", targetComment)
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

func startWithUpper(base string) string {
	result := strings.ToUpper(string(base[0]))
	if len(base) > 1 {
		result += base[1:]
	}
	return result
}

func (of *outputField) genTestFuncName(bf *baseFuncData) {
	of.TestFuncName = "Test" + startWithUpper(bf.FuncDecl.Name.Name)
}

func (of *outputField) genExecBaseCode(bf *baseFuncData) {
	funcName := bf.FuncDecl.Name.String()

	var use string
	if bf.Recv != nil {
		use = "test.Use."
	}

	var results string
	resultNames := make([]string, 0)
	for _, re := range bf.Results {
		resultNames = append(resultNames, re.Name)
	}
	if len(resultNames) != 0 {
		results = strings.Join(resultNames, ",") + ":="
	}

	paramNames := make([]string, 0)
	for _, param := range bf.Params {
		input := fmt.Sprintf("test.Input.%s", param.Name)
		paramNames = append(paramNames, input)
	}
	params := strings.Join(paramNames, ",")

	of.ExecBaseFunc = fmt.Sprintf("%s %s%s(%s)", results, use, funcName, params)
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

func typeToVarNames(t types.Type) (typeString string, varName string) {
	switch t := t.(type) {
	case *types.Basic:
		varName = t.Name()
	case *types.Array:
		if elem, _ := t.Elem().(*types.Basic); elem != nil {
			varName = elem.Name() + "List"
		} else {
			varName = "list"
		}
	case *types.Slice:
		if elem, _ := t.Elem().(*types.Basic); elem != nil {
			varName = elem.Name() + "List"
		} else {
			varName = "list"
		}
	case *types.Map:
		varName = "mp"
	case *types.Pointer:
		if elem, _ := t.Elem().(*types.Basic); elem != nil {
			varName = "p" + elem.Name()
		} else if elem, _ := t.Elem().(*types.Named); elem != nil {
			typeString, varName = typeToVarNames(elem)
			typeString = "*" + typeString
		} else {
			varName = "p"
		}
	case *types.Signature:
		varName = "fn"
	case *types.Chan:
		varName = "ch"
	case *types.Struct:
		varName = "st"
	case *types.Interface:
		varName = "in"
	case *types.Named:
		if t.Obj().Pkg() != nil {
			typeString = t.Obj().Pkg().Name() + "."
		}
		typeString += t.Obj().Name()

		objName := t.Obj().Name()
		varName = strings.ToLower(string(objName[0]))
		if len(objName) > 1 {
			varName += objName[1:]
		}
	}

	if typeString == "" {
		typeString = t.String()
	}

	return typeString, varName
}

func tupleToVarFields(tuple *types.Tuple, prefix string) ([]*varField, int) {
	varFields := make([]*varField, 0)
	nameMap := map[string]int{}
	errType := types.Universe.Lookup("error").Type()
	errCount := 0

	for i := 0; i < tuple.Len(); i++ {
		v := tuple.At(i)
		name := v.Name()

		typeString, varName := typeToVarNames(v.Type())
		if name == "" {
			name = varName
		}

		if prefix != "" {
			name = prefix + startWithUpper(name)
		}

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

func (bf *baseFuncData) setRecv() {
	recv := bf.Signature.Recv()
	if recv == nil {
		bf.Recv = nil
		return
	}

	typeName, varName := typeToVarNames(recv.Type())

	bf.Recv = &varField{
		Name:     varName,
		Type:     recv.Type(),
		TypeName: typeName,
		IsError:  false,
	}
}

func (of *outputField) genInputStruct(bf *baseFuncData) {
	if len(bf.Params) == 0 {
		return
	}

	paramDefs := make([]string, 0)
	for _, param := range bf.Params {
		paramDef := fmt.Sprintf("%s %s", param.Name, param.TypeName)
		paramDefs = append(paramDefs, paramDef)
	}

	of.InputStruct = fmt.Sprintf(
		`type input struct {
			%s
		}`,
		strings.Join(paramDefs, "\n"))
}

func (of *outputField) genExpectedStruct(bf *baseFuncData) {
	if len(bf.Results) == 0 {
		return
	}

	expectedDefs := make([]string, 0)
	for _, re := range bf.Results {
		expectedDef := fmt.Sprintf("%s %s", re.Name, re.TypeName)
		expectedDefs = append(expectedDefs, expectedDef)
	}

	of.ExpectedStruct = fmt.Sprintf(
		`type expected struct {
			%s
		}`,
		strings.Join(expectedDefs, "\n"))
}

func (of *outputField) genTestCasesDef(bf *baseFuncData) {
	elements := make([]string, 0)
	name := "Name string"
	elements = append(elements, name)

	if bf.Recv != nil {
		use := fmt.Sprintf("Use %s", bf.Recv.TypeName)
		elements = append(elements, use)
	}

	if len(bf.Params) != 0 {
		input := "Input input"
		elements = append(elements, input)
	}

	if len(bf.Results) != 0 {
		expected := "Expected expected"
		elements = append(elements, expected)
	}

	if bf.ResultErrorCount != 0 {
		want := "WantError bool"
		elements = append(elements, want)
	}

	if bf.IsRecvChenged {
		useExpected := fmt.Sprintf("UseExpected %s", bf.Recv.TypeName)
		elements = append(elements, useExpected)
	}

	of.TestCasesDef = fmt.Sprintf(
		`tests := []struct{
			%s
		}{
			// TODO: Add test cases.
		}`,
		strings.Join(elements, "\n"))
}

func (of *outputField) genAsserts(bf *baseFuncData) {
	checkErrs := make([]string, 0)
	equals := make([]string, 0)

	for _, v := range bf.Results {
		if v.IsError {
			checkErr := fmt.Sprintf(`
			if test.WantError {
				assert.Error(t, %s)
				if test.Expected.%s != nil {
					assert.EqualError(t, %s, test.Expected.%s.String())
				}
			} else {
				assert.NoError(t, %s)
			}
			`, v.Name, v.Name, v.Name, v.Name, v.Name)
			checkErrs = append(checkErrs, checkErr)
		} else {
			eq := fmt.Sprintf("assert.Equal(t, test.Expected.%s, %s)", v.Name, v.Name)
			equals = append(equals, eq)
		}
	}

	if bf.IsRecvChenged {
		useEq := "assert.Equal(t, test.UseExpected, test.Use)"
		equals = append(equals, useEq)
	}

	of.Asserts = strings.Join(checkErrs, "") + "\n" + strings.Join(equals, "\n")
}

func (of *outputField) genParalles(bf *baseFuncData) {
	if !parallelMode {
		return
	}

	of.Parallel = "t.Parallel()"
	of.Cleanup = "t.Cleanup()\n"
}

func (bf *baseFuncData) setRecvChenged(pass *analysis.Pass) {
	if bf.Recv == nil {
		return
	}
	r, _ := bf.Recv.Type.(*types.Pointer)
	if r == nil {
		return
	}

	recvIdent := bf.FuncDecl.Recv.List[0].Names[0]
	recvType := pass.TypesInfo.TypeOf(recvIdent)

	var changed bool

	ast.Inspect(bf.FuncDecl.Body, func(n ast.Node) bool {
		if changed {
			return false
		}

		switch n := n.(type) {
		case *ast.AssignStmt:
			for _, expr := range n.Lhs {
				var x *ast.Ident
				switch expr := expr.(type) {
				case *ast.SelectorExpr:
					x, _ = expr.X.(*ast.Ident)
				case *ast.Ident:
					x = expr
				}

				if x == nil {
					continue
				}

				xType := pass.TypesInfo.TypeOf(x)
				if types.Identical(xType, recvType) && x.Name == recvIdent.Name {
					changed = true
				}
			}
		}

		return true
	})

	bf.IsRecvChenged = changed
}

func outputTestCode(of *outputField) error {
	testCodeTemplate := `
func {{.TestFuncName}}(t *testing.T){
	{{.Parallel}}
	{{.InputStruct}}
	{{.ExpectedStruct}}
	{{.TestCasesDef}}
	
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			{{.Parallel}}
			{{.Cleanup}}
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
		"Parallel":       of.Parallel,
		"Cleanup":        of.Cleanup,
	}

	t, err := template.New("base").Parse(testCodeTemplate)
	if err != nil {
		return fmt.Errorf("func outputTestCode faild with parse template: %w", err)
	}
	buffer := &bytes.Buffer{}
	err = t.Execute(buffer, field)
	if err != nil {
		return fmt.Errorf("func outputTestCode faild with execute: %w", err)
	}
	result, err := format.Source([]byte(buffer.Bytes()))
	if err != nil {
		return fmt.Errorf("func outputTestCode faild with format generared test code\n%s\n: %w", buffer.String(), err)
	}

	output := string(result)
	if output == "" {
		return fmt.Errorf("error at outputTestCode")
	}
	fprint(output)
	return nil
}
