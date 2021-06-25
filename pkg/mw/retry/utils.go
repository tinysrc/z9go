package retry

import (
	"context"
	"time"

	"golang.org/x/net/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func waitRetryBackoff(attempt uint, ctx context.Context, opts *options) error {
	var wait time.Duration = 0
	if attempt > 0 {
		wait = opts.backoffFunc(ctx, attempt)
	}
	if wait > 0 {
		logTrace(ctx, "grpc retry attempt=%d wait=%v", attempt, wait)
		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctxToGrpcErr(ctx.Err())
		case <-timer.C:
		}
	}
	return nil
}

func isRetriable(err error, callOpts *options) bool {
	if isCtxErr(err) {
		return false
	}
	code := status.Code(err)
	for _, v := range callOpts.codes {
		if v == code {
			return true
		}
	}
	return false
}

func isCtxErr(err error) bool {
	code := status.Code(err)
	switch code {
	case codes.DeadlineExceeded:
	case codes.Canceled:
		return true
	}
	return false
}

func ctxToGrpcErr(err error) error {
	switch err {
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return status.Error(codes.Canceled, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

func logTrace(ctx context.Context, format string, a ...interface{}) {
	tr, ok := trace.FromContext(ctx)
	if !ok {
		return
	}
	tr.LazyPrintf(format, a...)
}
