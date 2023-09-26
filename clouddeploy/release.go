package clouddeploy

import (
	"context"

	"github.com/rs/zerolog/log"

	deploy "cloud.google.com/go/deploy/apiv1"
	"cloud.google.com/go/deploy/apiv1/deploypb"
)

func GetRelease(ctx context.Context, r *Rollout) (*deploypb.Release, error) {
	client, err := deploy.NewCloudDeployClient(ctx)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	defer client.Close()

	return client.GetRelease(ctx, &deploypb.GetReleaseRequest{
		Name: r.GetReleaseID(),
	})
}
