package bintest

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/action/order"
	"github.com/rusinov-artem/gophermart/test"
	"github.com/rusinov-artem/gophermart/test/utils"
	"github.com/rusinov-artem/gophermart/test/utils/writer"
)

type ServerTestsuite struct {
	suite.Suite
	dsn       string
	env       *utils.EnvManager
	cleanEnvs func()
}

func Test_Server(t *testing.T) {
	Run(t, &ServerTestsuite{})
}

func (s *ServerTestsuite) SetupSuite() {
	s.dsn = test.CreateTestDB("test_can_start_server")
}

func (s *ServerTestsuite) SetupTest() {
	name := s.T().Name()
	name = strings.ReplaceAll(name, "/", "-")
	SetupCoverDir(fmt.Sprintf("/app/test/bintest/coverdir/%s", name))
	s.env, s.cleanEnvs = utils.NewEnvManager()
}

func (s *ServerTestsuite) TearDownTest() {
	s.cleanEnvs()
}

func (s *ServerTestsuite) startServer(address, dsn string) *ServerUnderTest {
	server := NewServerUnderTest("./app",
		"-a", address,
		"-d", dsn,
	)
	err := server.Start("Listening")
	s.Require().NoError(err)
	return server
}

func (s *ServerTestsuite) stopServer(server *ServerUnderTest) {
	err := server.Stop()
	s.Require().NoError(err)
}

func (s *ServerTestsuite) Test_CanStartServer() {
	address := "127.0.0.1:8080"
	server := s.startServer(address, s.dsn)
	defer s.stopServer(server)
	url := fmt.Sprintf("http://%s/liveness", address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	s.Require().NoError(err)

	client := http.DefaultClient

	resp, err := client.Do(req)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()

	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *ServerTestsuite) Test_CanRegister() {
	address := "127.0.0.1:8080"
	server := s.startServer(address, s.dsn)
	defer s.stopServer(server)

	authToken := ""
	orderNr := order.OrderNr()

	s.T().Run("can register a user", func(t *testing.T) {
		finder := writer.NewFinder("/api/user/register")
		server.proxy.SetWriter(finder)

		url := fmt.Sprintf("http://%s/api/user/register", address)
		req, err := http.NewRequest(
			http.MethodPost,
			url,
			bytes.NewBufferString(`{"login": "user1", "password": "password1"}`))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/json")

		client := http.DefaultClient

		resp, err := client.Do(req)
		s.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Authorization"))
		authToken = resp.Header.Get("Authorization")
		s.NoError(finder.Wait(time.Second))
	})

	s.T().Run("registered user can login", func(t *testing.T) {
		finder := writer.NewFinder("/api/user/login")
		server.proxy.SetWriter(finder)

		url := fmt.Sprintf("http://%s/api/user/login", address)
		req, err := http.NewRequest(
			http.MethodPost,
			url,
			bytes.NewBufferString(`{"login": "user1", "password": "password1"}`))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := http.DefaultClient

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Authorization"))
		require.NoError(t, finder.Wait(time.Second))
		fmt.Println(resp.Body)
	})

	s.T().Run("user can add order number", func(t *testing.T) {
		finder := writer.NewFinder("/api/user/orders")
		server.proxy.SetWriter(finder)

		url := fmt.Sprintf("http://%s/api/user/orders", address)
		req, err := http.NewRequest(
			http.MethodPost,
			url,
			bytes.NewBufferString(orderNr))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Authorization", authToken)

		client := http.DefaultClient
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusAccepted, resp.StatusCode)
		require.NoError(t, finder.Wait(time.Second))
		fmt.Println(resp.Body)
	})

	s.T().Run("user can list orders", func(t *testing.T) {
		finder := writer.NewFinder("/api/user/orders")
		server.proxy.SetWriter(finder)

		url := fmt.Sprintf("http://%s/api/user/orders", address)
		req, err := http.NewRequest(
			http.MethodGet,
			url,
			nil,
		)

		require.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Authorization", authToken)

		client := http.DefaultClient
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		fmt.Println("BODY => ", resp.Body)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		require.NoError(t, finder.Wait(time.Second))
	})

	s.T().Run("user can check balance", func(t *testing.T) {
		finder := writer.NewFinder("/api/user/balance")
		server.proxy.SetWriter(finder)

		url := fmt.Sprintf("http://%s/api/user/balance", address)
		req, err := http.NewRequest(
			http.MethodGet,
			url,
			nil,
		)

		require.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Authorization", authToken)

		client := http.DefaultClient
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		fmt.Println("BODY => ", resp.Body)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		require.NoError(t, finder.Wait(time.Second))
	})

	s.T().Run("user can withdraw", func(t *testing.T) {
		t.Skip("broken")
		finder := writer.NewFinder("/api/user/balance/withdraw")
		server.proxy.SetWriter(finder)

		url := fmt.Sprintf("http://%s/api/balance/withdraw", address)
		req, err := http.NewRequest(
			http.MethodPost,
			url,
			nil,
		)

		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authToken)

		client := http.DefaultClient
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		fmt.Println("BODY => ", resp.Body)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		require.NoError(t, finder.Wait(time.Second))
	})
}

func (s *ServerTestsuite) Test_LogErrorIfUnableToBind() {
	address := "bad_address"
	server := *NewServerUnderTest("./app", "-a", address)
	err := server.Start("unable to listen")
	s.Require().NoError(err)
}

func (s *ServerTestsuite) Test_Exit1() {
	cmd := exec.Command("./app", "--unknwon")
	_ = cmd.Start()
	_ = cmd.Wait()
	s.Require().Equal(1, cmd.ProcessState.ExitCode())
}

func (s *ServerTestsuite) Test_MigrationError() {
	cmd := exec.Command("./app", "migrate",
		"-d", "badurl",
	)
	_ = cmd.Start()
	_ = cmd.Wait()
	s.Require().Equal(1, cmd.ProcessState.ExitCode())
}

func (s *ServerTestsuite) Test_MigrateWithoutMigrationDir() {
	s.env.Set("MIGRATION_DIR", "")
	cmd := exec.Command("./app", "migrate")
	proxy := writer.NewProxy()
	finder := writer.NewFinder("goose: unable to run migration")
	proxy.SetWriter(finder)
	cmd.Stdout = proxy
	cmd.Stderr = proxy
	_ = cmd.Start()
	_ = cmd.Wait()
	s.NoError(finder.Wait(time.Second))
}

func (s *ServerTestsuite) Test_ShutdownIn5Sec() {
	proxy := writer.NewProxy()
	finder := writer.NewFinder("Listening")
	proxy.SetWriter(finder)
	address := "127.0.0.1:8989"
	cmd := exec.Command("./app",
		"-a", address,
		"-d", s.dsn,
	)

	cmd.Stderr = proxy
	cmd.Stdout = proxy
	_ = cmd.Start()
	defer func() {
		if cmd.ProcessState != nil {
			if !cmd.ProcessState.Exited() {
				_ = cmd.Process.Signal(syscall.SIGKILL)
			}
		}
	}()

	err := finder.Wait(time.Second)
	s.Require().NoError(err)

	finder2 := writer.NewFinder("error while shutdown")
	proxy.SetWriter(finder2)
	ipAddr, _ := net.ResolveTCPAddr("tcp", address)
	c, err := net.DialTCP("tcp", nil, ipAddr)
	s.Require().NoError(err)
	_, err = c.Write([]byte("GET / HTTP/1.1"))
	time.Sleep(500 * time.Millisecond)
	s.Require().NoError(err)
	defer func() { _ = c.Close() }()

	_ = cmd.Process.Signal(syscall.SIGINT)
	err = finder2.Wait(6 * time.Second)
	s.Require().NoError(err)
	_ = cmd.Wait()
}
