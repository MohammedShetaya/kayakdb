package fixtures

import "github.com/MohammedShetaya/kayakdb/api"

// Given takes test data and makes any required validations
type Given struct {
	*Common
}

func (g *Given) Then() *Then {
	return &Then{Common: g.Common}
}

func (g *Given) When() *When {
	return &When{Common: g.Common}
}

func (g *Given) Payload(payload api.Payload) *Given {
	g.options["payload"] = payload
	return g
}
