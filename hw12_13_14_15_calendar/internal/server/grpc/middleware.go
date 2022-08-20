package grpc

import (
	"context"
	"fmt"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"time"
)

type interceptor struct {
	logger *logger.Logger
}

func (s *interceptor) handler(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	duration := time.Since(start)
	s.logger.Info(
		fmt.Sprintf(
			"Got request: start time %s, duration %s, method: %s, params: %s",
			start,
			duration,
			info.FullMethod,
			req,
		),
	)
	return h, err
}

func withServerUnaryInterceptor(logger *logger.Logger) grpc.ServerOption {
	i := interceptor{
		logger: logger,
	}
	return grpc.UnaryInterceptor(i.handler)
}
