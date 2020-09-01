package gentest

import "io"

func SetWriter(testWriter io.Writer) {
	writer = testWriter
}

func ParseFlags() {
	Analyzer.Flags.IntVar(&offset, "offset", offset, "offset")
}
