package clouddeploy

import (
	"context"
	"fmt"
	"time"

	deploy "cloud.google.com/go/deploy/apiv1"
	"github.com/rs/zerolog/log"
)

func AutoPromote(ctx context.Context, r *Rollout) error {
	client, err := deploy.NewCloudDeployClient(ctx)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	defer client.Close()

	nextTargetID, err := GetNextStageTargetID(ctx, client, r.GetPipelineResourceID(), r.TargetID)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	nextRollout := r
	r.TargetID = *nextTargetID
	r.RolloutID = getNextRolollouID(r.ReleaseID, *nextTargetID)

	return CreateRollout(ctx, client, nextRollout)
}

func getNextRolollouID(releaseID, targetID string) string {
	now := time.Now()
	return fmt.Sprintf("%s-to-%s-%d", releaseID, targetID, now.Unix())
}
