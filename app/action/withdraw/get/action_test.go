package get

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
	"github.com/rusinov-artem/gophermart/test/utils/logger"
)

type GetWithdrawalActionTestSuite struct {
	suite.Suite
	action  *Action
	logger  *zap.Logger
	logs    *bytes.Buffer
	storage *storage
}

func Test_GetWithdrawalsAction(t *testing.T) {
	suite.Run(t, &GetWithdrawalActionTestSuite{})
}

func (s *GetWithdrawalActionTestSuite) SetupTest() {
	s.logger, s.logs = logger.SpyLogger()
	s.storage = &storage{}
	s.action = New(s.logger, s.storage)
}

func (s *GetWithdrawalActionTestSuite) Test_UnableToGetWithdrawals() {
	s.storage.err = fmt.Errorf("storage error")

	login := "login"

	_, internalErr := s.action.GetWithdrawals(login)

	s.NotNil(internalErr)
	s.Equal(appError.ServiceUnavailable, internalErr.Code)
	s.Equal(login, s.storage.login)
}

func (s *GetWithdrawalActionTestSuite) Test_NoWithdrawals() {
	login := "login"

	_, internalErr := s.action.GetWithdrawals(login)

	s.NotNil(internalErr)
	s.Equal(appError.NoWithdrawals, internalErr.Code)
	s.Equal(login, s.storage.login)
}

func (s *GetWithdrawalActionTestSuite) Test_Success() {
	login := "login"

	s.storage.withdrawals = []dto.Withdrawal{
		{
			OrderNr:     "OrderNR",
			Sum:         500,
			ProcessedAt: time.Time{},
		},
	}

	withdrawals, internalErr := s.action.GetWithdrawals(login)

	s.Nil(internalErr)
	s.Equal(login, s.storage.login)
	s.Equal(s.storage.withdrawals, withdrawals)
}

type storage struct {
	login       string
	withdrawals []dto.Withdrawal
	err         error
}

func (s *storage) GetWithdrawals(login string) ([]dto.Withdrawal, error) {
	s.login = login
	return s.withdrawals, s.err
}
