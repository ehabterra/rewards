package main

import (
	"fmt"
	"net"
	"time"

	"github.com/ehabterra/rewards/internal/api"
	"github.com/ehabterra/rewards/internal/pb"
	"github.com/ehabterra/rewards/internal/types"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	var (
		cfg types.Config
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

	log.Debug("server listening..")

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

	srv := api.NewGRPCRewardsService(db, &cfg)

	s := grpc.NewServer()
	pb.RegisterRewardsServiceServer(s, srv)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
