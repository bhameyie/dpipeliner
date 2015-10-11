package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bhameyie/dpipeliner/data"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ControllerSuite struct{}

var _ = Suite(&ControllerSuite{})

const (
	snapper       = "candidateSnapper.json"
	compo         = "docker-compose.yml"
	snapJsContent = `[
      {
          "service": "boom",
          "version": "1"
      },
      {
          "service": "doom",
          "version": "12"
      }
  ]
  `
)

func (s *ControllerSuite) TearDownTest(c *C) {
	os.Remove(snapper)
	os.Remove(compo)
}

func (s *ControllerSuite) TestCanUpdateStateFromCandidateSnapshot(c *C) {
	rep := &AllGoodRepo{}
	sut := &Controller{Repo: rep}

	err := sut.updateStateForNonValidatedCandidates("Completed", snapJsContent)
	c.Assert(err, IsNil)

	ss := rep.Spies[0]
	sss := rep.Spies[1]

	c.Assert(ss.ServiceName, Equals, "boom")
	c.Assert(ss.ServiceVersion, Equals, "1")
	c.Assert(ss.StageName, Equals, "Completed")
	c.Assert(ss.StageCompleted, Equals, true)

	c.Assert(sss.ServiceName, Equals, "doom")
	c.Assert(sss.ServiceVersion, Equals, "12")
	c.Assert(sss.StageName, Equals, "Completed")
	c.Assert(sss.StageCompleted, Equals, true)
}

func (s *ControllerSuite) TestCanCompleteAStage(c *C) {
	rep := &AllGoodRepo{}
	sut := &Controller{Repo: rep}
	err := sut.CompleteStageFor("a", "po", "wolo")
	ss := rep.Spies[0]
	c.Assert(err, IsNil)
	c.Assert(ss.ServiceName, Equals, "a")
	c.Assert(ss.ServiceVersion, Equals, "po")
	c.Assert(ss.StageName, Equals, "wolo")
	c.Assert(ss.StageCompleted, Equals, true)
}

func (s *ControllerSuite) TestCanProduceCompositionAndSnapshotFiles(c *C) {
	sut := &Controller{Composer: &AllGoodComposer{}, Repo: &AllGoodRepo{}}

	err := sut.ProduceCompositionAndSnapshotFiles()
	c.Assert(err, IsNil)

	yml, errYml := ioutil.ReadFile(compo)
	c.Assert(errYml, IsNil)
	c.Assert(string(yml), Equals, "hoho")

	js, errJs := ioutil.ReadFile(snapper)
	c.Assert(errJs, IsNil)
	c.Assert(string(js), Equals, "hooo")

}

//stubs

type RepoSpy struct {
	StageCompleted bool
	StageName      string
	ServiceName    string
	ServiceVersion string
	ServiceImage   string
}

type AllGoodRepo struct {
	Spies []RepoSpy
}

func (s *AllGoodRepo) CompleteStage(name, version, stage string) error {
	ss := RepoSpy{}
	ss.StageCompleted = true
	ss.StageName = stage
	ss.ServiceName = name
	ss.ServiceVersion = version
	s.Spies = append(s.Spies, ss)
	return nil
}

func (s *AllGoodRepo) RegisterNewCandidate(name, image, version string) error {
	return nil
}
func (s *AllGoodRepo) AssignMarathonSpecToCandidate(name, version, specContent string) error {
	return nil
}
func (s *AllGoodRepo) MarkCandidateAsSucceeded(name, version string) error {
	return nil
}
func (s *AllGoodRepo) GetCandidatesForE2E() ([]data.DeploymentCandidate, error) {
	return nil, nil
}

func (s *AllGoodRepo) FindCandidate(name, version string) (data.DeploymentCandidate, error) {
	return data.DeploymentCandidate{}, nil
}

func (s *AllGoodRepo) Dispose() error {
	return nil
}

func (s *AllGoodRepo) GetTrackedServices() ([]string, error) {
	return nil, nil
}

type AllGoodComposer struct {
}

func (s *AllGoodComposer) PrepareComposerContent(candidates []data.DeploymentCandidate) ([]byte, error) {
	return []byte("hoho"), nil
}

func (s *AllGoodComposer) PrepareFinalizableCandidatesSnapshotContent(candidates []data.DeploymentCandidate) ([]byte, error) {
	return []byte("hooo"), nil
}
