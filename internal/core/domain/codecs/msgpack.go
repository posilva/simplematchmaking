package codecs

import (
	"github.com/vmihailenco/msgpack/v5"
)

// MsgPackCodec is a codec that uses msgpack
type MsgPackCodec struct{}

// NewMsgPackCodec creates a new MsgPackCodec
func NewMsgPackCodec() *MsgPackCodec {
	return &MsgPackCodec{}
}

// Encode marshals the given value
func (c *MsgPackCodec) Encode(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

// Decode unmarshals the given data into the given value
func (c *MsgPackCodec) Decode(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}
