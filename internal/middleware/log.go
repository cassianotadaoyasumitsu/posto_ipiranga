package middleware

import (
	"context"
	"fmt"

	"git.wealth-park.com/cassiano/posto_ipiranga/internal/log"
	"github.com/go-kit/kit/endpoint"
)

func EndpointLoggingMiddleware(l log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// This will be executed before the request is handled by the endpoint
			l.Info("request arrived", "request", fmt.Sprintf("%+v", request))
			defer func() {
				// This will be executed after the request is handled by the endpoint
				l.Info("response returned", "response", fmt.Sprintf("%+v", response), "error", err)
			}()
			return next(ctx, request)
		}
	}
}
