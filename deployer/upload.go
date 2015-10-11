package deployer

import (
	"encoding/json"
	"os"

	marathon "github.com/gambol99/go-marathon"
	"github.com/golang/glog"
	//"time"
)

//IDeployer deploys application
type IDeployer interface {
	Deploy(jsonContent []byte) (*ExpectedDeployment, error)
}

//MarathonDeployer deploys marathon apps
type MarathonDeployer struct {
	URL string
}

// ExpectedDeployment expected marathon deployment
type ExpectedDeployment struct {
	AppId         string
	NewDeployment bool
	DeploymentIds []string
}

func parseContent(jsonContent []byte) (res *marathon.Application, err error) {
	res = new(marathon.Application)
	err = json.Unmarshal(jsonContent, &res)
	return
}

func createNewApplication(client marathon.Marathon, app *marathon.Application) (deployed *ExpectedDeployment, err error) {
	deployed = &ExpectedDeployment{}
	if created, err := client.CreateApplication(app, false); err == nil {
		deployed.AppId = created.ID
		deployed.NewDeployment = true
		deps := created.DeploymentID
		ln := len(deps)
		deployed.DeploymentIds = make([]string, ln, ln)
		for i, el := range deps {
			deployed.DeploymentIds[i] = el["id"]
		}
	}

	return
}

func updateApplication(client marathon.Marathon, app *marathon.Application) (updated *ExpectedDeployment, err error) {
	updated = &ExpectedDeployment{}
	if updatedApp, err := client.UpdateApplication(app, false); err == nil {
		updated.AppId = app.ID
		updated.NewDeployment = false

		deps := updatedApp.DeploymentID
		ln := len(deps)
		updated.DeploymentIds = make([]string, ln, ln)
		for i, el := range deps {
			updated.DeploymentIds[i] = el["id"]
		}
	}
	return
}

// NewDeployer iniitializes a deployer
func NewDeployer(url string) IDeployer {
	return &MarathonDeployer{URL: url}
}

//Deploy deploys the marathon app
func (dep *MarathonDeployer) Deploy(jsonContent []byte) (*ExpectedDeployment, error) {
	if app, err := parseContent(jsonContent); err == nil {
		config := marathon.NewDefaultConfig()
		config.URL = dep.URL
		config.LogOutput = os.Stdout
		if client, err := marathon.NewClient(config); err == nil {
			if alreadyExists, err := client.HasApplication(app.ID); err == nil && alreadyExists {
				return updateApplication(client, app)
			} else if err == nil && !alreadyExists {
				return createNewApplication(client, app)
			} else {
				return nil, err
			}

		} else {

			glog.Fatalf("Failed to create a client for marathon, error: %s", err)
			return nil, err
		}

	}
	return nil, nil
}
