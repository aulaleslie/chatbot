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

func (handler *ClientHandler) SendMessageText(ctx context.Context, facebookConfig *config.FacebookConf, requestBody WebHookPayload) (bool, error) {
	log.Debug().Msg("[WhatsappCallback] entering whatsapp client text...")

	messageAndDestinationNumber := getMessageAndDestinationNumber(requestBody)
	if messageAndDestinationNumber == nil {
		return false, errors.New("could not obtain parameter")
	}

	messagePayload := GetMessagePayloadTypeText(messageAndDestinationNumber.DestinationNumber, messageAndDestinationNumber.Message)
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

func (handler *ClientHandler) SendMessageAudio(ctx context.Context, facebookConfig *config.FacebookConf, requestBody WebHookPayload) (bool, error) {
	log.Debug().Msg("[WhatsappCallback] entering whatsapp client audio...")

	messageAndDestinationNumber := getMessageAndDestinationNumber(requestBody)
	if messageAndDestinationNumber == nil {
		return false, errors.New("could not obtain parameter")
	}

	messagePayload := GetMessagePayloadTypeAudio(messageAndDestinationNumber.DestinationNumber, facebookConfig.AudioId)
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

func getMessageAndDestinationNumber(webhookPayload WebHookPayload) *MessageAndDestinationNumber {
	if len(webhookPayload.Entry) > 0 {
		entry := webhookPayload.Entry[0]

		if len(entry.Changes) > 0 {
			change := entry.Changes[0]

			if len(change.Value.Messages) > 0 {
				message := change.Value.Messages[0]

				log.Debug().Msg(message.From)
				log.Debug().Msg(message.Text.Body)

				return &MessageAndDestinationNumber{
					DestinationNumber: message.From,
					Message:           message.Text.Body,
				}
			}
		}
	}

	return nil
}
