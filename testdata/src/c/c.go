package c

import (
	"context"
	"net/http"
)

// offset_structFunc
func structFunc(input context.Context) context.Context {
	return input
}

// offset_interfaceFunc
func interfaceFunc(input http.Handler) http.Handler {
	return input
}

// offset_basicStruct
func basicStruct(input struct{ name string }) struct{ name string } {
	return input
}

// offset_basicInterface
func basicInterface(input interface{ hoge() }) interface{ hoge() } {
	return input
}
