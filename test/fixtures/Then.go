package fixtures

import (
	"fmt"
	"github.com/MohammedShetaya/kayakdb/api"
	"net"
)

// Then takes any action in for the test to be completed
type Then struct {
	*Common
}

func (t *Then) When() *When {
	return &When{Common: t.Common}
}

// SendRequest
// Expected options: ["payload"]
func (t *Then) SendRequest() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", t.server.Host, t.server.Port))
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		t.Error("Failed to connect to server", err)
	}
	// get the payload from the options
	payload, ok := t.options["payload"].(api.Payload)
	if !ok {
		t.Error("Failed to send payload", fmt.Errorf("payload not found in options"))
	}
	data, err := payload.Serialize()
	_, err = conn.Write(data)
}
