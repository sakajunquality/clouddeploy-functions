package clouddeployfunctions

import (
	"context"
	"time"

	"github.com/sakajunquality/clouddeploy-functions/clouddeploy"
	"github.com/sakajunquality/clouddeploy-functions/operations"
)

type PubSubMessage struct {
	ID          string            `json:"id,omitempty"`
	Data        []byte            `json:"Data,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	PublishTime time.Time         `json:"PublishTime,omitempty"`
}

// AutoPromote is an entrypoint of Cloud Functions
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
