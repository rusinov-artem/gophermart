package bintest

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

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

	w := &Writer{bytes.NewBuffer([]byte{})}

	s.cmd.Stdout = w
	s.cmd.Stderr = w

	err := s.cmd.Start()
	s.Require().NoError(err)

	err = s.cmd.Wait()
	s.Require().NoError(err)

	s.Contains(w.Logs.String(), "Hello World")
}

type Writer struct {
	Logs *bytes.Buffer
}

func (w *Writer) Write(p []byte) (n int, err error) {
	fmt.Println("App:", string(p))
	return w.Logs.Write(p)
}
