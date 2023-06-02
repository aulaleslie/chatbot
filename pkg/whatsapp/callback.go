package whatsapp

import (
	"chatbot/config"
	"chatbot/pkg/db"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const customerKey = "phone_number"

func NewCallbackHandler(config *config.FacebookConf, client IWhatsappClient, dbHandler db.IDatabaseAdapter) http.Handler {
	return &CallbackHandler{
		facebookConfig: config,
		client:         client,
		dbHandler:      dbHandler,
	}
}

func (handler CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("[WhatsappCallback] entering whatsapp callback...")
	var requestBody WebHookPayload

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Warn().Err(err).Msgf("[WhatsappCallback] failed to parse request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, entry := range requestBody.Entry {
		for _, change := range entry.Changes {
			phoneNumber := getPhoneNumberFromChange(change)

			_, err := handler.dbHandler.Any(&db.Customer{}, customerKey, phoneNumber)
			if err != nil && err == gorm.ErrRecordNotFound {
				isSuccess, err := handler.client.SendMessageText(r.Context(), handler.facebookConfig, requestBody)
				if !isSuccess && err != nil {
					log.Warn().Err(err).Msgf("[WhatsappCallback] failed to sent text")
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				isSuccess, err = handler.client.SendMessageAudio(r.Context(), handler.facebookConfig, requestBody)
				if !isSuccess && err != nil {
					log.Warn().Err(err).Msgf("[WhatsappCallback] failed to sent audio")
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				_, err = handler.dbHandler.Insert(&db.Customer{
					PhoneNumber: phoneNumber,
				})
				if err != nil {
					log.Warn().Err(err).Msgf("[WhatsappCallback] failed to insert phonenumber")
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	response := []byte("Ok!")
	_, err = w.Write(response)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		return
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
