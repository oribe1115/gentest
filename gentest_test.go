package gentest_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/oribe1115/gentest"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	buffer := &bytes.Buffer{}
	gentest.SetWriter(buffer)

	testdata := analysistest.TestData()

	tests := []struct {
		Label    string
		TestDir  string
		Offset   int
		Expected string
	}{
		{
			Label:    "simple func",
			TestDir:  "a",
			Offset:   30,
			Expected: "\n\tfunc TestF() { t * testing.T }\n\t",
		},
	}

	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			offset := strconv.Itoa(test.Offset)
			gentest.Analyzer.Flags.Set("offset", offset)
			analysistest.Run(t, testdata, gentest.Analyzer, test.TestDir)
			assert.Equal(t, test.Expected, buffer.String())
		})

	}
}
