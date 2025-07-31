package cmd

import (
	"encoding/binary"
	"fmt"
	"github.com/MohammedShetaya/kayakdb/cli/ui"
	"github.com/MohammedShetaya/kayakdb/types"
	"io"
	"net"
	"strconv"
	"strings"
)

func SendRequest(hostname string, port string, payload types.Payload) types.Payload {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", hostname, port))
	if err != nil {
		ui.Error("Connection Error", fmt.Sprintf("Failed to connect to server at %s:%s", hostname, port)).
			WithDetails(
				err.Error(),
				"Possible solutions:",
				"  • Check if the server is running",
				"  • Verify hostname and port are correct",
				"  • Check network connectivity",
			).PrintAndExit()
	}

	data, err := payload.Serialize()
	if err != nil {
		ui.Error("Serialization Error", "Failed to serialize request payload").
			WithDetails(err.Error()).
			PrintAndExit()
	}

	_, err = conn.Write(data)
	if err != nil {
		ui.Error("Network Error", "Failed to send data to server").
			WithDetails(err.Error()).
			PrintAndExit()
	}

	// Signal to the server that we have finished sending the request so it can start processing
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		_ = tcpConn.CloseWrite()
	}

	var resBuffer []byte
	buffer := make([]byte, 1024)

	for {
		n, e := conn.Read(buffer)

		if e != nil {
			if e != io.EOF {
				ui.Error("Buffer Error", "Error occurred while reading payload").WithDetails(err.Error()).PrintAndExit()
			}
			break
		}
		resBuffer = append(resBuffer, buffer[:n]...)
	}

	var res types.Payload
	err = res.Deserialize(resBuffer)
	if err != nil {
		ui.Error("Deserialization Error", "Failed to deserialize payload").
			WithDetails(err.Error()).
			PrintAndExit()
	}

	_ = conn.Close()

	return res
}

// FormatDataTypeError is a helper function to handle data type conversion errors consistently
func FormatDataTypeError(arg string, err error, context string) {
	ui.Error("Data Type Error", fmt.Sprintf("Failed to convert %s to valid data type", context)).
		WithCode(arg).
		WithDetails(
			err.Error(),
			"Supported formats:",
			"  • str:value    - for strings (e.g., str:hello)",
			"  • num:123      - for numbers (e.g., num:42)",
			"  • bool:true    - for booleans (e.g., bool:false)",
			"  • plain text   - auto-detected as number or string",
		).PrintAndExit()
}

func ConvertStringToDataType(data string) (types.Type, error) {
	switch {
	case strings.HasPrefix(data, "str:"):
		value := strings.TrimPrefix(data, "str:")
		return types.String(value), nil

	case strings.HasPrefix(data, "num:"):
		value := strings.TrimPrefix(data, "num:")
		num, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid num: %v", err)
		}
		byteArray := make([]byte, 8)
		binary.BigEndian.PutUint64(byteArray, uint64(num))
		return types.Number(byteArray), nil

	case strings.HasPrefix(data, "bool:"):
		value := strings.TrimPrefix(data, "bool:")
		switch value {
		case "true":
			return types.Bool([]byte{0x01}), nil
		case "false":
			return types.Bool([]byte{0x00}), nil
		default:
			return nil, fmt.Errorf("invalid bool value: %s", value)
		}
	}

	// Default behavior: try number, else string
	if num, err := strconv.Atoi(data); err == nil {
		byteArray := make([]byte, 8)
		binary.BigEndian.PutUint64(byteArray, uint64(num))
		return types.Number(byteArray), nil
	}

	return types.String(data), nil
}
