package cmd

import (
	"encoding/binary"
	"fmt"
	"github.com/MohammedShetaya/kayakdb/api"
	"net"
	"strconv"
)

func ConvertToDataType(data string) (api.Type, error) {
	if num, err := strconv.Atoi(data); err == nil {
		byteArray := make([]byte, 8)
		binary.BigEndian.PutUint64(byteArray, uint64(num))
		return api.Number(byteArray), nil
	} else {
		return api.String(data), nil
	}
}

func SendRequest(hostname string, port string, payload api.Payload) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", hostname, port))
	if err != nil {
		fmt.Errorf("Failed to connect to server: %v", err)
	}
	data, err := payload.Serialize()
	if err != nil {
		fmt.Errorf("Failed to serialize payload: %v", err)
	}
	_, err = conn.Write(data)
	conn.Close()
}
