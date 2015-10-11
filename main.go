package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/bhameyie/dpipeliner/composition"
	"github.com/bhameyie/dpipeliner/data"
	"github.com/bhameyie/dpipeliner/deployer"
)

func ensureValidSpec(serviceName, serviceVersion string) error {
	if serviceName == "-1" || serviceVersion == "-1" {
		return errors.New("invalid service or version")
	}
	return nil
}

func notNegative(s, e string) error {
	if s == "-1" {
		return errors.New(e)
	}
	return nil
}

func fileExists(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

func main() {

	modePtr := flag.String("mode", "deploy", "e.g. deploy, init_test, complete_state, compose")
	marathonPtr := flag.String("marathon", "-1", "marathon host")
	mongoPtr := flag.String("mongo", "-1", "location of mongo server")
	file := flag.String("file", "marathon.spec.js", "location of spec file for attach_spec mode")

	serviceImage := flag.String("image", "-1", "service image")
	serviceName := flag.String("service", "-1", "service name")
	serviceVersion := flag.String("version", "-1", "service version")
	stage := flag.String("stage", "-1", "e.g. unit, e2e, deployment")
	catalog := flag.String("catalog", "-1", "tracked service collection (e.g. fire_trackedservices)")

	flag.Parse()

	fmt.Println("DPipeliner")
	fmt.Println("")

	repo, err := data.NewRepository(*mongoPtr, *catalog)
	if err != nil {
		panic(err)
	}
	controller := &Controller{
		Repo:     repo,
		Deployer: deployer.NewDeployer(*marathonPtr),
		Composer: composition.NewComposer(),
	}

	defer controller.Dispose()

	validateSpec := ensureValidSpec(*serviceName, *serviceVersion)
	validateImage := notNegative(*serviceImage, "invalid image")
	validateStage := notNegative(*stage, "invalid stage")

	var e error
	switch *modePtr {
	case "deploy":
		fmt.Println("deploying")
		if validateSpec == nil {
			e = controller.TriggerCandidateDeployment(*serviceName, *serviceVersion)
		} else {
			e = validateSpec
		}

	case "compose":
		fmt.Println("composing")
		e = controller.ProduceCompositionAndSnapshotFiles()

	case "complete_snapshot":
		if fileExists(snapshotFile) {
			e = controller.CompleteCandidateSnapshot()
			if e == nil {
				e = controller.ChangeCandidateState("Deployed")
			}
		} else {
			e = errors.New(snapshotFile + " doesnt exist")
		}

	case "deploy_snapshot":
		if fileExists(snapshotFile) {
			e = controller.DeploySnapshot()
		} else {
			e = errors.New(snapshotFile + " doesnt exist")
		}

	case "accept_snapshot":
		if fileExists(snapshotFile) {
			e = controller.AcceptCandidateSnapshot()
		} else {
			e = errors.New(snapshotFile + " doesnt exist")
		}

	case "attach_spec":
		fmt.Println("Attaching marathon spec")
		if fileExists(*file) {
			e = controller.AssignMarathonSpecificationFor(*serviceName, *serviceVersion, *file)
		} else {
			e = errors.New("file doesnt exist")
		}

	case "init_test":
		fmt.Println("initiating tests")
		if validateSpec == nil {
			if validateImage == nil {
				fmt.Print(*serviceName + " - " + *serviceVersion)
				e = controller.StartPipeline(*serviceName, *serviceVersion, *serviceImage)
			} else {
				e = validateImage
			}
		} else {
			e = validateSpec
		}

	case "complete_state":
		if validateSpec == nil {
			if validateStage == nil {
				fmt.Print(*serviceName + " - " + *serviceVersion)
				e = controller.CompleteStageFor(*serviceName, *serviceVersion, *stage)
			} else {
				e = validateStage
			}
		} else {
			e = validateSpec
		}

	default:
		panic("unrecognized mode: " + *modePtr)
	}

	if e != nil {
		fmt.Println(e)
		panic(e.Error())
	}

}
