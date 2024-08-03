package bintest

import (
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/rusinov-artem/gophermart/test/utils/writer"
	"github.com/stretchr/testify/suite"
)

type ServerTestsuite struct {
	suite.Suite
	cmd *exec.Cmd
}

func Test_Server(t *testing.T) {
	Run(t, &ServerTestsuite{})
}

func (s *ServerTestsuite) Test_CanStartServer() {
	s.cmd = exec.Command("./app")

	p := writer.NewProxy()
	finder := writer.NewFinder("Hello World")
	p.SetWriter(finder)

	s.cmd.Stdout = p
	s.cmd.Stderr = p

	err := s.cmd.Start()
	s.Require().NoError(err)

	err = finder.Wait(time.Second)
	s.Require().NoError(err)

	finder2 := writer.NewFinder("Got signal: interrupt")
	p.SetWriter(finder2)

	err = s.cmd.Process.Signal(syscall.SIGINT)
	s.Require().NoError(err)

	err = s.cmd.Wait()
	s.Require().ErrorContains(err, "signal: interrupt")

	err = finder.Wait(time.Second)
	s.Require().NoError(err)
}
