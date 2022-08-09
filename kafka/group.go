package kafka

import (
	"sync"

	"github.com/Shopify/sarama"

	"github.com/opensourceways/community-robot-lib/mq"
)

// groupConsumer represents a Sarama consumer group consumer
type groupConsumer struct {
	handler mq.Handler
	subOpts mq.SubscribeOptions
	kOpts   mq.Options
	sess    sarama.ConsumerGroupSession

	ready chan bool
	once  sync.Once
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (gc *groupConsumer) Setup(sarama.ConsumerGroupSession) error {
	gc.once.Do(func() {
		close(gc.ready)
	})

	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (gc *groupConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a groupConsumer loop of ConsumerGroupClaim's Messages().
func (gc *groupConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	if gc.handler == nil {
		gc.handler = func(event mq.Event) error {
			return nil
		}
	}

	eh := gc.kOpts.ErrorHandler

	for msg := range claim.Messages() {
		var m mq.Message

		ke := &event{
			km:   msg,
			m:    &m,
			sess: session,
		}

		if err := gc.kOpts.Codec.Unmarshal(msg.Value, &m); err != nil {
			ke.err = err
			ke.m.Body = msg.Value

			if eh == nil {
				gc.kOpts.Log.Errorf("unmarshal kafka msg fail with error : %v", err)

				continue
			}

			if err := eh(ke); err != nil {
				gc.kOpts.Log.Error(err)
			}

		} else if err := gc.handler(ke); err != nil {
			ke.err = err

			if eh == nil {
				gc.kOpts.Log.Errorf("subscriber error: %v", err)

				continue
			}

			if err := eh(ke); err != nil {
				gc.kOpts.Log.Error(err)
			}
		}

		if gc.subOpts.AutoAck {
			_ = ke.Ack()
		}
	}

	return nil
}
