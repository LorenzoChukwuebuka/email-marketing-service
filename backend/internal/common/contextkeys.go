package common

import (
	"context"
	"net/http"
)

type ctxKey string

const requestKey ctxKey = "httpRequest"

func WithRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestKey, r)
}

func RequestFromCtx(ctx context.Context) *http.Request {
	r, _ := ctx.Value(requestKey).(*http.Request)
	return r
}
