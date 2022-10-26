package slacker

import (
	"context"
	"fmt"
	"strings"

	"github.com/sakajunquality/clouddeploy-functions/clouddeploy"
	"github.com/sakajunquality/clouddeploy-functions/operations"
	"github.com/sakajunquality/clouddeploy-functions/slacker/state"
	"github.com/slack-go/slack"
)

func (s *Slacker) NotifyReleaseUpdate(ctx context.Context, ops *operations.Operation) error {
	switch ops.Action {
	case operations.OperationsActionStart:
		return s.createReleasePost(ctx, ops)
	// TBD
	case operations.OperationsActionFailure:
	case operations.OperationsActionSucceed:
	}

	return nil
}

func (s *Slacker) createReleasePost(ctx context.Context, ops *operations.Operation) error {
	release, err := clouddeploy.GetRelease(ctx, &clouddeploy.Rollout{
		ProjectNumber: ops.ProjectNumber,
		PipelineID:    ops.DeliveryPipelineId,
		Location:      ops.Location,
		TargetID:      ops.TargetId,
		RolloutID:     ops.RolloutId,
		ReleaseID:     ops.ReleaseId,
	})
	if err != nil {
		return err
	}

	fields := make([]slack.AttachmentField, 0)
	fields = append(fields, slack.AttachmentField{
		Title: "Pipeline",
		Value: string(ops.DeliveryPipelineId),
		Short: false,
	})

	fields = append(fields, slack.AttachmentField{
		Title: "Status",
		Value: string(ops.Action),
		Short: true,
	})

	fields = append(fields, slack.AttachmentField{
		Title: "Version",
		Value: string(ops.ReleaseId),
		Short: true,
	})

	if release.Description != "" {
		fields = append(fields, slack.AttachmentField{
			Title: "Description",
			Value: release.Description,
			Short: false,
		})
	}

	fields = append(fields, slack.AttachmentField{
		Title: "Link",
		Value: fmt.Sprintf("<%s|Release> / <%s|Pipeline>", ops.GetReleaseURL(), ops.GetDeliveryPipelineURL()),
		Short: false,
	})

	var labels []string
	for k, v := range release.Labels {
		labels = append(labels, fmt.Sprintf("%s: %s", k, v))
	}
	if len(labels) > 0 {
		fields = append(fields, slack.AttachmentField{
			Title: "Lables",
			Value: strings.Join(labels, "\n"),
			Short: false,
		})
	}

	var deployerID string
	var ccIDs string
	var annotations []string
	for k, v := range release.Annotations {
		if k == "deployer-slack-id" {
			deployerID = v
			continue
		}
		if k == "cc-slack-group-ids" {
			ccIDs = v
			continue
		}

		annotations = append(annotations, fmt.Sprintf("%s: %s", k, v))
	}
	if len(annotations) > 0 {
		fields = append(fields, slack.AttachmentField{
			Title: "Annotations",
			Value: strings.Join(annotations, "\n"),
			Short: false,
		})
	}

	var txt string
	if deployerID != "" {
		txt = fmt.Sprintf("Hey <@%s>!\nYou have initiated the release of %s.\n", deployerID, ops.DeliveryPipelineId)
	}
	txt += "Check this thread for rollouts' status.\n"

	if ccIDs != "" {
		ids := strings.Split(ccIDs, ",")
		txt += "cc"
		for _, id := range ids {
			txt += fmt.Sprintf(" <!subteam^%s>", id)
		}
	}

	msg := slack.MsgOptionAttachments(
		slack.Attachment{
			Title:      fmt.Sprintf("Release has been started for %s", ops.DeliveryPipelineId),
			TitleLink:  ops.GetDeliveryPipelineURL(),
			Text:       txt,
			Fields:     fields,
			AuthorName: "Cloud Deploy",
		},
	)

	api := slack.New(s.token)
	_, ts, err := api.PostMessageContext(
		ctx,
		s.channel,
		msg,

		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		return err
	}

	return s.saveReleasePostTS(ctx, ops.DeliveryPipelineId, ops.ReleaseId, ts)
}

func (s *Slacker) NotifyRolloutUpdate(ctx context.Context, ops *operations.Operation) error {
	ts, err := s.getReleasePostTS(ctx, ops.DeliveryPipelineId, ops.ReleaseId)
	if err != nil {
		return err
	}

	var color string
	var possibleAction string
	switch ops.Action {
	case operations.OperationsActionStart:
		color = "warning"
		possibleAction = "abandon"
	case operations.OperationsActionSucceed:
		color = "good"
		possibleAction = "rollback"
	case operations.OperationsActionFailure:
		color = "danger"
		possibleAction = "logs"
	}

	fields := make([]slack.AttachmentField, 0)
	fields = append(fields, slack.AttachmentField{
		Title: "Status",
		Value: string(ops.Action),
		Short: true,
	})

	fields = append(fields, slack.AttachmentField{
		Title: "Stage",
		Value: fmt.Sprintf("<%s|%s>", ops.GetTargetURL(), ops.TargetId),
		Short: true,
	})

	msg := slack.MsgOptionAttachments(
		slack.Attachment{
			Color:      color,
			Title:      fmt.Sprintf("Rollout %s for %s ", ops.Action, ops.TargetId),
			TitleLink:  ops.GetDeliveryPipelineURL(),
			Text:       fmt.Sprintf("Rollout has been %s. Go to the Cloud Console for %s.", ops.Action, possibleAction),
			AuthorName: "Cloud Deploy",
			Fields:     fields,
		},
	)

	api := slack.New(s.token)
	_, _, err = api.PostMessageContext(
		ctx,
		s.channel,
		msg,
		slack.MsgOptionAsUser(true),
		slack.MsgOptionTS(*ts),
	)
	return err
}

func (s *Slacker) getReleasePostTS(ctx context.Context, pipelineID, releaseID string) (*string, error) {
	return state.NewReleaseStete(s.stateBucket, pipelineID, releaseID).GetTS(ctx)
}

func (s *Slacker) saveReleasePostTS(ctx context.Context, pipelineID, releaseID, ts string) error {
	return state.NewReleaseStete(s.stateBucket, pipelineID, releaseID).SaveTS(ctx, ts)
}
