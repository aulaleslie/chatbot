package whatsapp

import (
	"chatbot/config"
	"net/http"
)

func NewVerificationHandler(config *config.FacebookConf, client IWhatsappClient) http.Handler {
	return &VerificationHandler{
		facebookConfig: config,
		client:         client,
	}
}

func (handler VerificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	challenge := r.URL.Query().Get("hub.challenge")
	verifyToken := r.URL.Query().Get("hub.verify_token")

	// TODO: cek parameter verify token sama dengan yang ada di config
	if verifyToken != handler.facebookConfig.VerifyToken {
		http.Error(w, "invalid verify token", http.StatusBadRequest)
		return
	}

	// TODO: kirim balik value hub.challenge sebagai response
	w.Header().Set("Content-Type", "application/json")
	response := []byte(challenge)
	_, err := w.Write(response)
	w.WriteHeader(http.StatusCreated)
	if err != nil {
		return
	}

}
