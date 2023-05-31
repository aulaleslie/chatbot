package whatsapp

import (
	"bytes"
	"chatbot/config"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/context/ctxhttp"
)

func NewClientHandler(externalClient *http.Client) IWhatsappClient {
	return &ClientHandler{
		externalClient: externalClient,
	}
}

func (handler *ClientHandler) SendMessageText(ctx context.Context, facebookConfig *config.FacebookConf, destinationNumber, message string) (bool, error) {
	log.Debug().Msg("[WhatsappCallback] entering whatsapp client text...")

	messagePayload := GetMessagePayloadTypeText(destinationNumber, message)
	jsonByte, err := json.Marshal(messagePayload)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(http.MethodPost, facebookConfig.RequestUrl, bytes.NewBuffer(jsonByte))
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+facebookConfig.RequestToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ctxhttp.Do(ctx, handler.externalClient, req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return false, errors.New("status is not ok")
	}

	return true, nil
}

func (handler *ClientHandler) SendMessageAudio(ctx context.Context, facebookConfig *config.FacebookConf, destinationNumber, audioId string) (bool, error) {
	log.Debug().Msg("[WhatsappCallback] entering whatsapp client audio...")

	messagePayload := GetMessagePayloadTypeAudio(destinationNumber, audioId)
	jsonByte, err := json.Marshal(messagePayload)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(http.MethodPost, facebookConfig.RequestUrl, bytes.NewBuffer(jsonByte))
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+facebookConfig.RequestToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ctxhttp.Do(ctx, handler.externalClient, req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	log.Debug().Msg(resp.Status)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return false, errors.New("status is not ok")
	}

	return true, nil
}
