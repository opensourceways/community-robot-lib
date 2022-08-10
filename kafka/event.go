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
	km := e.km
	if km == nil {
		return nil
	}

	return map[string]interface{}{
		"time":       km.Timestamp,
		"offset":     km.Offset,
		"partition":  km.Partition,
		"block_time": km.BlockTimestamp,
	}
}
