package whatsapp

import (
	"chatbot/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewVerificationHandler(config *config.FacebookConf, client IWhatsappClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		challenge := c.Query("hub.challenge")
		verifyToken := c.Query("hub.verify_token")

		// TODO: cek parameter verify token sama dengan yang ada di config
		if verifyToken != config.VerifyToken {
			http.Error(c.Writer, "invalid verify token", http.StatusBadRequest)
			return
		}

		// TODO: kirim balik value hub.challenge sebagai response
		c.Writer.Header().Set("Content-Type", "application/json")
		response := []byte(challenge)
		_, err := c.Writer.Write(response)
		c.Writer.WriteHeader(http.StatusCreated)
		if err != nil {
			return
		}

	}

}
