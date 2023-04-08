package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ehabterra/rewards/internal/models"
	"github.com/ehabterra/rewards/internal/types"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	var (
		cfg types.Config
		db  *gorm.DB
		dsn string
		err error
	)

	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("migrating db..")

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DB)

	for tries := 0; tries < 5; tries++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Activity{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Create test users
	err = db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&models.User{
		ID:     1,
		Points: 0,
	}).Error
	if err != nil {
		log.Println(err)
	}

	err = db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&models.User{
		ID:     2,
		Points: 0,
	}).Error
	if err != nil {
		log.Println(err)
	}
}
