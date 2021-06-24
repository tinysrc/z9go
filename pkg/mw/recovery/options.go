package recovery

import "context"

var defaultOptions = &options{
	recoveryHandlerFunc: nil,
}

type options struct {
	recoveryHandlerFunc RecoveryHandlerContextFunc
}

func evalOptions(opts []Option) *options {
	out := &options{}
	*out = *defaultOptions
	for _, v := range opts {
		v(out)
	}
	return out
}

type Option func(*options)

func WithRecoveryHandler(f RecoveryHandlerFunc) Option {
	return func(o *options) {
		o.recoveryHandlerFunc = RecoveryHandlerContextFunc(func(ctx context.Context, p interface{}) (err error) {
			return f(p)
		})
	}
}

func WithRecoveryContextHandler(f RecoveryHandlerContextFunc) Option {
	return func(o *options) {
		o.recoveryHandlerFunc = f
	}
}
