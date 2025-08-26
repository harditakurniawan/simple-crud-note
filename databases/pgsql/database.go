package databases

import (
	"log"
	"simple-crud-notes/configs"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseInit(appConfig configs.AppConfig) {
	var err error

	log.Println("Connecting to the database : " + appConfig.DATABASE_URL)

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  appConfig.DATABASE_URL,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Panicln("failed to connect database: " + err.Error())
	}

	dbConfig, err := DB.DB()

	if err != nil {
		log.Panicln("Failed to get sql.DB: " + err.Error())
	}

	dbConfig.SetMaxIdleConns(10)
	dbConfig.SetMaxOpenConns(100)
	dbConfig.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to the database")
}
