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
		Label       string
		TestPackage string
		Offset      int
		Expected    string
	}{
		{
			Label:       "simple func",
			TestPackage: "a",
			Offset:      16,
			Expected: `
func TestF(t *testing.T) {

	tests := []struct{ Label string }{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			f()

		})
	}
}`,
		},
		{
			Label:       "simple int func",
			TestPackage: "a",
			Offset:      39,
			Expected: `
func TestReturnInt(t *testing.T) {

	type expected struct {
		gotint int
	}
	tests := []struct {
		Label    string
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			gotint := returnInt()

			assert.Equal(t, test.Expected.gotint, gotint)
		})
	}
}`,
		},
		{
			Label:       "multi int func",
			TestPackage: "a",
			Offset:      72,
			Expected: `
func TestReturnInts(t *testing.T) {

	type expected struct {
		gotint  int
		gotint2 int
	}
	tests := []struct {
		Label    string
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			gotint, gotint2 := returnInts()

			assert.Equal(t, test.Expected.gotint, gotint)
			assert.Equal(t, test.Expected.gotint2, gotint2)
		})
	}
}`,
		},
		{
			Label:       "int and error func",
			TestPackage: "a",
			Offset:      123,
			Expected: `
func TestReturnIntError(t *testing.T) {

	type expected struct {
		gotint   int
		goterror error
	}
	tests := []struct {
		Label    string
		Expected expected
		IsError  bool
	}{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			gotint, goterror := returnIntError()

			if test.Expected.IsError {
				assert.Error(t, goterror)
				return
			} else {
				assert.NoError(t, goterror)
			}

			assert.Equal(t, test.Expected.gotint, gotint)
		})
	}
}`,
		},
		{
			Label:       "input ints",
			TestPackage: "a",
			Offset:      172,
			Expected: `
func TestInputInts(t *testing.T) {
	type input struct {
		a int
		b int
	}

	tests := []struct {
		Label string
		Input input
	}{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			inputInts(test.Input.a, test.Input.b)

		})
	}
}`,
		},
		{
			Label:       "intlist func",
			TestPackage: "b",
			Offset:      17,
			Expected: `
func TestIntList(t *testing.T) {
	type input struct {
		list []int
	}
	type expected struct {
		gotintList []int
	}
	tests := []struct {
		Label    string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			gotintList := intList(test.Input.list)

			assert.Equal(t, test.Expected.gotintList, gotintList)
		})
	}
}`,
		},
		{
			Label:       "map func",
			TestPackage: "b",
			Offset:      68,
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
		Label    string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			gotmp, gotmp2 := mapFunc(test.Input.input)

			assert.Equal(t, test.Expected.gotmp, gotmp)
			assert.Equal(t, test.Expected.gotmp2, gotmp2)
		})
	}
}`,
		},
		{
			Label:       "pointer func",
			TestPackage: "b",
			Offset:      176,
			Expected: `
func TestPointer(t *testing.T) {
	type input struct {
		input *string
	}
	type expected struct {
		gotpstring *string
	}
	tests := []struct {
		Label    string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			gotpstring := pointer(test.Input.input)

			assert.Equal(t, test.Expected.gotpstring, gotpstring)
		})
	}
}`,
		},
		{
			Label:       "pointer func",
			TestPackage: "b",
			Offset:      236,
			Expected: `
func TestPointerList(t *testing.T) {
	type input struct {
		input []*string
	}
	type expected struct {
		gotlist []*string
	}
	tests := []struct {
		Label    string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			gotlist := pointerList(test.Input.input)

			assert.Equal(t, test.Expected.gotlist, gotlist)
		})
	}
}`,
		},
		{
			Label:       "function func",
			TestPackage: "b",
			Offset:      296,
			Expected: `
func TestFunction(t *testing.T) {
	type input struct {
		input func(i int) string
	}
	type expected struct {
		gotfn func(i int) string
	}
	tests := []struct {
		Label    string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			gotfn := function(test.Input.input)

			assert.Equal(t, test.Expected.gotfn, gotfn)
		})
	}
}`,
		},
		{
			Label:       "chanel func",
			TestPackage: "b",
			Offset:      372,
			Expected: `
func TestChanel(t *testing.T) {
	type input struct {
		input chan int
	}
	type expected struct {
		gotch chan int
	}
	tests := []struct {
		Label    string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			gotch := chanel(test.Input.input)

			assert.Equal(t, test.Expected.gotch, gotch)
		})
	}
}`,
		},
		{
			Label:       "mystruct func",
			TestPackage: "b",
			Offset:      459,
			Expected: `
func TestMyStructFunc(t *testing.T) {
	type input struct {
		ms b.myStruct
	}
	type expected struct {
		gotbmyStruct b.myStruct
	}
	tests := []struct {
		Label    string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			gotbmyStruct := myStructFunc(test.Input.ms)

			assert.Equal(t, test.Expected.gotbmyStruct, gotbmyStruct)
		})
	}
}`,
		},
		{
			Label:       "struct func",
			TestPackage: "c",
			Offset:      38,
			Expected: `
func TestStructFunc(t *testing.T) {
	type input struct {
		input context.Context
	}
	type expected struct {
		gotcontextContext context.Context
	}
	tests := []struct {
		Label    string
		Input    input
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			gotcontextContext := structFunc(test.Input.input)

			assert.Equal(t, test.Expected.gotcontextContext, gotcontextContext)
		})
	}
}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			gentest.SetWriter(buffer)
			gentest.SetOffset(test.Offset)
			analysistest.Run(t, testdata, gentest.Analyzer, test.TestPackage)
			assert.Equal(t, test.Expected, buffer.String())
		})

	}
}
