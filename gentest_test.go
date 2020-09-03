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
			Offset:  30,
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
			Offset:  142,
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
			Offset:  189,
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
			Offset:  244,
			Expected: `
func TestReturnIntError(t *testing.T) {
	type expected struct {
		gotint   int
		goterror error
	}
	tests := []struct {
		Label    string
		Expected expected
	}{}
	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			gotint, goterror := returnIntError()
			assert.Equal(t, test.Expected.gotint, gotint)
			assert.Equal(t, test.Expected.goterror, goterror)
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
