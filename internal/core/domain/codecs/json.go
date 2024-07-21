package codecs

import (
	"encoding/json"
)

// JSONCodec is a codec that uses JSON
type JSONCodec struct{}

// NewJSONCodec creates a new JSONCodec
func NewJSONCodec() *JSONCodec {
	return &JSONCodec{}
}

// Encode marshals the given value
func (c *JSONCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Decode unmarshals the given data into the given value
func (c *JSONCodec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
