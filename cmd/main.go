package main

import (
	"auth_service/api"
	"auth_service/configs"
	"auth_service/grpc"
	"auth_service/pkg/logger"
	"auth_service/storage"
	"context"

	"fmt"
	"net"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	cfg := configs.Load()

	log := logger.NewLogger(cfg.ServiceName, cfg.LoggerLevel, cfg.LogPath)
	defer logger.Cleanup(log)

	storage, err := storage.New(context.Background(), cfg, &log)
	if err != nil {
		log.Panic("error while creating storage in main", logger.Error(err))
		return
	}
	defer storage.Close()

	go ServiceRun(&wg, log, cfg, &storage)
	go RouterRun(&wg, log, cfg, storage)

	wg.Wait()
}

func ServiceRun(wg *sync.WaitGroup, log logger.ILogger, cfg *configs.Config, storage *storage.IStorage) {
	defer wg.Done()
	listener, err := net.Listen("tcp",
		cfg.UserServiceGrpcHost+cfg.UserServiceGrpcPort,
	)
	if err != nil {
		log.Panic("error while creating listener for user service", logger.Error(err))
		return
	}
	defer listener.Close()

	server := grpc.SetUpServer(storage, log)

	fmt.Printf("User service is listening on port %s...\n",
		cfg.UserServiceGrpcHost+cfg.UserServiceGrpcPort)
	if err := server.Serve(listener); err != nil {
		log.Fatal("Error with listening user server", logger.Error(err))
	}
}

func RouterRun(wg *sync.WaitGroup, log logger.ILogger, cfg *configs.Config, storage storage.IStorage) {
	defer wg.Done()

	r := api.NewRouter(log, storage)
	r.Run(cfg.UserServiceHttpHost + cfg.UserServiceHttpPort)
}
