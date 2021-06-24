package tags

var defaultOptions = &options{
	requestFieldsFunc:     nil,
	requestFieldsFromInit: false,
}

type options struct {
	requestFieldsFunc     RequestFieldExtractorFunc
	requestFieldsFromInit bool
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

func WithFieldExtractor(f RequestFieldExtractorFunc) Option {
	return func(o *options) {
		o.requestFieldsFunc = f
	}
}

func WithFieldExtractorFromInit(f RequestFieldExtractorFunc) Option {
	return func(o *options) {
		o.requestFieldsFunc = f
		o.requestFieldsFromInit = true
	}
}
