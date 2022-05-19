/*
Copyright 2022 The go-micro Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package mq is an interface used for asynchronous messaging, the default implementation is kafka
package mq

type MQ interface {
	Init(...Option) error
	Options() Options
	Address() string
	Connect() error
	Disconnect() error
	Publish(topic string, m *Message, opts ...PublishOption) error
	Subscribe(topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)
	String() string
}

// Event is given to a subscription handler for processing.
type Event interface {
	// Topic return the topic of the message
	Topic() string
	// Message return the message body
	Message() *Message
	// ACK message reply operation
	Ack() error
	// Error get the error of message consumption
	Error() error
	// Extra the important information other than the message body
	Extra() map[string]interface{}
}

// Subscriber is a convenience return type for the Subscribe method.
type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
}

// Message is the message entity.
type Message struct {
	key    string
	Header map[string]string
	Body   []byte
}

// SetMessageKey set a flag that represents the message
func (msg *Message) SetMessageKey(key string) {
	msg.key = key
}

// MessageKey get the flag that represents the message
func (msg Message) MessageKey() string {
	return msg.key
}

// Handler is used to process messages via a subscription of a topic.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type Handler func(Event) error
