package bintest

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/test/utils/writer"
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
}

func (s *ServerTestsuite) startServer(address string) {
	s.server = *NewServerUnderTest("./app", "-a", address)
	err := s.server.Start("Listening")
	s.Require().NoError(err)
}

func (s *ServerTestsuite) stopServer() {
	err := s.server.Stop()
	s.Require().NoError(err)
}

func (s *ServerTestsuite) Test_CanStartServer() {
	address := "127.0.0.1:8080"
	s.startServer(address)
	defer s.stopServer()
	url := fmt.Sprintf("http://%s/liveness", address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	s.Require().NoError(err)

	client := http.DefaultClient

	resp, err := client.Do(req)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()

	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *ServerTestsuite) Test_LogErrorIfUnableToBind() {
	address := "bad_address"
	s.server = *NewServerUnderTest("./app", "-a", address)
	err := s.server.Start("unable to listen")
	s.Require().NoError(err)
}

func (s *ServerTestsuite) Test_Exit1() {
	cmd := exec.Command("./app", "--unknwon")
	_ = cmd.Start()
	_ = cmd.Wait()
	s.Require().Equal(1, cmd.ProcessState.ExitCode())
}

func (s *ServerTestsuite) Test_ShutdownIn5Sec() {
	proxy := writer.NewProxy()
	finder := writer.NewFinder("Listening")
	proxy.SetWriter(finder)
	address := "127.0.0.1:8989"
	cmd := exec.Command("./app", "-a", address)
	cmd.Stderr = proxy
	cmd.Stdout = proxy
	_ = cmd.Start()
	defer func() {
		if !cmd.ProcessState.Exited() {
			_ = cmd.Process.Signal(syscall.SIGKILL)
		}
	}()

	err := finder.Wait(time.Second)
	s.Require().NoError(err)

	finder2 := writer.NewFinder("error while shutdown")
	proxy.SetWriter(finder2)
	ipAddr, _ := net.ResolveTCPAddr("tcp", address)
	c, err := net.DialTCP("tcp", nil, ipAddr)
	defer func() { _ = c.Close() }()
	s.Require().NoError(err)
	_ = cmd.Process.Signal(syscall.SIGINT)
	err = finder2.Wait(10 * time.Second)
	s.Require().NoError(err)
	_ = cmd.Wait()
}
