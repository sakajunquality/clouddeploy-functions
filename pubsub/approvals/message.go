package approvals

import (
	"errors"
	"fmt"
	"strings"
)

type ApprovalAction string

type Approval struct {
	Action        ApprovalAction
	Rollout       string
	ReleaseId     string
	RolloutId     string
	TargetId      string
	Location      string
	ProjectNumber string
}

const (
	ApprovalActionRequired ApprovalAction = "Required"
	ApprovalActionApproved ApprovalAction = "Approved"
	ApprovalActionRejected ApprovalAction = "Rejected"
)

func GetApprovalByAttributes(attributes map[string]string) *Approval {
	var approval Approval
	if v, ok := attributes["Action"]; ok {
		approval.Action = ApprovalAction(v)
	}
	if v, ok := attributes["Rollout"]; ok {
		approval.Rollout = v
	}
	if v, ok := attributes["ReleaseId"]; ok {
		approval.ReleaseId = v
	}
	if v, ok := attributes["RolloutId"]; ok {
		approval.RolloutId = v
	}
	if v, ok := attributes["TargetId"]; ok {
		approval.TargetId = v
	}
	if v, ok := attributes["Location"]; ok {
		approval.Location = v
	}
	if v, ok := attributes["ProjectNumber"]; ok {
		approval.ProjectNumber = v
	}

	return &approval
}

func (a *Approval) GetPipelineName() (*string, error) {
	parts := strings.Split(a.Rollout, "/")
	// projects/[project_id]/locations/[region]/deliveryPipelines/[pipeline name]/releases/[release name]/rollouts/[rollout name]
	if len(parts) != 10 {
		return nil, errors.New("failed to get pipeline name from rollout id")
	}
	return &parts[5], nil
}

func (a *ApprovalAction) GetPastParticiple() string {
	switch *a {
	case ApprovalActionRequired:
		return "required"
	case ApprovalActionApproved:
		return "approved"
	case ApprovalActionRejected:
		return "rejected"
	}

	return ""
}

func (a *Approval) GetReleaseURL() string {
	pipelineName, _ := a.GetPipelineName()
	base := fmt.Sprintf("https://console.cloud.google.com/deploy/delivery-pipelines/%s/%s/", a.Location, *pipelineName)
	return fmt.Sprintf("%sreleases/%s/rollouts?project=%s", base, a.ReleaseId, a.ProjectNumber)
}
