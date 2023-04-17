package main

import (
	"encoding/json"
	"io"
	"net/http"

	"cloud.google.com/go/pubsub"

	"github.com/sakajunquality/clouddeploy-functions/pubsub/approvals"
	"github.com/sakajunquality/clouddeploy-functions/pubsub/operations"
)

func HandleOperationsTopic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m, err := readPubSubMessage(r)
	if err != nil {
		w.WriteHeader(http.StatusOK) // ignoreing empty payload
		return
	}

	ops := operations.GetOperationByAttributes(m.Attributes)
	switch ops.ResourceType {
	case operations.ResourceTypeRelease:
		if err := client.NotifyReleaseUpdate(ctx, ops); err != nil {
			logger.Error().Err(err).Msg("failed to NotifyReleaseUpdate")
			return
		}
	case operations.ResourceTypeRollout:
		if err := client.NotifyRolloutUpdate(ctx, ops); err != nil {
			logger.Error().Err(err).Msg("failed to NotifyReleaseUpdate")
			return
		}
	default:
		// not supported
		logger.Debug().Msg("unsupported resource type")
	}

	w.WriteHeader(http.StatusOK)
}

func HandleApprovalsTopic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m, err := readPubSubMessage(r)
	if err != nil {
		w.WriteHeader(http.StatusOK) // ignoreing empty payload
		return
	}

	approval := approvals.GetApprovalByAttributes(m.Attributes)
	err = client.NotifyApprovalThread(ctx, approval)
	if err != nil {
		logger.Error().Err(err).Msg("failed to NotifyApprovalThread")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func readPubSubMessage(r *http.Request) (*pubsub.Message, error) {
	var m pubsub.Message
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &m); err != nil {
		return nil, err
	}

	return &m, nil
}
