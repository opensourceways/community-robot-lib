package mq

import (
	"context"
	"crypto/tls"

	"github.com/sirupsen/logrus"
)

type Options struct {
	Addresses []string
	Secure    bool
	Codec     Codecer

	// Handler executed when error happens in mq message processing
	ErrorHandler Handler

	TLSConfig *tls.Config

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	Log *logrus.Entry
}

type Option func(*Options)

// Addresses set the host addresses to be used by the mq
func Addresses(addrs ...string) Option {
	return func(o *Options) {
		o.Addresses = addrs
	}
}

// Secure communication with the mq
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// Codec sets the codec used for encoding/decoding used where
func Codec(c Codecer) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

// ErrorHandler set the error handler
func ErrorHandler(h Handler) Option {
	return func(o *Options) {
		o.ErrorHandler = h
	}
}

// SetTLSConfig Specify TLS Config
func SetTLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

func Context(c context.Context) Option {
	return func(o *Options) {
		o.Context = c
	}
}

func ContextWithValue(k, v interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, k, v)
	}
}

func Log(log *logrus.Entry) Option {
	return func(o *Options) {
		o.Log = log
	}
}

type PublishOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type PublishOption func(*PublishOptions)

// PublishContext set context
func PublishContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

type SubscribeOptions struct {
	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is receipt already.
	AutoAck bool
	// Subscribers with the same queue name
	// will create a shared subscription where each
	// receives a subset of messages.
	Queue string

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type SubscribeOption func(*SubscribeOptions)

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	opt := SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// DisableAutoAck will disable auto acking of messages
// after they have been handled.
func DisableAutoAck() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = false
	}
}

// Queue sets the name of the queue to share messages on
func Queue(name string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = name
	}
}

// SubscribeContext set context
func SubscribeContext(ctx context.Context) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Context = ctx
	}
}
