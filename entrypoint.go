package clouddeployfunctions

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sakajunquality/clouddeploy-functions/approvals"
	"github.com/sakajunquality/clouddeploy-functions/clouddeploy"
	"github.com/sakajunquality/clouddeploy-functions/operations"
	"github.com/sakajunquality/clouddeploy-functions/slacker"

	_ "github.com/GoogleCloudPlatform/berglas/pkg/auto"
)

type PubSubMessage struct {
	ID          string            `json:"id,omitempty"`
	Data        []byte            `json:"Data,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	PublishTime time.Time         `json:"PublishTime,omitempty"`
}

// AutoPromote is an entrypoint of Cloud Functions subscribing `clouddeploy-operations` topic.
func AutoPromote(ctx context.Context, m PubSubMessage) error {
	ops := operations.GetOperationByAttributes(m.Attributes)
	if !ops.DispachAutoPromote() {
		return nil
	}

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
	client := slacker.NewSlacker(os.Getenv("SLACK_TOKEN"), os.Getenv("SLACK_APPROVAL_CHANNEL"))
	return client.NotifyApproval(ctx, approval)
}

func NotifySlackWithThread(ctx context.Context, m PubSubMessage) error {
	client := slacker.NewSlacker(os.Getenv("SLACK_TOKEN"), os.Getenv("SLACK_CHANNEL"))
	client.SetStateBucket(os.Getenv("SLACK_BOT_STATE_BUCKET"))

	ops := operations.GetOperationByAttributes(m.Attributes)
	switch ops.ResourceType {
	case operations.ResourceTypeRelease:
		if err := client.NotifyReleaseUpdate(ctx, ops); err != nil {
			// fix logger
			fmt.Println(err)
		}
	case operations.ResourceTypeRollout:
		if err := client.NotifyRolloutUpdate(ctx, ops); err != nil {
			// fix logger
			fmt.Println(err)
		}
	default:
		// not supported
	}
	return nil
}
