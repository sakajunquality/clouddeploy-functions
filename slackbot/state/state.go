package state

import (
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"github.com/rs/zerolog/log"
)

type ReleaseState struct {
	bucketName string
	pipelineID string
	releaseID  string
}

func NewReleaseStete(bucketName, pipelineID, ReleaseID string) *ReleaseState {
	return &ReleaseState{
		bucketName: bucketName,
		pipelineID: pipelineID,
		releaseID:  ReleaseID,
	}
}

func (s *ReleaseState) GetTS(ctx context.Context) (*string, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	obj := client.Bucket(s.bucketName).Object(s.statePath())
	r, err := obj.NewReader(ctx)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	defer r.Close()

	tsBtytes, err := ioutil.ReadAll(r)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	ts := string(tsBtytes)
	return &ts, nil
}

func (s *ReleaseState) SaveTS(ctx context.Context, ts string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	w := client.Bucket(s.bucketName).Object(s.statePath()).NewWriter(ctx)
	defer w.Close()

	_, err = w.Write([]byte(ts))
	return err
}

func (s *ReleaseState) statePath() string {
	return fmt.Sprintf("state/%s/%s", s.pipelineID, s.releaseID)
}
