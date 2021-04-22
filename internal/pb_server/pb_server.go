package pb_server

import (
	"context"
	"fmt"
	"net/http"

	pb "github.com/Confialink/wallet-currencies/rpc/currencies"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-currencies/internal/config"
	"github.com/Confialink/wallet-currencies/internal/models"
	"github.com/Confialink/wallet-currencies/internal/services/currencies"
	"github.com/Confialink/wallet-currencies/internal/services/rates"
)

type PbServer struct {
	currenciesService *currencies.Service
	ratesService      *rates.Rates
	logger            log15.Logger
}

func NewPbServer(currenciesService *currencies.Service, ratesService *rates.Rates, logger log15.Logger) *PbServer {
	return &PbServer{currenciesService, ratesService, logger.New("context", "RPC Server")}
}

func (s *PbServer) Start() error {
	twirpHandler := pb.NewCurrencyFetcherServer(s, nil)
	mux := http.NewServeMux()
	mux.Handle(pb.CurrencyFetcherPathPrefix, twirpHandler)
	s.logger.Info("RPC server is starting", "port", config.ProtoBufConfig.Port)
	return http.ListenAndServe(fmt.Sprintf(":%s", config.ProtoBufConfig.Port), mux)
}

func (s *PbServer) GetCurrency(_ context.Context, req *pb.CurrencyReq) (*pb.CurrencyResp, error) {
	currenciesService := s.currenciesService
	if req.Code != "" {
		currency, err := currenciesService.FindByCode(req.Code)
		if err != nil {
			s.logger.Error("cannot obtain currency by code", "err", err)
			return &pb.CurrencyResp{}, err
		}
		return s.getCurrencyResp(currency), nil
	}
	if currency, err := currenciesService.FindById(req.Id); err != nil {
		s.logger.Error("cannot obtain currency by id", "err", err)
		return &pb.CurrencyResp{}, err
	} else {
		return s.getCurrencyResp(currency), nil
	}
}

func (s *PbServer) GetMain(_ context.Context, _ *pb.CurrencyReq) (*pb.CurrencyResp, error) {
	currency, err := s.currenciesService.Main()
	if err != nil {
		s.logger.Error("cannot obtain main currency", "err", err)
		return &pb.CurrencyResp{}, err
	}
	return s.getCurrencyResp(currency), nil
}

func (s *PbServer) GetCurrenciesRateByCodes(_ context.Context, req *pb.CurrenciesRateValueRequest) (*pb.CurrenciesRateResponse, error) {
	res, err := s.ratesService.CalculateRate(req.CurrencyCodeFrom, req.CurrencyCodeTo)
	if err != nil {
		s.logger.Error("cannot calculate rate", "err", err)
		return &pb.CurrenciesRateResponse{}, err
	}

	return &pb.CurrenciesRateResponse{Value: res.Rate.String(), ExchangeMargin: res.ExchangeMargin.String()}, nil
}

func (s *PbServer) GetCurrenciesRateValueByCodes(_ context.Context, req *pb.CurrenciesRateValueRequest) (*pb.CurrenciesRateValueResponse, error) {
	res, err := s.ratesService.CalculateRate(req.CurrencyCodeFrom, req.CurrencyCodeTo)
	if err != nil {
		s.logger.Error("cannot calculate rate", "err", err)
		return &pb.CurrenciesRateValueResponse{}, err
	}
	return &pb.CurrenciesRateValueResponse{Value: res.Rate.String()}, nil
}

func (s *PbServer) getCurrencyResp(currency *models.Currency) *pb.CurrencyResp {
	return &pb.CurrencyResp{Id: currency.ID,
		Code:          currency.Code,
		Active:        *currency.Active,
		DecimalPlaces: int32(currency.DecimalPlaces),
	}
}
