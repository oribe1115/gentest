package gentest_test

import (
	"bytes"
	"testing"

	"github.com/oribe1115/gentest"
	"github.com/sebdah/goldie/v2"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	tests := []struct {
		Name          string
		TestPackage   string
		OffsetComment string
		Parallel      bool
	}{
		{
			Name:          "f",
			TestPackage:   "a",
			OffsetComment: "offset_f",
		},
		{
			Name:          "returnInt",
			TestPackage:   "a",
			OffsetComment: "offset_returnInt",
		},
		{
			Name:          "returnInts",
			TestPackage:   "a",
			OffsetComment: "offset_returnInts",
		},
		{
			Name:          "returnIntError",
			TestPackage:   "a",
			OffsetComment: "offset_returnIntError",
		},
		{
			Name:          "inputInts",
			TestPackage:   "a",
			OffsetComment: "offset_inputInts",
		},
		{
			Name:          "intList",
			TestPackage:   "b",
			OffsetComment: "offset_intList",
		},
		{
			Name:          "mapFunc",
			TestPackage:   "b",
			OffsetComment: "offset_mapFunc",
		},
		{
			Name:          "pointer",
			TestPackage:   "b",
			OffsetComment: "offset_pointer",
		},
		{
			Name:          "pointerList",
			TestPackage:   "b",
			OffsetComment: "offset_pointerList",
		},
		{
			Name:          "function",
			TestPackage:   "b",
			OffsetComment: "offset_function",
		},
		{
			Name:          "chanel",
			TestPackage:   "b",
			OffsetComment: "offset_chanel",
		},
		{
			Name:          "myStructFunc",
			TestPackage:   "b",
			OffsetComment: "offset_myStructFunc",
		},
		{
			Name:          "basicStruct",
			TestPackage:   "c",
			OffsetComment: "offset_basicStruct",
		},
		{
			Name:          "basicInterface",
			TestPackage:   "c",
			OffsetComment: "offset_basicInterface",
		},
		{
			Name:          "namedStruct",
			TestPackage:   "c",
			OffsetComment: "offset_namedStruct",
		},
		{
			Name:          "namedInterface",
			TestPackage:   "c",
			OffsetComment: "offset_namedInterface",
		},
		{
			Name:          "basicRecv",
			TestPackage:   "d",
			OffsetComment: "offset_basicRecv",
		},
		{
			Name:          "pointerRecv",
			TestPackage:   "d",
			OffsetComment: "offset_pointerRecv",
		},
		{
			Name:          "parallel",
			TestPackage:   "d",
			OffsetComment: "offset_paralell",
			Parallel:      true,
		},
		{
			Name:          "assign",
			TestPackage:   "e",
			OffsetComment: "offset_assign",
		},
		{
			Name:          "recvChangedDirect",
			TestPackage:   "e",
			OffsetComment: "offset_recvChangedDirect",
		},
		{
			Name:          "sameTypeDiffVar",
			TestPackage:   "e",
			OffsetComment: "offset_sameTypeDiffVar",
		},
		{
			Name:          "assignInMethod",
			TestPackage:   "e",
			OffsetComment: "offset_assignInMethod",
		},
		{
			Name:          "assignInFunc",
			TestPackage:   "e",
			OffsetComment: "offset_assignInFunc",
		},
		{
			Name:          "assignInGoFunc",
			TestPackage:   "e",
			OffsetComment: "offset_assignInGoFunc",
		}, {
			Name:          "namedResults",
			TestPackage:   "f",
			OffsetComment: "offset_namedResults",
		},
		{
			Name:          "offsetFunc",
			TestPackage:   "g",
			OffsetComment: "offset_offsetFunc_beforeDec",
		},
		{
			Name:          "offsetFunc",
			TestPackage:   "g",
			OffsetComment: "offset_offsetFunc_beforeName",
		},
		{
			Name:          "offsetFunc",
			TestPackage:   "g",
			OffsetComment: "offset_offsetFunc_afterName",
		},
		{
			Name:          "offsetFunc",
			TestPackage:   "g",
			OffsetComment: "offset_offsetFunc_inBrackets",
		},
	}

	g := goldie.New(t,
		goldie.WithFixtureDir("testdata/golden"),
	)

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			gentest.SetWriter(buffer)
			gentest.SetOffsetComent(test.OffsetComment)
			gentest.SetPrallelMode(test.Parallel)
			analysistest.Run(t, testdata, gentest.Analyzer, test.TestPackage)

			g.Assert(t, test.Name, buffer.Bytes())
		})

	}
}
