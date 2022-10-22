package clouddeploy

import (
	"context"

	deploy "cloud.google.com/go/deploy/apiv1"
)

type Client struct {
	client deploy.CloudDeployClient
}

func NewClient(ctx context.Context) (*Client, error) {
	client, err := deploy.NewCloudDeployClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return &Client{
		client: *client,
	}, nil
}
