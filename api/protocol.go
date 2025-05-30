package api

/*
This module represents the TCP protocol specification that is used in the api endpoints
*/

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	MaxPayloadSize uint32 = 1 << 24 // 16MB //todo: make this configurable
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")

// Type This interface represents the payload to be sent to any endpoint
type Type interface {
	// Stringer returns the type as a string
	fmt.Stringer
	// Bytes returns the type as bytes
	Bytes() []byte
	// ReaderFrom reads from an io.Reader
	//io.ReaderFrom
}

// Binary Implementation the Binary data type
type Binary []byte

func (b Binary) String() string {
	if binary.BigEndian.Uint32(b) == 1 {
		return "True"
	}
	return "False"
}

func (b Binary) Bytes() []byte {
	return b
}

// Number TODO: support larger numbers and floating point
// Number Implementation of the Number data type
type Number []byte

func (n Number) String() string {
	// Convert the 4 bytes to uint32
	num := int64(binary.BigEndian.Uint64(n))
	return fmt.Sprintf("%d", num)
}

func (n Number) Bytes() []byte {
	return n
}

// String Implementation of the String data type
type String []byte

func (s String) String() string {
	return string(s)
}

func (s String) Bytes() []byte {
	return s
}

type Headers struct {
	Path string
}

func (h Headers) String() string {
	return fmt.Sprintf("Path: %s (Length: %d)", h.Path, len(h.Path))
}

type KeyValue struct {
	Key   Type
	Value Type
}

type Payload struct {
	Headers Headers
	Data    []KeyValue
}

func NewPayload() *Payload {
	payload := new(Payload)
	return payload
}

func (p Payload) String() string {
	var sb strings.Builder
	// Write headers to the string
	sb.WriteString(fmt.Sprintf("Headers:\n  %s\n", p.Headers.String()))

	// Iterate over Data and write each key-value pair
	sb.WriteString("Data:\n")
	for _, entry := range p.Data {
		if entry.Value == nil {
			sb.WriteString(fmt.Sprintf("  %s: nil\n", entry.Key.String()))
			continue
		}
		sb.WriteString(fmt.Sprintf("%s: %s\n", entry.Key.String(), entry.Value.String()))
	}

	return sb.String()
}

func (p *Payload) Serialize() ([]byte, error) {
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

func InitProtocol() {
	// register the types to be serialized
	gob.Register(Binary{})
	gob.Register(Number{})
	gob.Register(String{})
}

func ConvertStringKeyToDataType(data string) (Type, error) {
	// convert the string to number
	if num, err := strconv.Atoi(data); err == nil {
		byteArray := make([]byte, 8)
		binary.BigEndian.PutUint64(byteArray, uint64(num))
		return Number(byteArray), nil
	} else { // if not a number then convert to a string data-type
		return String(data), nil
	}
}
