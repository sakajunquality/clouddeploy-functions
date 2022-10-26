package operations

import "fmt"

type OperationsAction string
type ResourceType string

type Operation struct {
	Action             OperationsAction
	Resource           string
	ResourceType       ResourceType
	Location           string
	DeliveryPipelineId string
	ProjectNumber      string
	ReleaseId          string
	TargetId           string
	RolloutId          string
}

const (
	OperationsActionStart   OperationsAction = "Start"
	OperationsActionSucceed OperationsAction = "Succeed"
	OperationsActionFailure OperationsAction = "Failure"

	ResourceTypeDeliveryPipeline ResourceType = "DeliveryPipeline"
	ResourceTypeTarget           ResourceType = "Target"
	ResourceTypeRelease          ResourceType = "Release"
	ResourceTypeRollout          ResourceType = "Rollout"
	ResourceTypeJobRun           ResourceType = "JobRun"
)

func GetOperationByAttributes(attributes map[string]string) *Operation {
	var ops Operation
	if v, ok := attributes["Action"]; ok {
		ops.Action = OperationsAction(v)
	}
	if v, ok := attributes["Resource"]; ok {
		ops.Resource = v
	}
	if v, ok := attributes["ResourceType"]; ok {
		ops.ResourceType = ResourceType(v)
	}
	if v, ok := attributes["Location"]; ok {
		ops.Location = v
	}
	if v, ok := attributes["DeliveryPipelineId"]; ok {
		ops.DeliveryPipelineId = v
	}
	if v, ok := attributes["ProjectNumber"]; ok {
		ops.ProjectNumber = v
	}
	if v, ok := attributes["ReleaseId"]; ok {
		ops.ReleaseId = v
	}
	if v, ok := attributes["TargetId"]; ok {
		ops.TargetId = v
	}
	if v, ok := attributes["RolloutId"]; ok {
		ops.RolloutId = v
	}

	return &ops
}

func (o *Operation) GetDeliveryPipelineURL() string {
	base := o.getConsoleBaseURL()
	return fmt.Sprintf("%s?project=%s", base, o.ProjectNumber)
}

func (o *Operation) GetReleaseURL() string {
	base := o.getConsoleBaseURL()
	return fmt.Sprintf("%sreleases/%s/rollouts?project=%s", base, o.ReleaseId, o.ProjectNumber)
}

func (o *Operation) GetTargetURL() string {
	base := o.getConsoleBaseURL()
	return fmt.Sprintf("%stargets/%s?project=%s", base, o.TargetId, o.ProjectNumber)
}

func (o *Operation) getConsoleBaseURL() string {
	return fmt.Sprintf("https://console.cloud.google.com/deploy/delivery-pipelines/%s/%s/", o.Location, o.DeliveryPipelineId)
}

func (o *Operation) DispachAutoPromote() bool {
	return o.TargetId != "" && o.Action == OperationsActionSucceed && o.ResourceType == ResourceTypeRollout
}

func (o *OperationsAction) GetPastParticiple() string {
	switch *o {
	case OperationsActionStart:
		return "started"
	case OperationsActionSucceed:
		return "succeeded"
	case OperationsActionFailure:
		return "failed"
	}

	return ""
}
