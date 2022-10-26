package slacker

import (
	"context"
	"fmt"

	"github.com/sakajunquality/clouddeploy-functions/approvals"

	"github.com/slack-go/slack"
)

func (s *Slacker) NotifyApproval(ctx context.Context, approval *approvals.Approval) error {
	api := slack.New(s.token)
	_, _, err := api.PostMessage(
		s.channel,
		slack.MsgOptionAttachments(slack.Attachment{
			Color: getApprovalColor(string(approval.Action)),
			Text:  fmt.Sprintf("Rollout is now %s for %s", approval.Action, approval.Rollout),
		}),
		slack.MsgOptionAsUser(true),
	)

	return err
}

func getApprovalColor(action string) string {
	var color string
	switch approvals.ApprovalAction(action) {
	case approvals.ApprovalActionRequired:
		color = "#85C1E9"
	case approvals.ApprovalActionApproved:
		color = "#58D68D"
	case approvals.ApprovalActionRejected:
		color = "#F5B041"
	default:
		color = "#BFC9CA"
	}
	return color
}
