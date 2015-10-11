TEST?=./...
Version = 0.3
default: test

test:
	go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

updatedeps:
	go get -u github.com/gambol99/go-marathon
	go get -u github.com/golang/glog
	go get -u github.com/stretchr/testify/assert
	go get -u gopkg.in/mgo.v2
	go get -u gopkg.in/yaml.v2
	go get -u github.com/mitchellh/gox

release:
	mkdir _release
	go build
	mv dpipeliner _release/dpipeliner.$(Version)

clean:
	rm -r _release/
