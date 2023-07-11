package whatsapp

import (
	"chatbot/config"
	"chatbot/pkg/db"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const customerKey = "phone_number"

func NewCallbackHandler(config *config.FacebookConf, client IWhatsappClient, dbHandler db.IDatabaseAdapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug().Msg("[WhatsappCallback] entering whatsapp callback...")
		var requestBody WebHookPayload

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
		if err != nil {
			log.Warn().Err(err).Msgf("[WhatsappCallback] failed to parse request")
			http.Error(c.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		for _, entry := range requestBody.Entry {
			for _, change := range entry.Changes {
				phoneNumber := getPhoneNumberFromChange(change)

				_, err := dbHandler.Any(&db.Customer{}, customerKey, phoneNumber)
				if err != nil && err == gorm.ErrRecordNotFound {
					isSuccess, err := client.SendMessageText(c.Request.Context(), config, requestBody)
					if !isSuccess && err != nil {
						log.Warn().Err(err).Msgf("[WhatsappCallback] failed to sent text")
						http.Error(c.Writer, err.Error(), http.StatusBadRequest)
						return
					}

					isSuccess, err = client.SendMessageAudio(c.Request.Context(), config, requestBody)
					if !isSuccess && err != nil {
						log.Warn().Err(err).Msgf("[WhatsappCallback] failed to sent audio")
						http.Error(c.Writer, err.Error(), http.StatusBadRequest)
						return
					}

					_, err = dbHandler.Insert(&db.Customer{
						PhoneNumber: phoneNumber,
					})
					if err != nil {
						log.Warn().Err(err).Msgf("[WhatsappCallback] failed to insert phonenumber")
						http.Error(c.Writer, err.Error(), http.StatusBadRequest)
						return
					}
				}
			}
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		response := []byte("Ok!")
		_, err = c.Writer.Write(response)
		c.Writer.WriteHeader(http.StatusOK)
		if err != nil {
			return
		}
	}
}

func getPhoneNumberFromChange(change Change) string {
	if len(change.Value.Contacts) > 0 {
		contact := change.Value.Contacts[0]
		return contact.WaID
	}
	return ""
}

// func validateRequest(w http.ResponseWriter, r *http.Request) request
