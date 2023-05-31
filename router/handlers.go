package router

import (
	"chatbot/config"
	"chatbot/pkg/whatsapp"
	"net/http"

	"github.com/gorilla/mux"
)

func AddHandlers(router *mux.Router, externalClient *http.Client) error {

	facebookConf := config.NewFacebookConfig()
	whatsappClient := whatsapp.NewClientHandler(externalClient)
	router.Handle(whatsapp.WebhooksEndpoint, whatsapp.NewCallbackHandler(facebookConf, whatsappClient)).Methods(http.MethodPost)
	router.Handle(whatsapp.WebhooksEndpoint, whatsapp.NewVerificationHandler(facebookConf, whatsappClient)).Methods(http.MethodGet)
	return nil
}
