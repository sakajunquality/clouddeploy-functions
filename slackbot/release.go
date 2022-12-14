package slackbot

import (
	"context"
	"fmt"
	"strings"

	"github.com/sakajunquality/clouddeploy-functions/clouddeploy"
	"github.com/sakajunquality/clouddeploy-functions/pubsub/approvals"
	"github.com/sakajunquality/clouddeploy-functions/pubsub/operations"
	"github.com/sakajunquality/clouddeploy-functions/slackbot/state"
	"github.com/slack-go/slack"

	"github.com/rs/zerolog/log"
)

func (s *Slackbot) NotifyReleaseUpdate(ctx context.Context, ops *operations.Operation) error {
	switch ops.Action {
	case operations.OperationsActionStart:
		return s.createReleasePost(ctx, ops)
	// TBD
	case operations.OperationsActionFailure:
		fallthrough
	case operations.OperationsActionSucceed:
		log.Info().Msg(fmt.Sprintf("%s is not supported, skipping", ops.Action))
	}

	return nil
}

func (s *Slackbot) createReleasePost(ctx context.Context, ops *operations.Operation) error {
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
		Value: string(ops.Action.GetPastParticiple()),
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
			Title:      fmt.Sprintf("[Release] %s %s", ops.DeliveryPipelineId, ops.ReleaseId),
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

func (s *Slackbot) NotifyRolloutUpdate(ctx context.Context, ops *operations.Operation) error {
	ts, err := s.getReleasePostTS(ctx, ops.DeliveryPipelineId, ops.ReleaseId)
	if err != nil {
		log.Error().Err(err).Msg("failed to get release post ts")
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
		Value: string(ops.Action.GetPastParticiple()),
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
			Title:      fmt.Sprintf("[%s] rollout for %s ", ops.Action, ops.TargetId),
			TitleLink:  ops.GetDeliveryPipelineURL(),
			Text:       fmt.Sprintf("Rollout has been %s. Go to the Cloud Console for %s.", ops.Action.GetPastParticiple(), possibleAction),
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

func (s *Slackbot) getReleasePostTS(ctx context.Context, pipelineID, releaseID string) (*string, error) {
	return state.NewReleaseStete(s.stateBucket, pipelineID, releaseID).GetTS(ctx)
}

func (s *Slackbot) saveReleasePostTS(ctx context.Context, pipelineID, releaseID, ts string) error {
	return state.NewReleaseStete(s.stateBucket, pipelineID, releaseID).SaveTS(ctx, ts)
}

func (s *Slackbot) NotifyApprovalThread(ctx context.Context, approval *approvals.Approval) error {
	pipelineName, err := approval.GetPipelineName()
	if err != nil {
		return err
	}

	ts, err := s.getReleasePostTS(ctx, *pipelineName, approval.ReleaseId)
	if err != nil {
		return err
	}

	api := slack.New(s.token)
	_, _, err = api.PostMessageContext(
		ctx,
		s.channel,
		slack.MsgOptionAttachments(slack.Attachment{
			Color:      getApprovalColor(string(approval.Action)),
			Title:      fmt.Sprintf("[Approval] %s", approval.Action),
			TitleLink:  approval.GetReleaseURL(),
			Text:       fmt.Sprintf("Approval is now %s for %s", approval.Action.GetPastParticiple(), approval.TargetId),
			AuthorName: "Cloud Deploy",
		}),
		slack.MsgOptionAsUser(true),
		slack.MsgOptionTS(*ts),
	)
	return err
}
