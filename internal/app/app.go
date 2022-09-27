package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/tarkovskynik/Golang-ninja-project/internal/config"
	"github.com/tarkovskynik/Golang-ninja-project/internal/repository/psql"
	"github.com/tarkovskynik/Golang-ninja-project/internal/service"
	"github.com/tarkovskynik/Golang-ninja-project/internal/transport/rest"
	"github.com/tarkovskynik/Golang-ninja-project/pkg/database"
	"github.com/tarkovskynik/Golang-ninja-project/pkg/hash"
	"github.com/tarkovskynik/Golang-ninja-project/pkg/server/http"

	_ "github.com/lib/pq"
)

// @title File Manager App API
// @version 1.0
// @description API Server for File Manager Application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func Run(configDir string) {
	cfg, err := config.Init(configDir)
	if err != nil {
		logrus.Fatalf("Config initialization error %s", err)
	}

	db, err := database.NewPostgresConnection(cfg.Postgres.Host, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Name, cfg.Postgres.SSLMode, cfg.Postgres.Port)
	if err != nil {
		logrus.Fatalf("Postgres connection error %s", err)
	}

	usersRepo := psql.NewUsers(db)
	tokensRepo := psql.NewTokens(db)
	hasher := hash.NewSHA1Hasher(cfg.Auth.Salt)
	usersService := service.NewUsers(usersRepo, tokensRepo, hasher, []byte(cfg.Auth.Secret), cfg.Auth.AccessTokenTTL, cfg.Auth.RefreshTokenTTL)

	handler := rest.NewHandler(usersService)
	server := http.NewServer()

	go func() {
		if err := server.Run(cfg.Server.Host, cfg.Server.Port, handler.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running rest server: %s", err.Error())
		}
	}()

	logrus.Info("Http Server for fin manager service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Info("Http Server for fin manager service stopped")

	if err := db.Close(); err != nil {
		logrus.Errorf("Error occured on postgres connection close: %s", err.Error())
	}

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error occured on http server for fin manager service shutting down: %s", err.Error())
	}
}
