package backoff

import (
	"math"
	"math/rand"
	"time"
)

const maxInt64 = float64(math.MaxInt64 - 512)

type Options struct {
	min, max time.Duration
	factor   float64
}

type Option func(*Options)

type Backoff struct {
	opts    Options
	attempt int
}

func New(opts ...Option) Backoff {
	options := NewOptions(opts...)
	return Backoff{
		opts: options,
	}
}

func NewOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.min <= 0 {
		opt.min = 100 * time.Millisecond
	}

	if opt.max <= 0 {
		opt.max = 10 * time.Second
	}

	if opt.min > opt.max {
		opt.min, opt.max = opt.max, opt.min
	}

	if opt.factor <= 0 {
		opt.factor = 2
	}

	return opt
}

func Min(min time.Duration) Option {
	return func(o *Options) {
		o.min = min
	}
}

func Max(max time.Duration) Option {
	return func(o *Options) {
		o.max = max
	}
}

func Factor(factor float64) Option {
	return func(o *Options) {
		o.factor = factor
	}
}

func (b *Backoff) Duration() time.Duration {
	min := float64(b.opts.min)

	d := min * math.Pow(b.opts.factor, float64(b.attempt))
	d = rand.Float64()*(d-min) + min

	if d > maxInt64 {
		return b.opts.max
	}

	dur := time.Duration(d)

	if dur < b.opts.min {
		return b.opts.min
	}

	if dur > b.opts.max {
		return b.opts.max
	}

	b.attempt++
	return dur
}

func (b *Backoff) Attempt() int {
	return b.attempt
}

func (b *Backoff) Reset() {
	b.attempt = 0
}
