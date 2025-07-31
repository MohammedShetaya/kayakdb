package fixtures

import (
	"github.com/MohammedShetaya/kayakdb/types"
)

// Given takes test types and makes any required validations
type Given struct {
	*Common
}

func (g *Given) Then() *Then {
	return &Then{Common: g.Common}
}

func (g *Given) When() *When {
	return &When{Common: g.Common}
}

func (g *Given) Payload(payload types.Payload) *Given {
	g.options["payload"] = payload
	return g
}
