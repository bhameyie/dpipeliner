package data

// TrackedService represents a micro service that is allowed and tracked in the deployment pipeline
type TrackedService struct {
	Name        string `json:"Name" bson:"Name"`
	Description string `json:"Description" bson:"Description"`
}

// DeploymentCandidate represents candidate deployments that go through the deployment pipeline
type DeploymentCandidate struct {
	Image           string `json:"Image" bson:"Image"`
	Version         string `json:"Version" bson:"Version"`
	Started         int64  `json:"Started" bson:"Started"`
	Completed       bool   `json:"Completed" bson:"Completed"`
	Succeeded       bool   `json:"Succeeded" bson:"Succeeded"`
	MarathonVersion string `json:"MarathonVersion" bson:"MarathonVersion"`
	Unit            bool   `json:"Unit" bson:"Unit"`
	E2E             bool   `json:"E2E" bson:"E2E"`
	Deployed        bool   `json:"Deployed" bson:"Deployed"`
	MarathonSpec    string `json:"MarathonSpec" bson:"MarathonSpec"`
	ServiceName     string `json:"ServiceName" bson:"ServiceName"`
}

// IRepository defines the set of operations applicable to the tables/collection used through the pipeline
type IRepository interface {
	FindCandidate(name, version string) (DeploymentCandidate, error)
	CompleteStage(name, version, stage string) error
	RegisterNewCandidate(name, image, version string) error
	AssignMarathonSpecToCandidate(name, version, specContent string) error
	MarkCandidateAsSucceeded(name, version string) error
	GetCandidatesForE2E() ([]DeploymentCandidate, error)
	Dispose() error
}
