package bintest

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

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
	address   string
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
	s.address = "127.0.0.1:8080"
	server := s.startServer(s.address, s.dsn)
	defer s.stopServer(server)

	resp := s.liveness()
	defer func() { _ = resp.Body.Close() }()

	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *ServerTestsuite) Test_CanRegister() {
	s.address = "127.0.0.1:8080"
	server := s.startServer(s.address, s.dsn)
	defer s.stopServer(server)

	authToken := ""
	orderNr := order.Number()

	s.Run("can register a user", func() {
		finder := writer.NewFinder("/api/user/register")
		server.proxy.SetWriter(finder)

		resp := s.register(`{"login": "user1", "password": "password1"}`)
		defer func() { _ = resp.Body.Close() }()

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.NotEmpty(resp.Header.Get("Authorization"))
		authToken = resp.Header.Get("Authorization")
		s.NoError(finder.Wait(time.Second))
	})

	s.Run("registered user can login", func() {
		finder := writer.NewFinder("/api/user/login")
		server.proxy.SetWriter(finder)

		resp := s.login(`{"login": "user1", "password": "password1"}`)
		defer func() { _ = resp.Body.Close() }()

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.NotEmpty(resp.Header.Get("Authorization"))
		s.Require().NoError(finder.Wait(time.Second))
		fmt.Println(resp.Body)
	})

	s.T().Run("user can add order number", func(t *testing.T) {
		finder := writer.NewFinder("/api/user/orders")
		server.proxy.SetWriter(finder)

		resp := s.addOrder(authToken, orderNr)
		defer func() { _ = resp.Body.Close() }()

		s.Require().Equal(http.StatusAccepted, resp.StatusCode)
		s.Require().NoError(finder.Wait(time.Second))
		fmt.Println(resp.Body)
	})

	s.Run("user can list orders", func() {
		finder := writer.NewFinder("/api/user/orders")
		server.proxy.SetWriter(finder)

		resp := s.listOrders(authToken)
		defer func() { _ = resp.Body.Close() }()

		b, _ := io.ReadAll(resp.Body)
		fmt.Println("BODY => ", string(b))
		s.Require().NotEmpty(b)

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal("application/json", resp.Header.Get("Content-Type"))
		s.Require().NoError(finder.Wait(time.Second))
	})

	s.Run("user can check balance", func() {
		finder := writer.NewFinder("/api/user/balance")
		server.proxy.SetWriter(finder)

		resp := s.getBalance(authToken)
		defer func() { _ = resp.Body.Close() }()

		b, _ := io.ReadAll(resp.Body)
		fmt.Println("BODY => ", string(b))
		s.Require().NotEmpty(b)

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal("application/json", resp.Header.Get("Content-Type"))
		s.Require().NoError(finder.Wait(time.Second))
	})

	s.Run("user can withdraw", func() {
		finder := writer.NewFinder("/api/user/balance/withdraw")
		server.proxy.SetWriter(finder)

		resp := s.withdraw(authToken, fmt.Sprintf(`{"number":"%s", "sum": 0 }`, order.Number()))
		defer func() { _ = resp.Body.Close() }()

		b, _ := io.ReadAll(resp.Body)
		fmt.Println("BODY => ", string(b))
		s.Require().Empty(b)

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal("application/json", resp.Header.Get("Content-Type"))
		s.Require().NoError(finder.Wait(time.Second))
	})

	s.Run("user can list withdraw operations", func() {
		finder := writer.NewFinder("/api/user/withdrawals")
		server.proxy.SetWriter(finder)

		resp := s.listWithdrawals(authToken)
		defer func() { _ = resp.Body.Close() }()

		b, _ := io.ReadAll(resp.Body)
		fmt.Println("BODY => ", string(b))
		s.Require().NotEmpty(b)

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal("application/json", resp.Header.Get("Content-Type"))
		s.Require().NoError(finder.Wait(time.Second))
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

func (s *ServerTestsuite) url(path string) string {
	return fmt.Sprintf("http://%s%s", s.address, path)
}

func (s *ServerTestsuite) liveness() *http.Response {
	req, err := http.NewRequest(http.MethodGet, s.url("/liveness"), nil)
	s.Require().NoError(err)
	client := http.DefaultClient
	resp, err := client.Do(req)
	s.Require().NoError(err)
	return resp
}

func (s *ServerTestsuite) register(json string) *http.Response {
	req, err := http.NewRequest(
		http.MethodPost,
		s.url("/api/user/register"),
		bytes.NewBufferString(json))
	s.Require().NoError(err)

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient

	resp, err := client.Do(req)
	s.Require().NoError(err)
	return resp
}

func (s *ServerTestsuite) login(json string) *http.Response {
	req, err := http.NewRequest(
		http.MethodPost,
		s.url("/api/user/login"),
		bytes.NewBufferString(json))
	s.Require().NoError(err)

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient

	resp, err := client.Do(req)
	s.Require().NoError(err)
	return resp
}

func (s *ServerTestsuite) addOrder(token, orderNr string) *http.Response {
	req, err := http.NewRequest(
		http.MethodPost,
		s.url("/api/user/orders"),
		bytes.NewBufferString(orderNr))
	s.Require().NoError(err)

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", token)

	client := http.DefaultClient
	resp, err := client.Do(req)
	s.Require().NoError(err)
	return resp
}

func (s *ServerTestsuite) listOrders(token string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, s.url("/api/user/orders"), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", token)

	client := http.DefaultClient
	resp, err := client.Do(req)
	s.NoError(err)
	return resp
}

func (s *ServerTestsuite) getBalance(token string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, s.url("/api/user/balance"), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", token)

	client := http.DefaultClient
	resp, err := client.Do(req)
	s.NoError(err)
	return resp
}

func (s *ServerTestsuite) withdraw(token, json string) *http.Response {
	req, err := http.NewRequest(
		http.MethodPost,
		s.url("/api/user/balance/withdraw"),
		bytes.NewBufferString(json),
	)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := http.DefaultClient
	resp, err := client.Do(req)
	s.NoError(err)
	return resp
}

func (s *ServerTestsuite) listWithdrawals(token string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, s.url("/api/user/withdrawals"), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := http.DefaultClient
	resp, err := client.Do(req)
	s.NoError(err)
	return resp
}
