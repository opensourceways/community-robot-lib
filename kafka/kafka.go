package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/community-robot-lib/mq"
	"github.com/opensourceways/community-robot-lib/utils"
)

type kfkMQ struct {
	opts           mq.Options
	producer       sarama.SyncProducer
	consumers      []sarama.Client
	consumerGroups [] sarama.ConsumerGroup
	mutex          sync.RWMutex
	connected      bool
	log            *logrus.Entry
}

func (kMQ *kfkMQ) Init(opts ...mq.Option) error {
	kMQ.mutex.RLock()
	if kMQ.connected {
		return fmt.Errorf("mq is connected can't init")
	}
	kMQ.mutex.RUnlock()

	for _, o := range opts {
		o(&kMQ.opts)
	}

	kMQ.log = kMQ.opts.Log

	if kMQ.opts.Addresses == nil {
		kMQ.opts.Addresses = []string{"127.0.0.1:9092"}
	}

	if kMQ.opts.Context == nil {
		kMQ.opts.Context = context.Background()
	}

	if kMQ.opts.Codec == nil {
		kMQ.opts.Codec = mq.JsonCodec{}
	}

	return nil
}

func (kMQ *kfkMQ) Options() mq.Options {
	return kMQ.opts
}

func (kMQ *kfkMQ) Address() string {
	if len(kMQ.opts.Addresses) > 0 {
		return kMQ.opts.Addresses[0]
	}

	return ""
}

func (kMQ *kfkMQ) Connect() error {
	if kMQ.connected {
		return nil
	}

	kMQ.mutex.RLock()
	if kMQ.producer != nil {
		kMQ.mutex.RUnlock()

		return nil
	}
	kMQ.mutex.RUnlock()

	producer, err := sarama.NewSyncProducer(kMQ.opts.Addresses, kMQ.clusterConfig())
	if err != nil {
		return err
	}

	kMQ.mutex.Lock()
	kMQ.producer = producer
	kMQ.connected = true
	kMQ.mutex.Unlock()

	return nil
}

func (kMQ *kfkMQ) Disconnect() error {
	kMQ.mutex.Lock()
	defer kMQ.mutex.Unlock()

	mErr := utils.MultiError{}
	if kMQ.connected {
		mErr.AddError(kMQ.producer.Close())
	}

	for _, g := range kMQ.consumerGroups {
		mErr.AddError(g.Close())
	}

	for _, c := range kMQ.consumers {
		if !c.Closed() {
			mErr.AddError(c.Close())
		}
	}

	kMQ.connected = false

	return mErr.Err()
}

// Publish a message to a topic in the kafka cluster.
func (kMQ *kfkMQ) Publish(topic string, msg *mq.Message, opts ...mq.PublishOption) error {
	d, err := kMQ.opts.Codec.Marshal(msg)
	if err != nil {
		return err
	}

	pm := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(d),
	}

	if key := msg.MessageKey(); key != "" {
		pm.Key = sarama.StringEncoder(key)
	}

	_, _, err = kMQ.producer.SendMessage(pm)

	return err
}

// Subscribe to kafka message topics, each subscription generates a kafka groupConsumer group.
func (kMQ *kfkMQ) Subscribe(topics string, h mq.Handler, opts ...mq.SubscribeOption) (mq.Subscriber, error) {
	opt := mq.SubscribeOptions{
		AutoAck: true,
		Queue:   uuid.New().String(),
	}
	for _, o := range opts {
		o(&opt)
	}
	c, err := kMQ.saramaClusterClient()
	if err != nil {
		return nil, err
	}

	g, err := sarama.NewConsumerGroupFromClient(opt.Queue, c)
	if err != nil {
		return nil, err
	}

	gc := &groupConsumer{
		handler: h,
		subOpts: opt,
		kOpts:   kMQ.opts,
		ready:   make(chan bool),
	}

	go func() {
		for {
			select {
			case err := <-g.Errors():
				if err != nil {
					kMQ.log.Errorf("consumer error: %v", err)
				}
			default:
				err := g.Consume(kMQ.opts.Context, []string{topics}, gc)
				switch err {
				case sarama.ErrClosedConsumerGroup:
					return
				case nil:
					continue
				default:
					kMQ.log.Error(err)
				}
			}

			if kMQ.opts.Context.Err() != nil {
				return
			}
		}
	}()

	<-gc.ready

	kMQ.mutex.Lock()
	kMQ.consumerGroups = append(kMQ.consumerGroups, g)
	kMQ.mutex.Unlock()

	return &subscriber{cg: g, t: topics, opts: opt}, nil
}

func (kMQ *kfkMQ) String() string {
	return "kafka"
}

func (kMQ *kfkMQ) clusterConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 3

	if kMQ.opts.TLSConfig != nil {
		cfg.Net.TLS.Config = kMQ.opts.TLSConfig
		cfg.Net.TLS.Enable = true
	}

	if !cfg.Version.IsAtLeast(sarama.MaxVersion) {
		cfg.Version = sarama.MaxVersion
	}

	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	return cfg
}

func (kMQ *kfkMQ) saramaClusterClient() (sarama.Client, error) {
	cs, err := sarama.NewClient(kMQ.opts.Addresses, kMQ.clusterConfig())
	if err != nil {
		return nil, err
	}

	kMQ.mutex.Lock()
	defer kMQ.mutex.Unlock()
	kMQ.consumers = append(kMQ.consumers, cs)

	return cs, nil
}

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
	for msg := range claim.Messages() {
		var m mq.Message

		ke := &kEvent{km: msg, m: &m, sess: session}
		eHandler := gc.kOpts.ErrorHandler

		if err := gc.kOpts.Codec.Unmarshal(msg.Value, &m); err != nil {
			ke.err = err
			ke.m.Body = msg.Value

			if eHandler == nil {
				gc.kOpts.Log.Errorf("unmarshal kafka msg fail with error : %v", err)

				continue
			}

			if err := eHandler(ke); err != nil {
				gc.kOpts.Log.Error(err)
			}
		}

		if gc.handler == nil {
			gc.handler = func(event mq.Event) error {
				return fmt.Errorf("msg handler func is nil")
			}
		}

		if err := gc.handler(ke); err != nil {
			ke.err = err

			if eHandler == nil {
				gc.kOpts.Log.Errorf("subscriber error: %v", err)

				continue
			}

			if err := eHandler(ke); err != nil {
				gc.kOpts.Log.Error(err)
			}
		}

		if gc.subOpts.AutoAck {
			_ = ke.Ack()
		}

	}

	return nil
}

type kEvent struct {
	err error
	km  *sarama.ConsumerMessage
	m   *mq.Message

	sess sarama.ConsumerGroupSession
}

func (ke *kEvent) Topic() string {
	if ke.km != nil {
		return ke.km.Topic
	}

	return ""
}

func (ke *kEvent) Message() *mq.Message {
	return ke.m
}

func (ke *kEvent) Ack() error {
	ke.sess.MarkMessage(ke.km, "")

	return nil
}

func (ke *kEvent) Error() error {
	return ke.err
}

func (ke *kEvent) Extra() map[string]interface{} {
	em := make(map[string]interface{})
	em["offset"] = ke.km.Offset
	em["partition"] = ke.km.Partition
	em["time"] = ke.km.Timestamp
	em["block_time"] = ke.km.BlockTimestamp

	return em
}

type subscriber struct {
	cg   sarama.ConsumerGroup
	t    string
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

func NewMQ(opts ...mq.Option) mq.MQ {
	options := mq.Options{
		Codec:   mq.JsonCodec{},
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	if len(options.Addresses) == 0 {
		options.Addresses = []string{"127.0.0.1:9092"}
	}

	if options.Log == nil {
		options.Log = logrus.New().WithField("function", "kafka mq")
	}

	return &kfkMQ{
		opts: options,
	}
}
