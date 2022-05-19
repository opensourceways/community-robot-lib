package kafka

import (
	"github.com/opensourceways/community-robot-lib/mq"
	"testing"
)

func TestBroker(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("mq init error: %v", err)
	}

	if err := Connect(); err != nil {
		t.Fatalf("mq connect error: %v", err)
	}

	msg := mq.Message{
		Header: map[string]string{
			"Content-type": "application/json",
		},
		Body: []byte(`{"message":"broker_test"}`),
	}
	done := make(chan bool)

	sub, err := Subscribe("mq-test", func(event mq.Event) error {
		m := event.Message()
		if string(m.Body) != string(msg.Body) {
			t.Fatalf("Unexpected msg %s, expected %s", string(m.Body), string(msg.Body))
		}

		t.Logf("message head: %v , body: %s , extra: %v", m.Header, string(m.Body), event.Extra())

		close(done)

		return nil
	})
	if err != nil {
		t.Fatalf("Unexpected subscribe error: %v", err)
	}

	if err := Publish("mq-test", &msg); err != nil {
		t.Fatalf("Unexpected publish error: %v", err)
	}

	<-done
	_ = sub.Unsubscribe()

	if err := Disconnect(); err != nil {
		t.Fatalf("Unexpected disconnect error: %v", err)
	}
}

func TestTwoPartitionMultipleConsumerWithSameKey(t *testing.T) {
	// note: the topic of xwz has 2 partitions

	if err := Init(); err != nil {
		t.Fatalf("mq init error: %v", err)
	}

	if err := Connect(); err != nil {
		t.Fatalf("mq connect error: %v", err)
	}

	msg := mq.Message{
		Header: map[string]string{
			"Content-type": "application/json",
		},
		Body: []byte(`{"message":"broker_test"}`),
	}
	msg.SetMessageKey("2")

	shutdownMsg := mq.Message{
		Header: map[string]string{
			"Content-type": "application/json",
		},
		Body: []byte(`{"message":"shutdown"}`),
	}
	shutdownMsg.SetMessageKey("2")

	done := make(chan bool)
	sub1HasConsumed, sub2HasConsumed := false, false

	if _, err := Subscribe("xwz", func(event mq.Event) error {
		sub1HasConsumed = true
		m := event.Message()
		if m == nil {
			t.Fatal("msg is nil")
		}

		if string(m.Body) == `{"message":"shutdown"}` {
			close(done)
		}

		t.Logf("sub 1 get a msg -- > message head: %v , body: %s , extra: %v", m.Header, string(m.Body), event.Extra())

		return nil
	}, mq.Queue("test-xwz")); err != nil {
		t.Fatalf("Unexpected subscribe error: %v", err)
	}

	if _, err := Subscribe("xwz", func(event mq.Event) error {
		sub2HasConsumed = true
		m := event.Message()
		if m == nil {
			t.Fatal("msg is nil")
		}

		if string(m.Body) == `{"message":"shutdown"}` {
			close(done)
		}

		t.Logf("sub 2 get a msg -- > message head: %v , body: %s , extra: %v", m.Header, string(m.Body), event.Extra())

		return nil
	}, mq.Queue("test-xwz")); err != nil {
		t.Fatalf("Unexpected subscribe error: %v", err)
	}

	if err := Publish("xwz", &msg); err != nil {
		t.Fatalf("Unexpected publish error: %v", err)
	}

	if err := Publish("xwz", &shutdownMsg); err != nil {
		t.Fatalf("Unexpected publish error: %v", err)
	}

	<-done

	if err := Disconnect(); err != nil {
		t.Fatalf("Unexpected disconnect error: %v", err)
	}

	// only one groupConsumer can consume the msg that the msg in the same partition
	if sub1HasConsumed && sub2HasConsumed {
		t.Fatalf("Unexpected consumed logic")
	}
}

func TestTwoPartitionMultipleConsumerWithDiffKey(t *testing.T) {
	//note: the topic of xwz has 2 partitions

	if err := Init(); err != nil {
		t.Fatalf("mq init error: %v", err)
	}

	if err := Connect(); err != nil {
		t.Fatalf("mq connect error: %v", err)
	}

	msg := mq.Message{
		Header: map[string]string{
			"Content-type": "application/json",
		},
		Body: []byte(`{"message":"broker_test"}`),
	}
	msg.SetMessageKey("1")

	shutdownMsg := mq.Message{
		Header: map[string]string{
			"Content-type": "application/json",
		},
		Body: []byte(`{"message":"shutdown"}`),
	}
	shutdownMsg.SetMessageKey("2")

	done := make(chan bool)
	sub1HasConsumed, sub2HasConsumed := false, false

	if _, err := Subscribe("xwz", func(event mq.Event) error {
		sub1HasConsumed = true
		m := event.Message()
		if m == nil {
			t.Fatal("msg is nil")
		}

		if string(m.Body) == `{"message":"shutdown"}` {
			close(done)
		}

		t.Logf("sub 1 get a msg -- > message head: %v , body: %s , extra: %v", m.Header, string(m.Body), event.Extra())

		return nil
	}, mq.Queue("test-xwz")); err != nil {
		t.Fatalf("Unexpected subscribe error: %v", err)
	}

	if _, err := Subscribe("xwz", func(event mq.Event) error {
		sub2HasConsumed = true
		m := event.Message()
		if m == nil {
			t.Fatal("msg is nil")
		}

		if string(m.Body) == `{"message":"shutdown"}` {
			close(done)
		}

		t.Logf("sub 2 get a msg -- > message head: %v , body: %s , extra: %v", m.Header, string(m.Body), event.Extra())

		return nil
	}, mq.Queue("test-xwz")); err != nil {
		t.Fatalf("Unexpected subscribe error: %v", err)
	}

	if err := Publish("xwz", &msg); err != nil {
		t.Fatalf("Unexpected publish error: %v", err)
	}

	if err := Publish("xwz", &shutdownMsg); err != nil {
		t.Fatalf("Unexpected publish error: %v", err)
	}

	<-done

	if err := Disconnect(); err != nil {
		t.Fatalf("Unexpected disconnect error: %v", err)
	}

	// both consumers can consume messages because there are messages in both partitions
	if !sub1HasConsumed || !sub2HasConsumed {
		t.Fatalf("Unexpected consumed logic")
	}
}
