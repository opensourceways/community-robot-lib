package kafka

import (
	"github.com/Shopify/sarama"

	"github.com/opensourceways/community-robot-lib/mq"
)

type event struct {
	m    *mq.Message
	km   *sarama.ConsumerMessage
	err  error
	sess sarama.ConsumerGroupSession
}

func (e *event) Topic() string {
	if e.km != nil {
		return e.km.Topic
	}

	return ""
}

func (e *event) Message() *mq.Message {
	return e.m
}

func (e *event) Ack() error {
	e.sess.MarkMessage(e.km, "")

	return nil
}

func (e *event) Error() error {
	return e.err
}

func (e *event) Extra() map[string]interface{} {
	return map[string]interface{}{
		"time":       e.km.Timestamp,
		"offset":     e.km.Offset,
		"partition":  e.km.Partition,
		"block_time": e.km.BlockTimestamp,
	}
}

type subscriber struct {
	t    string
	cg   sarama.ConsumerGroup
	opts mq.SubscribeOptions
}

func (s *subscriber) Options() mq.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.t
}

func (s *subscriber) Unsubscribe() error {
	return s.cg.Close()
}
