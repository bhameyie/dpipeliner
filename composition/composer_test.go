package composition

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/bhameyie/dpipeliner/data"

	. "gopkg.in/check.v1"
	"gopkg.in/yaml.v2"
)

func Test(t *testing.T) { TestingT(t) }

type ComposerSuite struct{}

var _ = Suite(&ComposerSuite{})

var sut = &SystemComposer{}

func (s *ComposerSuite) TestPrepareComposerContentProducesValidFileWhenCandidatesAvailable(c *C) {
	cand1 := data.DeploymentCandidate{
		Image:       "hey",
		ServiceName: "yo",
	}

	cand2 := data.DeploymentCandidate{
		Image:       "look",
		ServiceName: "here",
	}
	arr := []data.DeploymentCandidate{cand1, cand2}
	content, err := sut.PrepareComposerContent(arr)
	c.Assert(err, IsNil)

	m := make(map[string]composeSpec)
	err = yaml.Unmarshal([]byte(content), &m)
	c.Assert(err, IsNil)

	c.Assert(m["yo"].Image, Equals, "hey")
	c.Assert(m["here"].Image, Equals, "look")
}

func (s *ComposerSuite) TestPrepareComposerContentProducesFailsWhenNoCandidateOrNil(c *C) {
	arr := []data.DeploymentCandidate{}
	content, err := sut.PrepareComposerContent(arr)
	c.Assert(content, IsNil)
	c.Assert(err, NotNil)

	content, err = sut.PrepareComposerContent(nil)
	c.Assert(content, IsNil)
	c.Assert(err, NotNil)
}

func (s *ComposerSuite) TestPrepareFinalizableCandidatesSnapshotContentProduceValidJsonForRelevantCandidates(c *C) {
	cand1 := data.DeploymentCandidate{
		Image:   "hey",
		Version: "1",
		E2E:     true,
	}

	cand2 := data.DeploymentCandidate{
		Version:     "2",
		ServiceName: "here",
	}

	cand3 := data.DeploymentCandidate{
		Version:     "3",
		ServiceName: "there",
		Completed:   true,
	}

	cand4 := data.DeploymentCandidate{
		Version:     "4",
		ServiceName: "nothere",
	}

	arr := []data.DeploymentCandidate{cand1, cand2, cand3, cand4}
	content, err := sut.PrepareFinalizableCandidatesSnapshotContent(arr)
	c.Assert(err, IsNil)

	var candidates []NonValidatedCandidates
	dec := json.NewDecoder(bytes.NewReader(content))
	err = dec.Decode(&candidates)
	c.Assert(err, IsNil)
	c.Assert(len(candidates), Equals, 2)
	c.Assert(candidates[0].Service, Equals, "here")
	c.Assert(candidates[0].Version, Equals, "2")
	c.Assert(candidates[1].Service, Equals, "nothere")
	c.Assert(candidates[1].Version, Equals, "4")
}

func (s *ComposerSuite) TestPrepareFinalizableCandidatesSnapshotContentFailsWhenNilOrEmpty(c *C) {
	arr := []data.DeploymentCandidate{}
	content, err := sut.PrepareFinalizableCandidatesSnapshotContent(arr)
	c.Assert(content, IsNil)
	c.Assert(err, NotNil)

	content, err = sut.PrepareFinalizableCandidatesSnapshotContent(nil)
	c.Assert(content, IsNil)
	c.Assert(err, NotNil)
}
