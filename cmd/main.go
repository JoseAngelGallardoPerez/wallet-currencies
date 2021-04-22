package main

import (
	"flag"

	"github.com/inconshreveable/log15"
	"github.com/jasonlvhit/gocron"
	"github.com/shopspring/decimal"
	"go.uber.org/dig"

	"github.com/Confialink/wallet-currencies/internal/config"
	"github.com/Confialink/wallet-currencies/internal/config/validator"
	"github.com/Confialink/wallet-currencies/internal/controllers"
	"github.com/Confialink/wallet-currencies/internal/db"
	"github.com/Confialink/wallet-currencies/internal/pb_server"
	"github.com/Confialink/wallet-currencies/internal/repositories"
	"github.com/Confialink/wallet-currencies/internal/router"
	"github.com/Confialink/wallet-currencies/internal/serializers"
	"github.com/Confialink/wallet-currencies/internal/services"
	"github.com/Confialink/wallet-currencies/internal/tasks"
	"github.com/Confialink/wallet-currencies/internal/workers"
)

func main() {
	config.Setup()

	decimal.DivisionPrecision = 32

	c := initContainer()
	if err := c.Invoke(config.Validate); err != nil {
		panic("cannot validate the app config: " + err.Error())
	}
	validator.Initialize(c)

	task := flag.String("t", "", "Pass task name to execute")
	flag.Parse()

	if *task == "" {
		app := NewApp(c)
		app.Run()
	} else {
		tasks.Execute(task, c)
	}
}

func initContainer() *dig.Container {
	container := dig.New()

	providers := []interface{}{
		//log15.Logger
		func() log15.Logger {
			return log15.New("service", "Currencies")
		},
		//*gorm.DB
		db.CreateConnection,
		//*gocron.Scheduler
		gocron.NewScheduler,
		pb_server.NewPbServer,
	}

	providers = append(providers, config.Providers()...)
	providers = append(providers, repositories.Providers()...)
	providers = append(providers, services.Providers()...)
	providers = append(providers, controllers.Providers()...)
	providers = append(providers, workers.Providers()...)
	providers = append(providers, serializers.Providers()...)

	for _, provider := range providers {
		err := container.Provide(provider)
		if err != nil {
			panic("unable to init container " + err.Error())
		}
	}

	namedProviders := map[string]interface{}{
		//ICurrencySerializer
		"currencySerializerAdmin": serializers.AdminCurrency,
		//ICurrencySerializer
		"currencySerializerDefault": serializers.Currency,
		//*gin.Engine
		"APIRouter": router.APIRouter,
	}

	for name, provider := range namedProviders {
		err := container.Provide(provider, dig.Name(name))
		if err != nil {
			panic("unable to init container " + err.Error())
		}
	}

	return container
}
