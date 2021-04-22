package config

import (
	"os"

	"github.com/Confialink/wallet-pkg-env_config"
	"github.com/Confialink/wallet-pkg-env_mods"
)

type General struct {
	Env  string
	Port string
}

type ProtoBuf struct {
	Port string
}

var GeneralConfig General
var CorsConfig *env_config.Cors
var DbConfig *env_config.Db
var ProtoBufConfig ProtoBuf

func init() {
	GeneralConfig = NewGeneral()
	CorsConfig = NewCors()
	DbConfig = NewDb()
	ProtoBufConfig = NewProtoBuf()
}

func NewCors() *env_config.Cors {
	defaultConfigReader := env_config.NewReader("currencies")
	return defaultConfigReader.ReadCorsConfig()
}

func NewDb() *env_config.Db {
	defaultConfigReader := env_config.NewReader("currencies")
	return defaultConfigReader.ReadDbConfig()
}

func NewGeneral() General {
	env := env_config.Env("ENV", env_mods.Development)
	port := os.Getenv("VELMIE_WALLET_CURRENCIES_PORT")

	return General{env, port}
}

func NewProtoBuf() ProtoBuf {
	return ProtoBuf{Port: os.Getenv("VELMIE_WALLET_CURRENCIES_PROTOBUF_PORT")}
}
