package approvals

import (
	"errors"
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
	parts := strings.Split(a.RolloutId, "/")
	// projects/[project_id]/locations/[region]/deliveryPipelines/[pipeline name]/releases/[release name]/rollouts/[rollout name]
	if len(parts) != 10 {
		return nil, errors.New("failed to get pipeline name from rollout id")
	}
	return &parts[5], nil
}
