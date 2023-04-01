package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ehabterra/rewards/internal/api"
	"github.com/ehabterra/rewards/internal/models"
	"github.com/ehabterra/rewards/internal/pb"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	var (
		cfg api.Config
		lis net.Listener
		db  *gorm.DB
		dsn string
		err error
	)

	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	lis, err = net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("server listening..")

	//dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DB)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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

	srv := api.NewGRPCRewardsService(db, &cfg)

	s := grpc.NewServer()
	pb.RegisterRewardsServiceServer(s, srv)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
