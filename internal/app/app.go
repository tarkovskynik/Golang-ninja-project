package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tarkovskynik/Golang-ninja-project/internal/config"
	"github.com/tarkovskynik/Golang-ninja-project/internal/repository/psql"

	s3repo "github.com/tarkovskynik/Golang-ninja-project/internal/repository/s3"
	"github.com/tarkovskynik/Golang-ninja-project/internal/service"
	"github.com/tarkovskynik/Golang-ninja-project/internal/transport/rest"
	"github.com/tarkovskynik/Golang-ninja-project/pkg/database"
	"github.com/tarkovskynik/Golang-ninja-project/pkg/hash"
	"github.com/tarkovskynik/Golang-ninja-project/pkg/logger"
	"github.com/tarkovskynik/Golang-ninja-project/pkg/s3"
	"github.com/tarkovskynik/Golang-ninja-project/pkg/server/http"

	_ "github.com/lib/pq"
)

func init() {
	logger.InitLogParams()
}

func Run(configDir string) {
	cfg, err := config.Init(configDir)
	if err != nil {
		logger.Fatalf("Config initialization error %s", err)
	}

	db, err := database.NewPostgresConnection(cfg.Postgres.Host, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Name, cfg.Postgres.SSLMode, cfg.Postgres.Port)
	if err != nil {
		logger.Fatalf("Postgres connection error %s", err)
	}

	usersRepo := psql.NewUsers(db)
	tokensRepo := psql.NewTokens(db)
	hasher := hash.NewSHA1Hasher(cfg.Auth.Salt)
	usersService := service.NewUsers(usersRepo, tokensRepo, hasher, []byte(cfg.Auth.Secret), cfg.Auth.AccessTokenTTL, cfg.Auth.RefreshTokenTTL)

	cfg.File.MaxUploadSize = cfg.File.MaxUploadSize << 20 // 10 megabytes = 10 << 20
	cfg.File.CheckTypes = make(map[string]interface{})    // "image/jpeg": nil, "image/png": nil, ...
	for _, t := range cfg.File.Types {
		cfg.File.CheckTypes[t] = nil
	}

	filesRepo := psql.NewFiles(db)
	s3, err := s3.NewFileStorage(context.Background(), cfg.File.Storage)
	if err != nil {
		logger.Fatalf("error occurred while running s3 server: %s", err.Error())
	}
	s3FilesStorage := s3repo.NewS3FilesStorage(s3)
	filesService := service.NewServiceFiles(filesRepo, s3FilesStorage)

	handler := rest.NewHandler(usersService, filesService, cfg.File.MaxUploadSize, cfg.File.CheckTypes)
	server := http.NewServer()

	go func() {
		if err := server.Run(cfg.Server.Host, cfg.Server.Port, handler.InitRoutes()); err != nil {
			logger.Fatalf("error occurred while running rest server: %s", err.Error())
		}
	}()

	logger.Info("Http Server for fin manager service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("Http Server for fin manager service stopped")

	if err := db.Close(); err != nil {
		logger.Errorf("Error occurred on postgres connection close: %s", err.Error())
	}

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Errorf("Error occurred on http server for fin manager service shutting down: %s", err.Error())
	}
}
