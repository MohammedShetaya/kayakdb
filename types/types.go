package types

/*
 This module represents the types used in the api and storage backend
*/

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
)

const (
	MaxPayloadSize uint32 = 1 << 24 // 16MB //todo: make this configurable
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")

// Type This interface represents any data type
type Type interface {
	// Stringer returns the type as a string
	fmt.Stringer
	// Bytes returns the type as bytes
	Bytes() []byte
}

// Bool ------------------------------------------------------------------------------------------------------
// Bool Implementation the Bool types
type Bool []byte

func (b Bool) String() string {
	if binary.BigEndian.Uint32(b) == 1 {
		return "True"
	}
	return "False"
}

func (b Bool) Bytes() []byte {
	return b
}

// Number ------------------------------------------------------------------------------------------------------
// Number TODO: support larger numbers and floating point
// Number Implementation of the Number type
type Number []byte

func (n Number) String() string {
	// Convert the 4 bytes to uint32
	num := int64(binary.BigEndian.Uint64(n))
	return fmt.Sprintf("%d", num)
}

func (n Number) Bytes() []byte {
	return n
}

// String ------------------------------------------------------------------------------------------------------
// String Implementation of the String type
type String string

func (s String) String() string {
	return string(s)
}

func (s String) Bytes() []byte {
	return []byte(s)
}

type Headers struct {
	Path String
	// TODO: fix this to have more meaning if needed. Now it is just sucdess 0 code or fail 1
	Status  uint
	Message string
}

func (h Headers) String() string {
	return fmt.Sprintf("Path: %s (Length: %d)", h.Path, len(h.Path))
}

func (h Headers) Bytes() []byte {
	return h.Path.Bytes()
}

// KeyValue ------------------------------------------------------------------------------------------------------
type KeyValue struct {
	Key   Type
	Value Type
}

func (kv KeyValue) String() string {
	if kv.Value == nil {
		return fmt.Sprintf("%s: nil", kv.Key.String())
	}
	return fmt.Sprintf("%s: %s", kv.Key.String(), kv.Value.String())
}

func (kv KeyValue) Bytes() []byte {
	var buffer bytes.Buffer

	// Write key bytes
	keyBytes := kv.Key.Bytes()
	buffer.Write(keyBytes)

	// Write separator
	buffer.WriteByte(':')

	// Write value bytes
	if kv.Value != nil {
		valueBytes := kv.Value.Bytes()
		buffer.Write(valueBytes)
	}

	return buffer.Bytes()
}

// Payload ------------------------------------------------------------------------------------------------------
type Payload struct {
	Headers Headers
	Data    []Type
}

func (p Payload) String() string {
	var sb strings.Builder
	// Add headers
	sb.WriteString("Headers: ")
	sb.WriteString(p.Headers.String())
	sb.WriteString("\n")

	// Add data
	sb.WriteString("Data:\n")
	if len(p.Data) == 0 {
		sb.WriteString("  (empty)\n")
	} else {
		for i, item := range p.Data {
			sb.WriteString(fmt.Sprintf("  [%d] %s\n", i, item.String()))
		}
	}

	return sb.String()
}

func (p Payload) Bytes() []byte {
	var buffer bytes.Buffer

	// Write headers bytes
	headerBytes := p.Headers.Bytes()
	buffer.Write(headerBytes)

	// Write separator
	buffer.WriteByte('\n')

	// Write data bytes
	for _, kv := range p.Data {
		kvBytes := kv.Bytes()
		buffer.Write(kvBytes)
		buffer.WriteByte('\n')
	}

	return buffer.Bytes()
}

// TODO fix this shit to match the type interface
func (p Payload) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(p); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (p *Payload) Deserialize(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	if err := decoder.Decode(p); err != nil {
		return err
	}

	return nil
}

func RegisterDataTypes() {
	// register the types to be serialized
	gob.Register(Bool{})
	gob.Register(Number{})
	gob.Register(String(""))
	gob.Register(KeyValue{})
	gob.Register(Headers{})
	gob.Register(Payload{})
}
