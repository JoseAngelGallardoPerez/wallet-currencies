// Package router contains actions for routing
package router

import (
	"net/http"

	"github.com/Confialink/wallet-currencies/internal/authentication"
	"github.com/Confialink/wallet-currencies/internal/middleware"
	"github.com/Confialink/wallet-currencies/internal/version"
	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-currencies/internal/config"
	"github.com/Confialink/wallet-currencies/internal/controllers"
)

// Runs router
func APIRouter(
	currenciesController *controllers.Currencies,
	ratesController *controllers.Rates,
	settingsController *controllers.Settings,
	logger log15.Logger,
) *gin.Engine {
	r := gin.New()

	r.GET("/currencies/health-check", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/currencies/build", func(c *gin.Context) {
		c.JSON(http.StatusOK, version.BuildInfo)
	})

	trackableGroup := r.Group("")
	errorsMiddleware := errors.ErrorHandler(logger.New("Middleware", "Errors"))
	trackableGroup.Use(gin.Recovery(), gin.Logger(), corsMiddleware(), errorsMiddleware)

	trackableGroup.OPTIONS("/*cors")

	mwAdminOnly := middleware.AdminOnly

	mainGroup := trackableGroup.Group("/currencies")

	authMiddleware := authentication.Middleware(logger.New("Middleware", "Auth"))
	privateGroup := mainGroup.Group("/private", authMiddleware)
	publicGroup := mainGroup.Group("/public")
	v1Group := privateGroup.Group("/v1")
	publicV1Group := publicGroup.Group("/v1")
	adminGroup := v1Group.Group("/admin", mwAdminOnly)

	v1Group.GET("/currencies", currenciesController.Index)
	v1Group.GET("/currencies/:id", currenciesController.Show)

	publicV1Group.GET("/currency-logo/:code", currenciesController.GetLogo)

	v1Group.GET("/rates/main", ratesController.Index)
	v1Group.GET("/rates/pair", ratesController.GetForCurrencies)

	adminGroup.POST("/currencies", currenciesController.AdminCreate)
	update(adminGroup, "/currencies", currenciesController.AdminUpdateList)
	adminGroup.GET("/currencies", currenciesController.AdminIndex)
	adminGroup.GET("/currencies/:id", currenciesController.AdminShow)

	adminGroup.GET("/settings/main", settingsController.ShowMain)
	update(adminGroup, "/settings/main", settingsController.UpdateMain)
	update(adminGroup, "/rates", ratesController.UpdateList)

	return r
}

// Adds PUT and PATCH methods for one handler
func update(group *gin.RouterGroup, relativePath string, handler gin.HandlerFunc) {
	group.PUT(relativePath, handler)
	group.PATCH(relativePath, handler)
}

func corsMiddleware() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = config.CorsConfig.Methods

	for _, origin := range config.CorsConfig.Origins {
		if origin == "*" {
			corsConfig.AllowAllOrigins = true
		}
	}
	if !corsConfig.AllowAllOrigins {
		corsConfig.AllowOrigins = config.CorsConfig.Origins
	}
	corsConfig.AllowHeaders = config.CorsConfig.Headers

	return cors.New(corsConfig)
}
