package composition

import (
	"dpipeliner/data"
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v2"
)

//SystemComposer produces content needed for creation of docker compose file
type SystemComposer struct {
}

//NewComposer initializes a new composer
func NewComposer() IComposer {
	return &SystemComposer{}
}

// PrepareComposerContent produce docker-compose file content
func (sc *SystemComposer) PrepareComposerContent(candidates []data.DeploymentCandidate) ([]byte, error) {
	if candidates == nil || len(candidates) == 0 {
		return nil, errors.New("No candidates found")
	}

	com := make(map[string]composeSpec)
	for _, candidate := range candidates {
		com[candidate.ServiceName] = composeSpec{Image: candidate.Image}
	}

	return yaml.Marshal(com)
}

// PrepareFinalizableCandidatesSnapshotContent produces candidateSnapper content
func (sc *SystemComposer) PrepareFinalizableCandidatesSnapshotContent(candidates []data.DeploymentCandidate) ([]byte, error) {
	if candidates == nil || len(candidates) == 0 {
		return nil, errors.New("No candidates found")
	}

	var com []NonValidatedCandidates
	for _, candidate := range candidates {
		if !candidate.E2E && !candidate.Completed {
			com = append(com,
				NonValidatedCandidates{
					Service: candidate.ServiceName,
					Version: candidate.Version,
				})
		}
	}
	return json.Marshal(com)
}
