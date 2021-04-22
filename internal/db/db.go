package db

import (
	"fmt"
	"log"

	"github.com/Confialink/wallet-currencies/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// CreateConnection creates connection with database
func CreateConnection() *gorm.DB {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		config.DbConfig.User, config.DbConfig.Password, config.DbConfig.Host, config.DbConfig.Port, config.DbConfig.Schema)
	db, err := gorm.Open(
		"mysql",
		connectionString,
	)

	if err != nil {
		log.Fatalf("Could not connect to DB: %v\n", err)
		return nil
	}

	if config.DbConfig.IsDebugMode {
		return db.Debug()
	}

	return db
}
