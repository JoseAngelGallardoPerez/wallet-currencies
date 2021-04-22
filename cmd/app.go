package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"golang.org/x/sync/errgroup"

	"github.com/Confialink/wallet-currencies/internal/config"
	"github.com/Confialink/wallet-currencies/internal/pb_server"
	"github.com/Confialink/wallet-currencies/internal/workers"
)

type runPublicAPIParams struct {
	dig.In

	GeneralConfig config.General
	Router        *gin.Engine `name:"APIRouter"`
}

type App struct {
	container *dig.Container
}

func NewApp(container *dig.Container) *App {
	return &App{container: container}
}

func (a *App) Run() {
	var errg errgroup.Group
	err := a.container.Invoke(func(publicParams runPublicAPIParams, rpcServer *pb_server.PbServer, scheduler *workers.Jobs) {
		if err := scheduler.Start(); err != nil {
			log.Fatal(err)
		}

		errg.Go(rpcServer.Start)

		errg.Go(func() error {
			return a.runPublicAPI(publicParams)
		})
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := errg.Wait(); err != nil {
		log.Fatal(err)
	}
}

func (a *App) runPublicAPI(params runPublicAPIParams) error {
	return params.Router.Run(":" + params.GeneralConfig.Port)
}
