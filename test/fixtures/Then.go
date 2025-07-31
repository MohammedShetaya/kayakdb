package fixtures

import (
	"fmt"
	"github.com/MohammedShetaya/kayakdb/types"
	. "github.com/MohammedShetaya/kayakdb/test/fixtures/test_data"
	"net"
	"time"
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
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", KayakdbHost, KayakdbPort))
	if err != nil {
		t.Error("Failed to connect to server", err)
	}
	// get the payload from the options
	payload, ok := t.options["payload"].(types.Payload)
	if !ok {
		t.Error("Failed to send payload", fmt.Errorf("payload not found in options"))
	}
	data, err := payload.Serialize()
	_, err = conn.Write(data)
	conn.Close()
	// TODO: remove this sleep after implementing responses from the server
	time.Sleep(time.Second)
}
