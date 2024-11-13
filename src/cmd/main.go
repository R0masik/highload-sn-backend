package main

import (
	"errors"
	"net/http"
	"os"

	"highload-sn-backend/config"
	"highload-sn-backend/db/postgres"
	"highload-sn-backend/internal/log"
	httpTransport "highload-sn-backend/transport/http"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Logger().Error(err)
		os.Exit(1)
	}

	err = postgres.InitClient()
	if err != nil {
		log.Logger().Error(err)
		os.Exit(1)
	}

	server, port := httpTransport.NewServer()
	log.Logger().Infof("listening at :%s", port)
	err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Logger().Info("server closed")
		return
	}

	if err != nil {
		log.Logger().Errorf("error starting server: %v\n", err)
		os.Exit(1)
	}
}
