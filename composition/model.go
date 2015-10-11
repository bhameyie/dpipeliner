package composition

import "github.com/bhameyie/dpipeliner/data"

//IComposer produces content needed for creation of docker compose file
type IComposer interface {
	PrepareComposerContent(candidates []data.DeploymentCandidate) ([]byte, error)
	PrepareFinalizableCandidatesSnapshotContent(candidates []data.DeploymentCandidate) ([]byte, error)
}

//NonValidatedCandidates represents candidates for E2E testing
type NonValidatedCandidates struct {
	Service string
	Version string
}

type composeSpec struct {
	Image string
}
