package main

import (
	"log"
	"os"
	"tlab-wallet/internal/config"
	"tlab-wallet/internal/repository"
	"tlab-wallet/internal/repository/pg"
	"tlab-wallet/internal/server"

	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type application struct {
	DSN      string
	HttpPort string
	Db       repository.Repo
	Log      *logrus.Logger
	Jwt      *jwtauth.JWTAuth
}

func main() {

	curDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	err = godotenv.Load(curDir + "/.env")
	if err != nil {
		log.Fatal("can't load .env file from current directory: " + curDir)
	}

	var app application
	app.DSN = os.Getenv("DB_SOURCE")
	app.HttpPort = os.Getenv("PORT")

	app.Log = config.NewLogger()
	app.Jwt = config.GenerateAuthToken()

	conn, err := app.connectDb()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	app.Db = &pg.PgRepo{
		DB: conn,
	}

	httpServer := server.NewHttpServer(
		app.HttpPort,
		app.Db,
		app.Log,
		app.Jwt,
	)
	err = httpServer.Run()
	if err != nil {
		panic(err)
	}
}
