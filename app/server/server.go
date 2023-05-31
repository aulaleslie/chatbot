package server

import (
	"chatbot/config"
	"chatbot/router"
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	chatbotlog "chatbot/pkg/log"

	"github.com/gorilla/mux"
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

	muxRouter := mux.NewRouter()
	if err := router.AddHandlers(muxRouter, externalClient); err != nil {
		log.Error().Err(err).Msg("unable to setup server handler. exiting")
		return err
	}

	var srv *http.Server
	if srv, err = createServer(muxRouter); err != nil {
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

func createServer(h http.Handler) (*http.Server, error) {
	if h == nil {
		return nil, errors.New("missing server mux")
	}

	serverPort := ":" + config.GetListenPort()
	srv := &http.Server{
		Addr:         serverPort,
		Handler:      h,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return srv, nil
}
