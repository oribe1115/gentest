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
		Label    string
		TestDir  string
		Offset   int
		Expected string
	}{
		{
			Label:   "simple func",
			TestDir: "a",
			Offset:  17,
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
			Label:   "simple int func",
			TestDir: "a",
			Offset:  117,
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
			Label:   "multi int func",
			TestDir: "a",
			Offset:  156,
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
			Label:   "int and error func",
			TestDir: "a",
			Offset:  198,
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
			Label:   "input ints",
			TestDir: "a",
			Offset:  257,
			Expected: `
func TestInputInts(t *testing.T) {
	type input struct {
		a int
		b int
	}

	tests := []struct{ Label string }{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			inputInts(test.Input.a, test.Input.b)

		})
	}
}`,
		},
		{
			Label:   "intlist func",
			TestDir: "a",
			Offset:  290,
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
	}

	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			gentest.SetWriter(buffer)
			gentest.SetOffset(test.Offset)
			analysistest.Run(t, testdata, gentest.Analyzer, test.TestDir)
			assert.Equal(t, test.Expected, buffer.String())
		})

	}
}
