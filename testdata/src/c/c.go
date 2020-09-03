package c

import (
	"context"
	"net/http"
)

func structFunc(input context.Context) context.Context {
	return input
}

func interfaceFunc(input http.Handler) http.Handler {
	return input
}
