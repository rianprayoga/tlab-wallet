package server

import (
	"fmt"
	"log"
	"net/http"
	"tlab-wallet/internal/handler"
	"tlab-wallet/internal/repository"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
)

type httpServer struct {
	port        string
	httpHandler handler.HttpHandler
}

func NewHttpServer(
	port string,
	db repository.Repo,
	logger *logrus.Logger,
	jwt *jwtauth.JWTAuth) *httpServer {

	h := handler.HttpHandler{
		Db:  db,
		Log: logger,
		Jwt: jwt,
	}

	return &httpServer{
		port:        port,
		httpHandler: h,
	}
}

func (s *httpServer) Run() error {

	err := http.ListenAndServe(fmt.Sprintf(":%s", s.port), s.httpHandler.Routes())

	if err != nil {
		log.Fatal(err)
	}

	return nil

}
