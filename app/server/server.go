package server

import (
	"chatbot/config"
	"chatbot/router"
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"chatbot/pkg/db"
	chatbotlog "chatbot/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Run() error {

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := config.Reload(nil)
	if err != nil {
		log.Error().Err(err).Msg("error processing config data")
		return err
	}

	chatbotlog.InitLogging()

	database, dbHandler, err := db.ConnectDB(config.NewDatabaseConfig())
	if err != nil {
		log.Error().Err(err).Msg("error connecting database")
		return err
	}

	models := db.GetAllModels()
	err = database.AutoMigrate(models...)
	if err != nil {
		log.Error().Err(err).Msg("error run migration")
		return err
	}

	externalClientTLSConfig := &tls.Config{
		Renegotiation: tls.RenegotiateOnceAsClient,
	}

	if config.GetRuntimeEnvironment() == "local" {
		externalClientTLSConfig.InsecureSkipVerify = true
	}

	externalClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: externalClientTLSConfig,
		},
	}

	// muxRouter := mux.NewRouter()
	ginRouter := gin.Default()
	if err := router.AddHandlers(ginRouter, externalClient, dbHandler); err != nil {
		log.Error().Err(err).Msg("unable to setup server handler. exiting")
		return err
	}

	var srv *http.Server
	if srv, err = createServer(ginRouter); err != nil {
		log.Error().Err(err).Msg("unable to create server. exiting.")
		return err
	}

	log.Info().Str("address", srv.Addr).Msg("starting server")
	if serverErr := srv.ListenAndServe(); serverErr != nil {
		log.Error().Err(serverErr).Msg("error in ListenAndServe")
		return serverErr
	}

	return nil
}

func createServer(router *gin.Engine) (*http.Server, error) {
	serverPort := ":" + config.GetListenPort()

	srv := &http.Server{
		Addr:         serverPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	return srv, nil
}
