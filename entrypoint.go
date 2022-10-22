package clouddeployfunctions

import (
	"context"
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
	client := slacker.Slacker{
		Token:   os.Getenv("SLACK_TOKEN"),
		Channel: os.Getenv("SLACK_APPROVAL_CHANNEL"),
	}
	return client.NotifyApproval(ctx, approval)
}
