package data

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CandidateRepository is the repository of deployment candidates
type CandidateRepository struct {
	Session *mgo.Session
	Catalog string
}

const dbName = "pipeline"

// FindCandidate Retrieves a candidate based on the given criterias
func (r *CandidateRepository) FindCandidate(name, version string) (DeploymentCandidate, error) {
	res := DeploymentCandidate{}
	coll := r.Session.DB(dbName).C(name)
	err := coll.Find(bson.M{"Version": version}).One(&res)
	return res, err
}

func ensureValidStage(stage string) (string, error) {
	if strings.EqualFold(stage, "Completed") {
		return "Completed", nil
	}
	if strings.EqualFold(stage, "Succeeded") {
		return "Succeeded", nil
	}
	if strings.EqualFold(stage, "Unit") {
		return "Unit", nil
	}
	if strings.EqualFold(stage, "E2E") {
		return "E2E", nil
	}
	if strings.EqualFold(stage, "Deployed") {
		return "Deployed", nil
	}
	return "", errors.New(stage + " is not a valid stage")
}

// CompleteStage mark a give stage on the pipeline as completed
func (r *CandidateRepository) CompleteStage(name, version, stage string) error {
	realStage, err := ensureValidStage(stage)
	if err != nil {
		return err
	}

	c := r.Session.DB(dbName).C(name)
	return c.Update(bson.M{"Version": version}, bson.M{"$set": bson.M{realStage: true}})
}

// RegisterNewCandidate starts a new deployment candidate
func (r *CandidateRepository) RegisterNewCandidate(name, image, version string) error {
	cand := &DeploymentCandidate{
		ServiceName: name,
		Image:       image,
		Version:     version,
		Started:     time.Now().Unix(),
	}

	c := r.Session.DB(dbName).C(name)
	index := mgo.Index{
		Key:      []string{"Version"},
		Unique:   true,
		DropDups: true,
		Sparse:   true,
	}
	c.EnsureIndex(index)
	return c.Insert(cand)
}

// AssignMarathonSpecToCandidate attaches a MarathonSpec to a given candidate
func (r *CandidateRepository) AssignMarathonSpecToCandidate(name, version, specContent string) error {
	c := r.Session.DB(dbName).C(name)
	return c.Update(bson.M{"Version": version}, bson.M{"$set": bson.M{"MarathonSpec": specContent}})
}

// MarkCandidateAsSucceeded mark a candidate deployment as having succeeded
func (r *CandidateRepository) MarkCandidateAsSucceeded(name, version string) error {
	return r.CompleteStage(name, version, "Completed")
}

// GetCandidatesForE2E gets candidates that have passed unit testing and have a marathon spec
func (r *CandidateRepository) GetCandidatesForE2E() ([]DeploymentCandidate, error) {
	var candidates []DeploymentCandidate
	servs, err := r.getTrackedServices()
	if err != nil {
		return nil, err
	}
	for _, s := range servs {
		found, accErr := r.accumulate(s.Name)
		if accErr != nil {
			return nil, accErr
		}
		candidates = append(candidates, found...)
	}
	return candidates, nil
}

func (r *CandidateRepository) accumulate(name string) ([]DeploymentCandidate, error) {
	var found []DeploymentCandidate
	c := r.Session.DB(dbName).C(name)
	crit := bson.M{
		"Unit":         true,
		"MarathonSpec": bson.M{"$ne": ""},
	}
	err := c.Find(crit).Sort("Started").All(&found)
	return found, err
}

func (r *CandidateRepository) getTrackedServices() ([]TrackedService, error) {
	c := r.Session.DB(dbName).C(r.Catalog + "_trackedservices")
	var res []TrackedService
	err := c.Find(bson.M{}).All(&res)
	return res, err
}

//Dispose closes the open session
func (r *CandidateRepository) Dispose() error {
	//todo should possibly surround with a recover
	r.Session.Close()
	return nil
}

//NewRepository creates a new repository
func NewRepository(url, catalog string) (IRepository, error) {
	fmt.Println("Connecting to server: " + url)
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &CandidateRepository{Session: session, Catalog: catalog}, nil
}
