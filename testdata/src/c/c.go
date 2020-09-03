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
