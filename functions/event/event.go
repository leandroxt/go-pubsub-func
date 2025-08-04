package event

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	functions.CloudEvent("ProcessEvent", ProcessEvent)
}

type MessagePublishedData struct {
	Message PubSubMessage
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type EventData struct {
	Email string `json:"email"`
}

type MessageWrapper struct {
	Subscription string `json:"subscription"`
	Message      struct {
		Data string `json:"data"` // Base64-encoded
	} `json:"message"`
}

func ProcessEvent(ctx context.Context, e event.Event) error {
	logger.Info("receiving data", slog.String("id", e.ID()))

	var msg MessageWrapper
	if err := e.DataAs(&msg); err != nil {
		logger.Error("cannot read DataAs")
		return fmt.Errorf("event.DataAs: %v", err)
	}

	logger.Info("message value", slog.String("value", msg.Message.Data))

	payload, err := base64.StdEncoding.DecodeString(msg.Message.Data)
	if err != nil {
		logger.Error("error decoding payload", slog.String("err", err.Error()))
		return err
	}

	logger.Info("payload value", slog.String("value", string(payload)))

	var msgDec MessageWrapper
	if err := json.Unmarshal(payload, &msgDec); err != nil {
		logger.Error("cannot parse the email json", slog.String("error", err.Error()))
		return err
	}

	var evtData EventData
	emailData, err := base64.StdEncoding.DecodeString(msgDec.Message.Data)
	if err := json.Unmarshal(emailData, &evtData); err != nil {
		logger.Error("cannot parse the email json unmarshalled", slog.String("error", err.Error()))
		return err
	}

	logger.Info("Hello", slog.String("email", evtData.Email))
	return nil
}
