package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sakajunquality/clouddeploy-functions/pubsub/approvals"
	"github.com/sakajunquality/clouddeploy-functions/pubsub/operations"
)

type PubSubMessage struct {
	Message struct {
		Attributes  map[string]string `json:"attributes,omitempty"`
		MessageId   string            `json:"messageId,omitempty"`
		PublishTime time.Time         `json:"PublishTime,omitempty"`
	} `json:"message,omitempty"`
}

func HandleOperationsTopic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m, err := readPubSubMessage(r)
	if err != nil {
		log.Error().Err(err).Msg("failed to readPubSubMessage")
		w.WriteHeader(http.StatusOK) // ignore empty payload
		return
	}

	ops := operations.GetOperationByAttributes(m.Message.Attributes)
	switch ops.ResourceType {
	case operations.ResourceTypeRelease:
		if err := client.NotifyReleaseUpdate(ctx, ops); err != nil {
			logger.Error().Err(err).Msg("failed to NotifyReleaseUpdate")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case operations.ResourceTypeRollout:
		if err := client.NotifyRolloutUpdate(ctx, ops); err != nil {
			logger.Error().Err(err).Msg("failed to NotifyRolloutUpdate")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		// not supported
		logger.Info().Msg(fmt.Sprintf("unsupported resource type: %s", ops.ResourceType))
	}

	w.WriteHeader(http.StatusOK)
}

func HandleApprovalsTopic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m, err := readPubSubMessage(r)
	if err != nil {
		w.WriteHeader(http.StatusOK) // ignore empty payload
		return
	}

	approval := approvals.GetApprovalByAttributes(m.Message.Attributes)
	err = client.NotifyApprovalThread(ctx, approval)
	if err != nil {
		logger.Error().Err(err).Msg("failed to NotifyApprovalThread")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func readPubSubMessage(r *http.Request) (*PubSubMessage, error) {
	var m PubSubMessage
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, err
	}

	return &m, nil
}
