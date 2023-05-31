package whatsapp

import (
	"chatbot/config"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func NewCallbackHandler(config *config.FacebookConf, client IWhatsappClient) http.Handler {
	return &CallbackHandler{
		facebookConfig: config,
		client:         client,
	}
}

func (handler CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("[WhatsappCallback] entering whatsapp callback...")
	var p CallbackPayload

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	handler.client.SendMessageText(r.Context(), handler.facebookConfig, p.DestinationNumber, p.Message)
	handler.client.SendMessageAudio(r.Context(), handler.facebookConfig, p.DestinationNumber, handler.facebookConfig.AudioId)

	w.Header().Set("Content-Type", "application/json")
	response := []byte("Ok!")
	_, err = w.Write(response)
	w.WriteHeader(http.StatusCreated)
	if err != nil {
		return
	}
}

// func validateRequest(w http.ResponseWriter, r *http.Request) request
