package fixtures

import (
	"github.com/MohammedShetaya/kayakdb/api"
	"github.com/MohammedShetaya/kayakdb/config"
	. "github.com/MohammedShetaya/kayakdb/test/fixtures/test_data"
	"github.com/MohammedShetaya/kayakdb/utils"
	"github.com/stretchr/testify/suite"
	"time"
)

type E2ESuite struct {
	suite.Suite
	Common
}

func (s *E2ESuite) SetupSuite() {
	s.Common.t = s.Suite.T()
	s.options = make(map[string]interface{})
	// start the server
	logger := utils.InitLogger("debug")
	defer func() {
		_ = logger.Sync()
	}()
	// Start the server in a separate goroutine
	go func() {
		cfg := &config.Configuration{KayakPort: KayakdbPort}
		s.Common.server = api.NewServer(cfg, logger)
		s.Common.server.Start()
	}()

	// wait for the server to start
	// TODO: replace this with a better way to wait for the server after implementing responses from the server
	for s.Common.server == nil {
		time.Sleep(time.Second)
	}
}

func (s *E2ESuite) Given() *Given {
	return &Given{
		Common: &s.Common,
	}
}
