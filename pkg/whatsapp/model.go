package whatsapp

import (
	"chatbot/config"
	"context"
	"net/http"
)

const WebhooksEndpoint = "/whatsapp/webhooks"

type WebHookPayload struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Value Value  `json:"value"`
	Field string `json:"field"`
}

type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Contact struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

type Profile struct {
	Name string `json:"name"`
}

type Message struct {
	From      string   `json:"from"`
	Id        string   `json:"id"`
	Timestamp string   `json:"timestamp"`
	Text      Text     `json:"text,omitempty"`
	Reaction  Reaction `json:"reaction,omitempty"`
	Image     Image    `json:"image,omitempty"`
	Type      string   `json:"type"`
}

type Text struct {
	Body string `json:"body"`
}

type Reaction struct {
	MessageId string `json:"message_id"`
	Emoji     string `json:"emoji"`
}

type Image struct {
	Caption  string `json:"caption"`
	MimeType string `json:"mime_type"`
	Sha256   string `json:"sha256"`
	Id       string `json:"id"`
}

type CallbackHandler struct {
	facebookConfig *config.FacebookConf
	client         IWhatsappClient
}

type VerificationHandler struct {
	facebookConfig *config.FacebookConf
	client         IWhatsappClient
}

type CallbackPayload struct {
	Message           string `json:"message"`
	DestinationNumber string `json:"destination_number"`
}

type IWhatsappClient interface {
	SendMessageText(ctx context.Context, facebookConfig *config.FacebookConf, requestBody WebHookPayload) (bool, error)
	SendMessageAudio(ctx context.Context, facebookConfig *config.FacebookConf, requestBody WebHookPayload) (bool, error)
}

type ClientHandler struct {
	externalClient *http.Client
}

type CallbackRequest struct {
}

type MediaType string

const (
	MediaImage    MediaType = "WhatsApp Image Keys"
	MediaVideo    MediaType = "WhatsApp Video Keys"
	MediaAudio    MediaType = "WhatsApp Audio Keys"
	MediaDocument MediaType = "WhatsApp Document Keys"
)

type MessagePayload struct {
	MessagingProdut string        `json:"messaging_product"`
	RecipientType   string        `json:"recipient_type"`
	To              string        `json:"to"`
	Type            string        `json:"type"`
	Text            *TextPayload  `json:"text,omitempty"`
	Audio           *AudioPayload `json:"audio,omitempty"`
}

type TextPayload struct {
	PreviewUrl bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type AudioPayload struct {
	Id string `json:"id"`
}

func GetMessagePayloadTypeText(destinationNumber, message string) MessagePayload {
	return MessagePayload{
		MessagingProdut: "whatsapp",
		RecipientType:   "individual",
		To:              destinationNumber,
		Type:            "text",
		Text: &TextPayload{
			PreviewUrl: false,
			Body:       message,
		},
	}
}

func GetMessagePayloadTypeAudio(destinationNumber, audioId string) MessagePayload {
	return MessagePayload{
		MessagingProdut: "whatsapp",
		RecipientType:   "individual",
		To:              destinationNumber,
		Type:            "audio",
		Audio: &AudioPayload{
			Id: audioId,
		},
	}
}

type MessageAndDestinationNumber struct {
	Message           string
	DestinationNumber string
}
