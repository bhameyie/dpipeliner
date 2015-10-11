package data

import (
	"testing"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
)

func Test(t *testing.T) { TestingT(t) }

type RepoSuite struct{}

var _ = Suite(&RepoSuite{})

var session *mgo.Session
var sut = &CandidateRepository{Catalog: "testy"}

func buildSession() (*mgo.Session, error) {
	return mgo.Dial("localhost")
}

func (s *RepoSuite) SetUpSuite(c *C) {
	var err error
	if session, err = buildSession(); err == nil {
		sut.Session = session
		c := session.DB(dbName).C("testy_trackedservices")
		c.Insert(&TrackedService{Name: "cans"})
		c.Insert(&TrackedService{Name: "bottles"})
	} else {
		c.Fatal(err)
	}
}

func (s *RepoSuite) TearDownSuite(c *C) {
	if session != nil {
		c := session.DB(dbName).C("testy_trackedservices")
		c.DropCollection()
		session.Close()
	}
}

func (s *RepoSuite) TearDownTest(c *C) {
	if session != nil {
		db := session.DB(dbName)
		db.C("cans").DropCollection()
		db.C("bottles").DropCollection()
	}
}

func (s *RepoSuite) TestCanCompleteStageFromCriteria(c *C) {
	coll1 := session.DB(dbName).C("cans")
	ser1 := &DeploymentCandidate{Version: "v1"}
	c.Assert(coll1.Insert(ser1), IsNil)

	err := sut.CompleteStage("cans", "v1", "deployed")
	c.Assert(err, IsNil)

	cand, err2 := sut.FindCandidate("cans", "v1")
	c.Assert(err2, IsNil)
	c.Assert(cand.Deployed, Equals, true)
}

func (s *RepoSuite) TestCanCompleteStageFromCriteriaIndependentOfCasing(c *C) {
	coll1 := session.DB(dbName).C("cans")
	ser1 := &DeploymentCandidate{Version: "v1"}
	c.Assert(coll1.Insert(ser1), IsNil)

	err := sut.CompleteStage("cans", "v1", "DeplOyed")
	c.Assert(err, IsNil)

	cand, err2 := sut.FindCandidate("cans", "v1")
	c.Assert(err2, IsNil)
	c.Assert(cand.Deployed, Equals, true)
}

func (s *RepoSuite) TestCannotCompleteStageFromCriteriaWhenItemMissing(c *C) {
	err := sut.CompleteStage("bottles", "nada", "DeplOyed")
	c.Assert(err, NotNil)
}

func (s *RepoSuite) TestCannotCompleteStageFromCriteriaWhenStageInvalid(c *C) {
	err := sut.CompleteStage("bottles", "nada", "Mwahaha")
	c.Assert(err, NotNil)
}

func (s *RepoSuite) TestCanRegisterNewCandidatec(c *C) {
	err := sut.RegisterNewCandidate("bottles", "wolo", "loo")
	c.Assert(err, IsNil)

	cand, err2 := sut.FindCandidate("bottles", "loo")
	c.Assert(err2, IsNil)

	c.Assert(cand.Completed, Equals, false)
	c.Assert(cand.Unit, Equals, false)
	c.Assert(cand.E2E, Equals, false)
	c.Assert(cand.Succeeded, Equals, false)
	c.Assert(cand.Deployed, Equals, false)

	c.Assert(cand.Started, Not(Equals), 0)
	c.Assert(cand.ServiceName, Equals, "bottles")
	c.Assert(cand.MarathonSpec, Equals, "")
	c.Assert(cand.MarathonVersion, Equals, "")
	c.Assert(cand.Image, Equals, "wolo")
	c.Assert(cand.Version, Equals, "loo")
}

func (s *RepoSuite) TestCanFindCandidateFromCriteria(c *C) {
	coll1 := session.DB(dbName).C("cans")
	coll2 := session.DB(dbName).C("bottles")
	ser1 := &DeploymentCandidate{Version: "v1"}
	ser2 := &DeploymentCandidate{Version: "v2"}
	ser3 := &DeploymentCandidate{Version: "v3"}

	c.Assert(coll1.Insert(ser1), IsNil)
	c.Assert(coll1.Insert(ser2), IsNil)
	c.Assert(coll2.Insert(ser3), IsNil)

	cand, err := sut.FindCandidate("cans", "v2")
	c.Assert(err, IsNil)
	c.Assert(cand.Version, Equals, "v2")
}

func (s *RepoSuite) TestCanAssignMarathonSpecToCandidate(c *C) {
	coll1 := session.DB(dbName).C("cans")
	ser1 := &DeploymentCandidate{Version: "v1"}
	c.Assert(coll1.Insert(ser1), IsNil)

	err := sut.AssignMarathonSpecToCandidate("cans", "v1", "spec")
	c.Assert(err, IsNil)
	cand, err2 := sut.FindCandidate("cans", "v1")
	c.Assert(err2, IsNil)
	c.Assert(cand.MarathonSpec, Equals, "spec")
}

func (s *RepoSuite) TestFailsOnFindCandidateWhenNonePresent(c *C) {
	_, err := sut.FindCandidate("cans", "v5")
	c.Assert(err, NotNil)
}

func (s *RepoSuite) TestCanGetCandidatesForE2E(c *C) {
	coll1 := session.DB(dbName).C("cans")
	coll2 := session.DB(dbName).C("bottles")
	ser1 := &DeploymentCandidate{Version: "v1", Unit: true, MarathonSpec: "p"}
	ser2 := &DeploymentCandidate{Version: "v2", Unit: true}
	ser3 := &DeploymentCandidate{Version: "v3", Unit: true, MarathonSpec: "pp"}

	c.Assert(coll1.Insert(ser1), IsNil)
	c.Assert(coll1.Insert(ser2), IsNil)
	c.Assert(coll2.Insert(ser3), IsNil)

	cands, err := sut.GetCandidatesForE2E()
	c.Assert(err, IsNil)
	c.Assert(len(cands), Equals, 2)

}

func (s *RepoSuite) TestCanGetCandidatesForE2EEvenWhenNone(c *C) {
	cands, err := sut.GetCandidatesForE2E()
	c.Assert(err, IsNil)
	c.Assert(len(cands), Equals, 0)
}
