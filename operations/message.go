package operations

type operationsAction string
type resourceType string

type Operation struct {
	Action             operationsAction
	Resource           string
	ResourceType       resourceType
	Location           string
	DeliveryPipelineId string
	ProjectNumber      string
	ReleaseId          string
	TargetId           string
	RolloutId          string
}

const (
	operationsActionStart   operationsAction = "Start"
	operationsActionSucceed operationsAction = "Succeed"
	operationsActionFailure operationsAction = "Failure"

	resourceTypeDeliveryPipeline resourceType = "DeliveryPipeline"
	resourceTypeDeliveryTarget   resourceType = "Target"
	resourceTypeDeliveryRelease  resourceType = "Release"
	resourceTypeDeliveryRollout  resourceType = "Rollout"
	resourceTypeDeliveryJobRun   resourceType = "JobRun"
)

func GetOperationByAttributes(attributes map[string]string) *Operation {
	var ops Operation
	if v, ok := attributes["Action"]; ok {
		ops.Action = operationsAction(v)
	}
	if v, ok := attributes["Resource"]; ok {
		ops.Resource = v
	}
	if v, ok := attributes["ResourceType"]; ok {
		ops.ResourceType = resourceType(v)
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

func (o *Operation) DispachAutoPromote() bool {
	return o.TargetId != "" && o.Action == operationsActionSucceed && o.ResourceType == resourceTypeDeliveryRollout
}
