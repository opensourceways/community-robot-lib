package mq

import "encoding/json"

//Codec is a simple encoding interface used for the mq/transport
type Codecer interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}

type JsonCodec struct{}

func (jc JsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (jc JsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (jc JsonCodec) String() string {
	return "json"
}