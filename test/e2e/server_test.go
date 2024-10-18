package e2e

import (
	"github.com/MohammedShetaya/kayakdb/test/fixtures"
	"github.com/MohammedShetaya/kayakdb/test/fixtures/test_data"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServerSuite struct {
	fixtures.E2ESuite
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

func (s *ServerSuite) TestServerCanReceivesPayload() {
	s.Given().
		Payload(test_data.GetPayload).
		Then().
		SendRequest()
}
