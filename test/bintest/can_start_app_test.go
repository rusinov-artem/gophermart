package bintest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServerTestsuite struct {
	suite.Suite
	server ServerUnderTest
}

func Test_Server(t *testing.T) {
	Run(t, &ServerTestsuite{})
}

func (s *ServerTestsuite) SetupTest() {
	s.server = *NewServerUnderTest("./app")
	err := s.server.Start()
	s.Require().NoError(err)
}

func (s *ServerTestsuite) TearDownTest() {
	err := s.server.Stop()
	s.Require().NoError(err)
}

func (s *ServerTestsuite) Test_CanStartServer() {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:7777/liveness", nil)
	s.Require().NoError(err)
	client := http.DefaultClient
	resp, err := client.Do(req)
	s.Require().NoError(err)
	defer func() {
		resp.Body.Close()
	}()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}
