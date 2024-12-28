package fixtures

import (
	"github.com/MohammedShetaya/kayakdb/api"
	"testing"
)

// Common holds common logic and types of tests
type Common struct {
	server  *api.Server
	options map[string]interface{}
	t       *testing.T
}

func (c *Common) Error(msg string, err error) {
	if err != nil {
		c.t.Fatalf("%v\n%v", msg, err)
	}
}
