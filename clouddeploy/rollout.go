package clouddeploy

import (
	"context"
	"errors"
	"fmt"

	deploy "cloud.google.com/go/deploy/apiv1"
	deploypb "google.golang.org/genproto/googleapis/cloud/deploy/v1"
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

func (r *Rollout) GetParent() string {
	return fmt.Sprintf("projects/%s/locations/%s/deliveryPipelines/%s/releases/%s", r.ProjectNumber, r.Location, r.PipelineID, r.ReleaseID)
}

func CreateRollout(ctx context.Context, client *deploy.CloudDeployClient, r *Rollout) error {
	promoteReq := &deploypb.CreateRolloutRequest{
		Parent:    r.GetParent(),
		RolloutId: r.RolloutID,
		Rollout: &deploypb.Rollout{
			TargetId: r.TargetID,
		},
	}
	op, err := client.CreateRollout(ctx, promoteReq)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	_, err = op.Wait(ctx)
	fmt.Println(err)
	return err
}

func GetNextStageTargetID(ctx context.Context, client *deploy.CloudDeployClient, pipelineID string, currentStage string) (*string, error) {
	pipeline, err := client.GetDeliveryPipeline(ctx, &deploypb.GetDeliveryPipelineRequest{
		Name: pipelineID,
	})
	if err != nil {
		fmt.Println(err)
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
