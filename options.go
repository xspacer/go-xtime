package xtime

type options struct {
	timeLayout string
	nullLayout string
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func TimeLayout(timeLayout string) Option {
	return optionFunc(func(o *options) {
		o.timeLayout = timeLayout
	})
}

func NullLayout(nullLayout string) Option {
	return optionFunc(func(o *options) {
		o.nullLayout = nullLayout
	})
}
