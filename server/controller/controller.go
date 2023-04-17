package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"

	"github.com/sakajunquality/clouddeploy-functions/pubsub/approvals"
	"github.com/sakajunquality/clouddeploy-functions/pubsub/operations"
	"github.com/sakajunquality/clouddeploy-functions/slackbot"
)

var (
	client *slackbot.Slackbot
)

func init() {
	// @todo refactor
	client = slackbot.NewSlackbot(os.Getenv("SLACK_TOKEN"), os.Getenv("SLACK_CHANNEL"))
	client.SetStateBucket(os.Getenv("SLACK_BOT_STATE_BUCKET"))
}

func HandleOperationsTopic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m, err := readPubSubMessage(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ops := operations.GetOperationByAttributes(m.Attributes)
	switch ops.ResourceType {
	case operations.ResourceTypeRelease:
		if err := client.NotifyReleaseUpdate(ctx, ops); err != nil {
			log.Error().Err(err).Msg("failed to NotifyReleaseUpdate")
			return
		}
	case operations.ResourceTypeRollout:
		if err := client.NotifyRolloutUpdate(ctx, ops); err != nil {
			log.Error().Err(err).Msg("failed to NotifyReleaseUpdate")
			return
		}
	default:
		// not supported
		log.Debug().Msg("unsupported resource type")
	}
}

func HandleApprovalsTopic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m, err := readPubSubMessage(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	approval := approvals.GetApprovalByAttributes(m.Attributes)
	client.NotifyApprovalThread(ctx, approval)
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
