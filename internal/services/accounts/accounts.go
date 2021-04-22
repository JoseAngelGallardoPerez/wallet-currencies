package accounts

import (
	"context"

	pb "github.com/Confialink/wallet-accounts/rpc/accounts"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-currencies/internal/config/connections"
)

type Service struct {
	logger log15.Logger
}

func NewService(logger log15.Logger) *Service {
	return &Service{logger: logger.New("Service", "Accounts")}
}

func (s *Service) CanDisableCurrency(code string) bool {
	connection, err := connections.GetAccountsClient()
	if err != nil {
		s.logger.Error("failed to connect to accounts service", "error", err)
		return false
	}

	request := pb.DisableCurrencyReq{Code: code}
	resp, err := connection.CanDisableCurrency(context.Background(), &request)
	if err != nil {
		s.logger.Error("failed to check currency", "error", err)
		return false
	}

	return resp.Can
}
