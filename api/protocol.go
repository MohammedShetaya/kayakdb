package api

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"strings"
)

/*
This represents the TCP protocol specification that is used in the api endpoints
*/
const (
	// BinaryType Binary constant that represents the binary type
	BinaryType uint8 = iota
	// NumberType Number NumberType constant that represents the numbers type
	NumberType
	// StringType String constant that represents the string type
	StringType

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
	io.ReaderFrom
}

// Binary Implementation the Binary data type
type Binary []byte

func (b Binary) String() string {
	return string(b)
}

func (b Binary) Bytes() []byte {
	return b
}

// ReadFrom reads from am io.reader for example a socket
func (b *Binary) ReadFrom(r io.Reader) (int64, error) {
	// read one byte
	*b = make([]byte, 1)
	err := binary.Read(r, binary.BigEndian, *b)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// Number TODO: support larger numbers and floating point
// Number Implementation of the Number data type
type Number []byte

func (n Number) String() string {
	// Convert the 4 bytes to uint32
	num := int32(binary.BigEndian.Uint32(n))
	return fmt.Sprintf("%d", num)
}

func (n Number) Bytes() []byte {
	return n
}

func (n *Number) ReadFrom(r io.Reader) (int64, error) {
	// read 4 bytes
	*n = make([]byte, 4)
	err := binary.Read(r, binary.BigEndian, *n)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

// String Implementation of the String data type
type String []byte

func (s String) String() string {
	return string(s)
}

func (s String) Bytes() []byte {
	return s
}

// ReadFrom reads a string from the io.Reader
func (s *String) ReadFrom(r io.Reader) (int64, error) {
	// Read the first 4 bytes which represent the string length
	var length uint32
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return 0, err
	}

	// Check if the length exceeds the maximum payload size
	if length > MaxPayloadSize {
		return 4, ErrMaxPayloadSize
	}

	// Read the string data based on the length
	*s = make([]byte, length)
	n, err := io.ReadFull(r, *s)
	if err != nil {
		return int64(n + 4), err
	}

	return int64(n + 4), nil
}

type Headers struct {
	// number of bytes for the path
	PathLength uint32
	// the actual path
	Path string
}

func (h Headers) String() string {
	return fmt.Sprintf("Path: %s (Length: %d)", h.Path, h.PathLength)
}

type Payload struct {
	Headers Headers
	Data    map[Type]Type
}

func (p Payload) String() string {
	var sb strings.Builder

	// Write headers to the string
	sb.WriteString(fmt.Sprintf("Headers: %s\n", p.Headers.String()))

	// Iterate over Data and write each key-value pair
	sb.WriteString("Data:\n")
	for key, value := range p.Data {
		sb.WriteString(fmt.Sprintf("%s: %s\n", key.String(), value.String()))
	}

	return sb.String()
}

// decodes and returns Payload from an io.Reader
func decode(logger *zap.Logger, r io.Reader) (*Payload, error) {
	headers := Headers{}
	// read the path chunk
	pathLength := make([]byte, 4)
	n, err := r.Read(pathLength)
	if err != nil {
		logger.Error("Unable to read path length")
		return nil, err
	}
	if n <= 0 {
		logger.Error("Path length is less than or equal 0")
		return nil, nil
	}
	err = binary.Read(bytes.NewReader(pathLength), binary.BigEndian, &headers.PathLength)
	if err != nil {
		logger.Error("Unable to convert path len byte array to int")
		return nil, err
	}

	logger.Info("Path Length", zap.Int("Path Length is", int(headers.PathLength)))

	// read the path
	path := make([]byte, headers.PathLength)
	n, err = r.Read(path)
	if err != nil {
		logger.Error("Unable to read path")
		return nil, err
	}
	if uint32(n) != headers.PathLength {
		logger.Error("Path length is less than or equal 0")
		return nil, io.ErrUnexpectedEOF
	}
	headers.Path = string(path)
	// TODO: implement different handlers based on the bath
	logger.Info("Path is:", zap.String("Path", headers.Path))
	// Initialize the Payload
	payload := &Payload{
		Headers: headers,
		Data:    map[Type]Type{},
	}

	totalDataSize := 0

	for {
		// Read key type (1 byte)
		var keyType uint8
		if err := binary.Read(r, binary.BigEndian, &keyType); err != nil {
			if err == io.EOF {
				break // Exit when all data is read
			}
			logger.Error("Failed to read key type", zap.Error(err))
			return nil, err
		}
		logger.Debug("", zap.Int("key type", int(keyType)))

		// Read value type (1 byte)
		var valueType uint8
		if err := binary.Read(r, binary.BigEndian, &valueType); err != nil {
			logger.Error("Failed to read value type", zap.Error(err))
			return nil, err
		}

		logger.Info("value type", zap.Int("value type", int(valueType)))
		// Track total data size
		totalDataSize += 2 // 1 for keyType, 1 for valueType

		// Initialize the key based on the key type
		keyData, err := readType(logger, r, keyType)
		if err != nil {
			logger.Error("Cannot read key data", zap.Error(err))
		}

		// Initialize the value based on the value type
		valueData, err := readType(logger, r, valueType)
		if err != nil {
			logger.Error("Cannot read value data", zap.Error(err))
		}

		totalDataSize += len(keyData.Bytes()) + len(valueData.Bytes())
		if totalDataSize > int(MaxPayloadSize) {
			logger.Error("Payload exceeds MaxPayloadSize")
			return nil, ErrMaxPayloadSize
		}
		// Append the key and value to the payload data
		payload.Data[keyData] = valueData
	}

	return payload, nil
}

func readType(logger *zap.Logger, reader io.Reader, dataType uint8) (Type, error) {
	var data Type

	switch dataType {
	case BinaryType:
		data = &Binary{}
		if _, err := data.ReadFrom(reader); err != nil {
			logger.Error("Failed to read key as Binary", zap.Error(err))
			return nil, err
		}
	case StringType:
		data = &String{}
		if _, err := data.ReadFrom(reader); err != nil {
			logger.Error("Failed to read key as String", zap.Error(err))
			return nil, err
		}
	case NumberType:
		data = &Number{}
		if _, err := data.ReadFrom(reader); err != nil {
			logger.Error("Failed to read key as Number", zap.Error(err))
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown key type: %d", dataType)
	}

	return data, nil
}
