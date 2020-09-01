package gentest

import "io"

func SetWriter(testWriter io.Writer) {
	writer = testWriter
}
