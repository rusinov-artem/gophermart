package bintest

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"github.com/rusinov-artem/gophermart/test/utils/writer"
)

type ServerUnderTest struct {
	cmd   *exec.Cmd
	proxy *writer.ProxyWriter
}

func NewServerUnderTest(name string, args ...string) *ServerUnderTest {
	server := &ServerUnderTest{}

	server.cmd = exec.Command("./app", args...)
	server.proxy = writer.NewProxy()

	server.cmd.Stdout = server.proxy
	server.cmd.Stderr = server.proxy

	return server
}

func (s *ServerUnderTest) Start() error {
	finder := writer.NewFinder("Hello World")
	s.proxy.SetWriter(finder)

	err := s.cmd.Start()
	if err != nil {
		return fmt.Errorf("unable to start server: %w", err)
	}

	err = finder.Wait(time.Second)
	if err != nil {
		return fmt.Errorf("unable to find starting long entry: %w", err)
	}

	return nil
}

func (s *ServerUnderTest) Stop() error {
	finder2 := writer.NewFinder("Got signal: interrupt")
	s.proxy.SetWriter(finder2)

	errSigInt := s.cmd.Process.Signal(syscall.SIGINT)
	if errSigInt != nil {
		err := fmt.Errorf("unable to send SIGINT: %w", errSigInt)
		errSigKill := s.cmd.Process.Signal(syscall.SIGKILL)
		if errSigKill != nil {
			err := fmt.Errorf("unable to send SIGKILL: %w: %w", errSigKill, err)
			return err
		}
		return err
	}

	finderErr := finder2.Wait(1 * time.Second)
	if finderErr != nil {
		err := fmt.Errorf("unable to find finish log entry: %w", finderErr)
		errSigKill := s.cmd.Process.Signal(syscall.SIGKILL)
		if errSigKill != nil {
			err := fmt.Errorf("unable to send SIGKILL: %w: %w", errSigKill, err)
			return err
		}
	}
	_ = s.cmd.Wait()

	return nil
}
