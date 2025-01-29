package db

import (
	"github.com/data_base_messages_saver/configs"
	m "github.com/maneulf/messages_models/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToMariaDB() (*gorm.DB, error) {
	if DB == nil {
		config := configs.ConfigFromEnv("")
		dsn := config.Service.DataBaseUrl
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		DB = db
		if err != nil {
			return DB, err
		}

	}

	return DB, nil
}

func AutoMigrate() {
	db, err := ConnectToMariaDB()

	if err != nil {
		panic(err.Error())
	}

	var csmlDatabaseMessage m.CsmlDataBaseMessageModelDB
	err = db.AutoMigrate(csmlDatabaseMessage)

	if err != nil {
		panic(err.Error())
	}

}
