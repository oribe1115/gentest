package c

import (
	"context"
	"net/http"
)

// offset_basicStruct
func basicStruct(input struct{ name string }) struct{ name string } {
	return input
}

// offset_basicInterface
func basicInterface(input interface{ hoge() }) interface{ hoge() } {
	return input
}

// offset_namedStruct
func namedStruct(input context.Context) context.Context {
	return input
}

// offset_namedInterface
func namedInterface(input http.Handler) http.Handler {
	return input
}
