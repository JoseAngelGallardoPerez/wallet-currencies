package config

import (
	"github.com/Confialink/wallet-pkg-env_config"
	"github.com/inconshreveable/log15"
)

func Validate(logger log15.Logger) {
	validator := env_config.NewValidator(logger)
	validator.ValidateCors(CorsConfig, logger)
	validator.ValidateDb(DbConfig, logger)
	validator.CriticalIfEmpty(GeneralConfig.Port, "VELMIE_WALLET_CURRENCIES_PORT", logger)
	validator.CriticalIfEmpty(ProtoBufConfig.Port, "VELMIE_WALLET_CURRENCIES_PROTOBUF_PORT", logger)
}
