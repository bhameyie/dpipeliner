package main

import (
	"dpipeliner/composition"
	"dpipeliner/data"
	"dpipeliner/deployer"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

const snapshotFile = "candidateSnapper.json"
const composerFile = "docker-compose.yml"

// Controller performs defined operations using internal constructs
type Controller struct {
	Repo     data.IRepository
	Composer composition.IComposer
	Deployer deployer.IDeployer
}

func readNonValidatedCandidates(content string) (candidates []composition.NonValidatedCandidates, err error) {
	dec := json.NewDecoder(strings.NewReader(content))
	err = dec.Decode(&candidates)
	return
}

func (c *Controller) updateStateForNonValidatedCandidates(state, fileContent string) error {
	candidates, err := readNonValidatedCandidates(fileContent)
	if err != nil {
		return err
	}

	for _, candidate := range candidates {
		if err := c.CompleteStageFor(candidate.Service, candidate.Version, state); err != nil {
			return err
		}
	}
	return nil
}

// Dispose terminates all the resources associated with the controller
func (c *Controller) Dispose() error {
	return c.Repo.Dispose()
}

// DeploySnapshot deploys all candidates from the snapshotFile
func (c *Controller) DeploySnapshot() error {
	b, err := ioutil.ReadFile(snapshotFile)
	if err != nil {
		return err
	}
	content := string(b)
	cands, err2 := readNonValidatedCandidates(content)
	if err2 != nil {
		return err2
	}

	for _, cc := range cands {
		if err := c.TriggerCandidateDeployment(cc.Service, cc.Version); err != nil {
			return err
		}
	}
	return nil
}

// AcceptCandidateSnapshot marks as succesful all the versions listed in the snapshot file
func (c *Controller) AcceptCandidateSnapshot() error {
	return c.ChangeCandidateState("Succeeded")
}

// CompleteCandidateSnapshot marks as completed all the versions listed in the snapshot file
func (c *Controller) CompleteCandidateSnapshot() error {
	return c.ChangeCandidateState("Completed")
}

// ChangeCandidateState updates the state candidates in a snapsot
func (c *Controller) ChangeCandidateState(state string) error {
	b, err := ioutil.ReadFile(snapshotFile)
	if err != nil {
		return err
	}

	return c.updateStateForNonValidatedCandidates(state, string(b))
}

// AssignMarathonSpecificationFor assigns the content of a marathonspec to a candidate
func (c *Controller) AssignMarathonSpecificationFor(name, version, marathonSpec string) error {
	b, err := ioutil.ReadFile(marathonSpec)
	if err != nil {
		return err
	}
	return c.Repo.AssignMarathonSpecToCandidate(name, version, string(b))
}

func writeSnapshotFile(c composition.IComposer, candidates []data.DeploymentCandidate) error {
	content, err := c.PrepareFinalizableCandidatesSnapshotContent(candidates)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(snapshotFile, content, 0644)
}

func writeDockerComposeFile(c composition.IComposer, candidates []data.DeploymentCandidate) error {
	content, err := c.PrepareComposerContent(candidates)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(composerFile, content, 0644)
}

// ProduceCompositionAndSnapshotFiles produces a docker compose file and candidate snapshot file
func (c *Controller) ProduceCompositionAndSnapshotFiles() error {
	candidates, err := c.Repo.GetCandidatesForE2E()
	if err != nil {
		return err
	}
	if err := writeDockerComposeFile(c.Composer, candidates); err != nil {
		return err
	}
	return writeSnapshotFile(c.Composer, candidates)
}

// CompleteStageFor marks a given stage as completed for the chosen candidate
func (c *Controller) CompleteStageFor(name, version, stage string) error {
	return c.Repo.CompleteStage(name, version, stage)
}

// TriggerCandidateDeployment attempts to deploy a candidate to marathon
func (c *Controller) TriggerCandidateDeployment(name, version string) error {
	candidate, err := c.Repo.FindCandidate(name, version)
	if err != nil {
		return err
	}
	if deployment, err := c.Deployer.Deploy([]byte(candidate.MarathonSpec)); err != nil {
		return err
	} else {
		fmt.Println("Deployed " + deployment.AppId + " with version " + version)
		return c.CompleteStageFor(name, version, "Deployed")
	}
}

// StartPipeline initiates candidate registration
func (c *Controller) StartPipeline(name, version, image string) error {
	return c.Repo.RegisterNewCandidate(name, image, version)
}
