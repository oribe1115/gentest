package gentest

import (
	"io"
)

func SetWriter(testWriter io.Writer) {
	writer = testWriter
}

func SetOffsetComent(oc string) {
	offsetComment = oc
}

func SetPrallelMode(pm bool) {
	parallelMode = pm
}
