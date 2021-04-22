package config

import (
	"github.com/Confialink/wallet-pkg-env_mods"
	"github.com/gin-gonic/gin"
)

func Setup() {
	ginMode := env_mods.GetMode(GeneralConfig.Env)
	gin.SetMode(ginMode)
}
