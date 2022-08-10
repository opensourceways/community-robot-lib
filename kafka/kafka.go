package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/community-robot-lib/mq"
)

type kfkMQ struct {
	opts     mq.Options
	producer sarama.SyncProducer

	mutex     sync.RWMutex
	connected bool
}

func (kMQ *kfkMQ) Init(opts ...mq.Option) error {
	if kMQ.isConnected() {
		return fmt.Errorf("mq is connected and can't do init")
	}

	kMQ.mutex.Lock()
	defer kMQ.mutex.Unlock()

	if kMQ.connected {
		return nil
	}

	for _, o := range opts {
		o(&kMQ.opts)
	}

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

func (kMQ *kfkMQ) isConnected() (b bool) {
	kMQ.mutex.RLock()
	b = kMQ.connected
	kMQ.mutex.RUnlock()

	return
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
	if kMQ.isConnected() {
		return nil
	}

	kMQ.mutex.Lock()
	defer kMQ.mutex.Unlock()

	if kMQ.connected {
		return nil
	}

	producer, err := sarama.NewSyncProducer(kMQ.opts.Addresses, kMQ.clusterConfig())
	if err != nil {
		return err
	}

	kMQ.producer = producer
	kMQ.connected = true

	return nil
}

func (kMQ *kfkMQ) Disconnect() error {
	if !kMQ.isConnected() {
		return nil
	}

	kMQ.mutex.Lock()
	defer kMQ.mutex.Unlock()

	if !kMQ.connected {
		return nil
	}

	kMQ.connected = false

	return kMQ.producer.Close()
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
		c.Close()

		return nil, err
	}

	gc := groupConsumer{
		handler: h,
		subOpts: opt,
		kOpts:   kMQ.opts,
	}

	s := newSubscriber(topics, c, g, gc)
	s.start()

	return s, nil
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
	return sarama.NewClient(kMQ.opts.Addresses, kMQ.clusterConfig())
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
