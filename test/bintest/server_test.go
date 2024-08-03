package bintest

import (
	"fmt"
	"net/http"
	"strings"
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
	name := s.T().Name()
	name = strings.ReplaceAll(name, "/", "-")
	SetupCoverDir(fmt.Sprintf("/app/test/bintest/coverdir/%s", name))
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
	defer func() { resp.Body.Close() }()

	s.Require().Equal(http.StatusOK, resp.StatusCode)
}
