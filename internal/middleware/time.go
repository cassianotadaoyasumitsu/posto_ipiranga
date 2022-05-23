package middleware

import (
	"context"
	"time"

	"git.wealth-park.com/cassiano/posto_ipiranga/internal/log"
	"github.com/go-kit/kit/endpoint"
)

func EndpointTimeMiddleware(l log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// mark the start time when the request arrives
			start := time.Now()
			defer func(begin time.Time) {
				l.Info("time", "elapsed", time.Since(begin).String())
			}(start)
			return next(ctx, request)
		}
	}
}
