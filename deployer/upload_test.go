package deployer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const jsonContent = `{
    "id": "elApp",
    "cmd": "env && sleep 300",
    "args": ["/bin/sh", "-c", "env && sleep 300"],
    "cpus": 1.5,
    "mem": 256.0,
    "ports": [
        8080,
        9000
    ],
    "requirePorts": false,
    "instances": 3,
    "executor": "",
    "container": {
        "type": "DOCKER",
        "docker": {
            "image": "group/image",
            "network": "BRIDGE",
            "portMappings": [
                {
                    "containerPort": 8080,
                    "hostPort": 0,
                    "servicePort": 9000,
                    "protocol": "tcp"
                },
                {
                    "containerPort": 161,
                    "hostPort": 0,
                    "protocol": "udp"
                }
            ],
            "privileged": false,
            "parameters": [
                { "key": "a-docker-option", "value": "xxx" },
                { "key": "b-docker-option", "value": "yyy" }
            ]
        },
        "volumes": [
            {
                "containerPath": "/etc/a",
                "hostPath": "/var/data/a",
                "mode": "RO"
            },
            {
                "containerPath": "/etc/b",
                "hostPath": "/var/data/b",
                "mode": "RW"
            }
        ]
    },
    "env": {
        "LD_LIBRARY_PATH": "/usr/local/lib/myLib"
    },
    "constraints": [
        ["attribute", "OPERATOR", "value"]
    ],
    "acceptedResourceRoles": [
        "role1", "*"
    ],
    "labels": {
        "environment": "staging"
    },
    "uris": [
        "https://raw.github.com/mesosphere/marathon/master/README.md"
    ],
    "dependencies": ["/product/db/mongo", "/product/db", "../../db"],
    "healthChecks": [
        {
            "protocol": "HTTP",
            "path": "/health",
            "gracePeriodSeconds": 3,
            "intervalSeconds": 10,
            "portIndex": 0,
            "timeoutSeconds": 10,
            "maxConsecutiveFailures": 3
        },
        {
            "protocol": "TCP",
            "gracePeriodSeconds": 3,
            "intervalSeconds": 5,
            "portIndex": 1,
            "timeoutSeconds": 5,
            "maxConsecutiveFailures": 3
        },
        {
            "protocol": "COMMAND",
            "command": { "value": "curl -f -X GET http://$HOST:$PORT0/health" },
            "maxConsecutiveFailures": 3
        }
    ],
    "backoffSeconds": 1,
    "backoffFactor": 1.15,
    "maxLaunchDelaySeconds": 3600,
    "upgradeStrategy": {
        "minimumHealthCapacity": 0.5,
        "maximumOverCapacity": 0.2
    }
}`

func TestCanProduceApplicationFromJson(t *testing.T) {
	bb := []byte(jsonContent)
	app, err := parseContent(bb)
	assert.Nil(t, err, "should not throw")

	assert.Equal(t, "elApp", app.ID, "should have the same id")
}
