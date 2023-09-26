package clouddeploy

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	deploy "cloud.google.com/go/deploy/apiv1"
	"cloud.google.com/go/deploy/apiv1/deploypb"
)

type Rollout struct {
	ProjectNumber string
	PipelineID    string
	Location      string
	TargetID      string
	RolloutID     string
	ReleaseID     string
}

func (r *Rollout) GetPipelineResourceID() string {
	return fmt.Sprintf("projects/%s/locations/%s/deliveryPipelines/%s", r.ProjectNumber, r.Location, r.PipelineID)
}

func (r *Rollout) GetReleaseID() string {
	return fmt.Sprintf("projects/%s/locations/%s/deliveryPipelines/%s/releases/%s", r.ProjectNumber, r.Location, r.PipelineID, r.ReleaseID)
}

func CreateRollout(ctx context.Context, client *deploy.CloudDeployClient, r *Rollout) error {
	promoteReq := &deploypb.CreateRolloutRequest{
		Parent:    r.GetReleaseID(),
		RolloutId: r.RolloutID,
		Rollout: &deploypb.Rollout{
			TargetId: r.TargetID,
		},
	}
	op, err := client.CreateRollout(ctx, promoteReq)
	if err != nil {
		log.Error().Err(err)
		return nil
	}

	_, err = op.Wait(ctx)
	log.Error().Err(err)
	return err
}

func GetNextStageTargetID(ctx context.Context, client *deploy.CloudDeployClient, pipelineID string, currentStage string) (*string, error) {
	pipeline, err := client.GetDeliveryPipeline(ctx, &deploypb.GetDeliveryPipelineRequest{
		Name: pipelineID,
	})
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	stages := pipeline.GetSerialPipeline().GetStages()
	var nextStageIndex int
	for i, s := range stages {
		if s.GetTargetId() == currentStage {
			nextStageIndex = i + 1
			break
		}
	}

	if nextStageIndex == 0 || nextStageIndex > len(stages) {
		return nil, errors.New("next stage not found")
	}

	return &stages[nextStageIndex].TargetId, nil
}
