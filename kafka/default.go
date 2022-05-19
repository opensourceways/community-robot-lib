package kafka

import "github.com/opensourceways/community-robot-lib/mq"

var (
	DefaultMQ = NewMQ()
)

func Init(opts ...mq.Option) error {
	return DefaultMQ.Init(opts...)
}

func Connect() error {
	return DefaultMQ.Connect()
}

func Disconnect() error {
	return DefaultMQ.Disconnect()
}

func Publish(topic string, msg *mq.Message, opts ...mq.PublishOption) error {
	return DefaultMQ.Publish(topic, msg, opts...)
}

func Subscribe(topic string, handler mq.Handler, opts ...mq.SubscribeOption) (mq.Subscriber, error) {
	return DefaultMQ.Subscribe(topic, handler, opts...)
}

func String() string {
	return DefaultMQ.String()
}
