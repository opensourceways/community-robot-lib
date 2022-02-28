package mq

import (
	"context"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/community-robot-lib/utils"
)

type kfkMQ struct {
	opts           Options
	producer       sarama.SyncProducer
	consumers      []sarama.Client
	consumerGroups [] sarama.ConsumerGroup
	mutex          sync.RWMutex
	connected      bool
	log            *logrus.Entry
}

func (kMQ *kfkMQ) Init(opts ...Option) error {
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
		kMQ.opts.Codec = JsonCodec{}
	}

	return nil
}

func (kMQ *kfkMQ) Options() Options {
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
func (kMQ *kfkMQ) Publish(topic string, msg *Message, opts ...PublishOption) error {
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

// Subscribe to kafka message topics, each subscription generates a kafka consumer group.
func (kMQ *kfkMQ) Subscribe(topics string, h Handler, opts ...SubscribeOption) (Subscriber, error) {
	opt := SubscribeOptions{
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

	cg := &consumerGroup{
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
				err := g.Consume(kMQ.opts.Context, []string{topics}, cg)
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

			cg.ready = make(chan bool)
		}
	}()

	<-cg.ready

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

// consumerGroup represents a Sarama consumer group consumer
type consumerGroup struct {
	handler Handler
	subOpts SubscribeOptions
	kOpts   Options
	sess    sarama.ConsumerGroupSession

	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (cg *consumerGroup) Setup(sarama.ConsumerGroupSession) error {
	close(cg.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (cg *consumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (cg *consumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var m Message

		ke := &kEvent{km: msg, m: &m, sess: session}
		eHandler := cg.kOpts.ErrorHandler

		if err := cg.kOpts.Codec.Unmarshal(msg.Value, &m); err != nil {
			ke.err = err
			ke.m.Body = msg.Value

			if eHandler == nil {
				cg.kOpts.Log.Errorf("unmarshal kafka msg fail with error : %v", err)

				continue
			}

			if err := eHandler(ke); err != nil {
				cg.kOpts.Log.Error(err)
			}
		}

		if cg.handler == nil {
			cg.handler = func(event Event) error {
				return fmt.Errorf("msg handler func is nil")
			}
		}

		if err := cg.handler(ke); err != nil {
			ke.err = err

			if eHandler == nil {
				cg.kOpts.Log.Errorf("subscriber error: %v", err)

				continue
			}

			if err := eHandler(ke); err != nil {
				cg.kOpts.Log.Error(err)
			}
		}

		if cg.subOpts.AutoAck {
			_ = ke.Ack()
		}

	}

	return nil
}

type kEvent struct {
	err error
	km  *sarama.ConsumerMessage
	m   *Message

	sess sarama.ConsumerGroupSession
}

func (ke *kEvent) Topic() string {
	if ke.km != nil {
		return ke.km.Topic
	}

	return ""
}

func (ke *kEvent) Message() *Message {
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
	opts SubscribeOptions
}

func (s *subscriber) Options() SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.t
}

func (s *subscriber) Unsubscribe() error {
	return s.cg.Close()
}

func NewMQ(opts ...Option) MQ {
	options := Options{
		Codec:   JsonCodec{},
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
