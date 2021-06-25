package retry

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	DefaultRetryCodes = []codes.Code{
		codes.ResourceExhausted,
		codes.Unavailable,
	}
	defaultOptions = &options{
		max:       0,
		timeout:   0,
		incHeader: true,
		codes:     DefaultRetryCodes,
		backoffFunc: BackoffContextFunc(func(ctx context.Context, attempt uint) time.Duration {
			return BackoffLinear(50 * time.Microsecond)(attempt)
		}),
	}
)

type BackoffFunc func(attempt uint) time.Duration
type BackoffContextFunc func(ctx context.Context, attempt uint) time.Duration

func Disable() CallOption {
	return WithMax(0)
}

func WithMax(count uint) CallOption {
	return CallOption{
		applyFunc: func(opts *options) {
			opts.max = count
		},
	}
}

func WithBackoff(f BackoffFunc) CallOption {
	return CallOption{
		applyFunc: func(opts *options) {
			opts.backoffFunc = BackoffContextFunc(func(ctx context.Context, attempt uint) time.Duration {
				return f(attempt)
			})
		},
	}
}

func WithBackoffContext(f BackoffContextFunc) CallOption {
	return CallOption{
		applyFunc: func(opts *options) {
			opts.backoffFunc = f
		},
	}
}

func WithCodes(codes ...codes.Code) CallOption {
	return CallOption{
		applyFunc: func(opts *options) {
			opts.codes = codes
		},
	}
}

func WithTimeout(timeout time.Duration) CallOption {
	return CallOption{
		applyFunc: func(opts *options) {
			opts.timeout = timeout
		},
	}
}

type options struct {
	max         uint
	timeout     time.Duration
	incHeader   bool
	codes       []codes.Code
	backoffFunc BackoffContextFunc
}

type CallOption struct {
	grpc.EmptyCallOption
	applyFunc func(opts *options)
}

func reuseOrNewWithCallOptions(opts *options, callOpts []CallOption) *options {
	if len(callOpts) == 0 {
		return opts
	}
	out := &options{}
	*out = *opts
	for _, v := range callOpts {
		v.applyFunc(out)
	}
	return out
}

func filterCallOptions(callOpts []grpc.CallOption) (grpcOpts []grpc.CallOption, retryOpts []CallOption) {
	for _, v := range callOpts {
		if co, ok := v.(CallOption); ok {
			retryOpts = append(retryOpts, co)
		} else {
			grpcOpts = append(grpcOpts, v)
		}
	}
	return
}
