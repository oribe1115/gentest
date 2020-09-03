package gentest_test

import (
	"bytes"
	"testing"

	"github.com/oribe1115/gentest"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	tests := []struct {
		Name          string
		TestPackage   string
		OffsetComment string
		Expected      string
	}{
		{
			Name:          "simple func",
			TestPackage:   "a",
			OffsetComment: "offset_f",
			Expected: `
func TestF(t *testing.T) {

	tests := []struct{ Name string }{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			f()

		})
	}
}`,
		},
		{
			Name:          "simple int func",
			TestPackage:   "a",
			OffsetComment: "offset_returnInt",
			Expected: `
func TestReturnInt(t *testing.T) {

	type expected struct {
		gotint int
	}
	tests := []struct {
		Name     string
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotint := returnInt()

			assert.Equal(t, test.Expected.gotint, gotint)
		})
	}
}`,
		},
		{
			Name:          "multi int func",
			TestPackage:   "a",
			OffsetComment: "offset_returnInts",
			Expected: `
func TestReturnInts(t *testing.T) {

	type expected struct {
		gotint  int
		gotint2 int
	}
	tests := []struct {
		Name     string
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotint, gotint2 := returnInts()

			assert.Equal(t, test.Expected.gotint, gotint)
			assert.Equal(t, test.Expected.gotint2, gotint2)
		})
	}
}`,
		},
		{
			Name:          "int and error func",
			TestPackage:   "a",
			OffsetComment: "offset_returnIntError",
			Expected: `
func TestReturnIntError(t *testing.T) {

	type expected struct {
		gotint   int
		goterror error
	}
	tests := []struct {
		Name      string
		Expected  expected
		wantError bool
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotint, goterror := returnIntError()

			if test.wantError {
				assert.Error(t, goterror)
				if test.Expected.goterror != nil {
					assert.EqualError(t, goterror, test.Expected.goterror.String())
				}
			} else {
				assert.NoError(t, goterror)
			}

			assert.Equal(t, test.Expected.gotint, gotint)
		})
	}
}`,
		},
		{
			Name:          "input ints",
			TestPackage:   "a",
			OffsetComment: "offset_inputInts",
			Expected: `
func TestInputInts(t *testing.T) {
	type input struct {
		a int
		b int
	}

	tests := []struct {
		Name  string
		Input input
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			inputInts(test.Input.a, test.Input.b)

		})
	}
}`,
		},
		{
			Name:          "intlist func",
			TestPackage:   "b",
			OffsetComment: "offset_intList",
			Expected: `
func TestIntList(t *testing.T) {
	type input struct {
		list []int
	}
	type expected struct {
		gotintList []int
	}
	tests := []struct {
		Name     string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotintList := intList(test.Input.list)

			assert.Equal(t, test.Expected.gotintList, gotintList)
		})
	}
}`,
		},
		{
			Name:          "map func",
			TestPackage:   "b",
			OffsetComment: "offset_mapFunc",
			Expected: `
func TestMapFunc(t *testing.T) {
	type input struct {
		input map[int]string
	}
	type expected struct {
		gotmp  map[int]string
		gotmp2 map[string]error
	}
	tests := []struct {
		Name     string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotmp, gotmp2 := mapFunc(test.Input.input)

			assert.Equal(t, test.Expected.gotmp, gotmp)
			assert.Equal(t, test.Expected.gotmp2, gotmp2)
		})
	}
}`,
		},
		{
			Name:          "pointer func",
			TestPackage:   "b",
			OffsetComment: "offset_pointer",
			Expected: `
func TestPointer(t *testing.T) {
	type input struct {
		input *string
	}
	type expected struct {
		gotpstring *string
	}
	tests := []struct {
		Name     string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotpstring := pointer(test.Input.input)

			assert.Equal(t, test.Expected.gotpstring, gotpstring)
		})
	}
}`,
		},
		{
			Name:          "pointer func",
			TestPackage:   "b",
			OffsetComment: "offset_pointerList",
			Expected: `
func TestPointerList(t *testing.T) {
	type input struct {
		input []*string
	}
	type expected struct {
		gotlist []*string
	}
	tests := []struct {
		Name     string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotlist := pointerList(test.Input.input)

			assert.Equal(t, test.Expected.gotlist, gotlist)
		})
	}
}`,
		},
		{
			Name:          "function func",
			TestPackage:   "b",
			OffsetComment: "offset_function",
			Expected: `
func TestFunction(t *testing.T) {
	type input struct {
		input func(i int) string
	}
	type expected struct {
		gotfn func(i int) string
	}
	tests := []struct {
		Name     string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotfn := function(test.Input.input)

			assert.Equal(t, test.Expected.gotfn, gotfn)
		})
	}
}`,
		},
		{
			Name:          "chanel func",
			TestPackage:   "b",
			OffsetComment: "offset_chanel",
			Expected: `
func TestChanel(t *testing.T) {
	type input struct {
		input chan int
	}
	type expected struct {
		gotch chan int
	}
	tests := []struct {
		Name     string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotch := chanel(test.Input.input)

			assert.Equal(t, test.Expected.gotch, gotch)
		})
	}
}`,
		},
		{
			Name:          "mystruct func",
			TestPackage:   "b",
			OffsetComment: "offset_myStructFunc",
			Expected: `
func TestMyStructFunc(t *testing.T) {
	type input struct {
		ms b.myStruct
	}
	type expected struct {
		gotmyStruct b.myStruct
	}
	tests := []struct {
		Name     string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotmyStruct := myStructFunc(test.Input.ms)

			assert.Equal(t, test.Expected.gotmyStruct, gotmyStruct)
		})
	}
}`,
		},
		{
			Name:          "basic struct func",
			TestPackage:   "c",
			OffsetComment: "offset_basicStruct",
			Expected: `
func TestBasicStruct(t *testing.T) {
	type input struct {
		input struct{ name string }
	}
	type expected struct {
		gotst struct{ name string }
	}
	tests := []struct {
		Name     string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotst := basicStruct(test.Input.input)

			assert.Equal(t, test.Expected.gotst, gotst)
		})
	}
}`,
		},
		{
			Name:          "basic interface func",
			TestPackage:   "c",
			OffsetComment: "offset_basicInterface",
			Expected: `
func TestBasicInterface(t *testing.T) {
	type input struct {
		input interface{ hoge() }
	}
	type expected struct {
		gotin interface{ hoge() }
	}
	tests := []struct {
		Name     string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotin := basicInterface(test.Input.input)

			assert.Equal(t, test.Expected.gotin, gotin)
		})
	}
}`,
		},
		{
			Name:          "named struct func",
			TestPackage:   "c",
			OffsetComment: "offset_namedStruct",
			Expected: `
func TestNamedStruct(t *testing.T) {
	type input struct {
		input context.Context
	}
	type expected struct {
		gotcontext context.Context
	}
	tests := []struct {
		Name     string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotcontext := namedStruct(test.Input.input)

			assert.Equal(t, test.Expected.gotcontext, gotcontext)
		})
	}
}`,
		},
		{
			Name:          "named interface func",
			TestPackage:   "c",
			OffsetComment: "offset_namedInterface",
			Expected: `
func TestNamedInterface(t *testing.T) {
	type input struct {
		input http.Handler
	}
	type expected struct {
		gothandler http.Handler
	}
	tests := []struct {
		Name     string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gothandler := namedInterface(test.Input.input)

			assert.Equal(t, test.Expected.gothandler, gothandler)
		})
	}
}`,
		},
		{
			Name:          "basic reciever func",
			TestPackage:   "d",
			OffsetComment: "offset_basicRecv",
			Expected: `
func TestBasicRecv(t *testing.T) {

	tests := []struct {
		Name string
		Use  d.T
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Use.basicRecv()

		})
	}
}`,
		},
		{
			Name:          "pointer reciever func",
			TestPackage:   "d",
			OffsetComment: "offset_pointerRecv",
			Expected: `
func TestPointerRecv(t *testing.T) {

	tests := []struct {
		Name string
		Use  *d.T
	}{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Use.pointerRecv()

		})
	}
}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			gentest.SetWriter(buffer)
			gentest.SetOffsetComent(test.OffsetComment)
			analysistest.Run(t, testdata, gentest.Analyzer, test.TestPackage)
			assert.Equal(t, test.Expected, buffer.String())
		})

	}
}
