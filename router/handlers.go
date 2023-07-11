package router

import (
	"chatbot/config"
	"chatbot/pkg/db"
	"chatbot/pkg/whatsapp"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddHandlers(router *gin.Engine, externalClient *http.Client, databaseHandler db.IDatabaseAdapter) error {

	facebookConf := config.NewFacebookConfig()
	whatsappClient := whatsapp.NewClientHandler(externalClient)
	router.POST(whatsapp.WebhooksEndpoint, whatsapp.NewCallbackHandler(facebookConf, whatsappClient, databaseHandler))
	router.GET(whatsapp.WebhooksEndpoint, whatsapp.NewVerificationHandler(facebookConf, whatsappClient))
	return nil
}
