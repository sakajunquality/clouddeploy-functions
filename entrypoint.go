package clouddeployfunctions

import (
	"context"
	"os"
	"time"

	"github.com/sakajunquality/clouddeploy-functions/clouddeploy"
	"github.com/sakajunquality/clouddeploy-functions/pubsub/approvals"
	"github.com/sakajunquality/clouddeploy-functions/pubsub/operations"
	"github.com/sakajunquality/clouddeploy-functions/slackbot"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/GoogleCloudPlatform/berglas/pkg/auto"
)

type PubSubMessage struct {
	ID          string            `json:"id,omitempty"`
	Data        []byte            `json:"Data,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	PublishTime time.Time         `json:"PublishTime,omitempty"`
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// AutoPromote is an entrypoint of Cloud Functions subscribing `clouddeploy-operations` topic.
func AutoPromote(ctx context.Context, m PubSubMessage) error {
	ops := operations.GetOperationByAttributes(m.Attributes)
	if !ops.DispachAutoPromote() {
		return nil
	}

	log.Debug().Msg("triggering AutoPromote")
	rollout := &clouddeploy.Rollout{
		ProjectNumber: ops.ProjectNumber,
		PipelineID:    ops.DeliveryPipelineId,
		Location:      ops.Location,
		TargetID:      ops.TargetId,
		RolloutID:     ops.RolloutId,
		ReleaseID:     ops.ReleaseId,
	}

	return clouddeploy.AutoPromote(ctx, rollout)
}

// NotifyApprovalRequestSlackSimple is an entrypoint of Cloud Functions subscribing `clouddeploy-approvals`
func NotifyApprovalRequestSlackSimple(ctx context.Context, m PubSubMessage) error {
	approval := approvals.GetApprovalByAttributes(m.Attributes)
	client := slackbot.NewSlackbot(os.Getenv("SLACK_TOKEN"), os.Getenv("SLACK_APPROVAL_CHANNEL"))
	return client.NotifyApproval(ctx, approval)
}

func NotifySlackWithThread(ctx context.Context, m PubSubMessage) error {
	log.Debug().Msg("running NotifySlackWithThread")
	client := slackbot.NewSlackbot(os.Getenv("SLACK_TOKEN"), os.Getenv("SLACK_CHANNEL"))
	client.SetStateBucket(os.Getenv("SLACK_BOT_STATE_BUCKET"))

	ops := operations.GetOperationByAttributes(m.Attributes)
	switch ops.ResourceType {
	case operations.ResourceTypeRelease:
		if err := client.NotifyReleaseUpdate(ctx, ops); err != nil {
			log.Error().Err(err).Msg("failed to NotifyReleaseUpdate")
			return err
		}
	case operations.ResourceTypeRollout:
		if err := client.NotifyRolloutUpdate(ctx, ops); err != nil {
			log.Error().Err(err).Msg("failed to NotifyReleaseUpdate")
			return err
		}
	default:
		// not supported
	}
	return nil
}

func NotifySlackApprovalWithThread(ctx context.Context, m PubSubMessage) error {
	log.Debug().Msg("running NotifySlackApprovalWithThread")

	approval := approvals.GetApprovalByAttributes(m.Attributes)
	client := slackbot.NewSlackbot(os.Getenv("SLACK_TOKEN"), os.Getenv("SLACK_CHANNEL"))
	client.SetStateBucket(os.Getenv("SLACK_BOT_STATE_BUCKET"))
	return client.NotifyApprovalThread(ctx, approval)
}
