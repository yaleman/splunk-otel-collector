// Copyright Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build integration

package tests

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/signalfx/splunk-otel-collector/tests/testutils"
)

func TestContainerEndpointProperties(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("unable to share sockets between mac and d4m vm: https://github.com/docker/for-mac/issues/483#issuecomment-758836836")
	}

	containers := []testutils.Container{
		testutils.NewContainer().WithImage("redis").WithLabels(
			map[string]string{
				"io.opentelemetry.collector.receiver-creator.smartagent/redis.rule":          `type == "container"`,
				"io.opentelemetry.collector.receiver-creator.smartagent/redis.config.type":   "collectd/redis",
				"io.opentelemetry.collector.receiver-creator.smartagent/redis.res_attrs": `{"attr.one": "attr_one_value", "attr.two": "attr_two_value"}`,
			},
		).WithName("redis-container-name").WillWaitForLogs("Ready to accept connections"),
	}

	testutils.AssertAllMetricsReceived(
		t, "docker-observer-smart-agent-redis.yaml", "docker-observer-config.yaml",
		containers, []testutils.CollectorBuilder{
			func(c testutils.Collector) testutils.Collector {
				cc := c.(*testutils.CollectorContainer)
				cc.Container = cc.Container.WithBinds("/var/run/docker.sock:/var/run/docker.sock:ro")
				cc.Container = cc.Container.WithUser(fmt.Sprintf("999:%d", testutils.GetDockerGID(t)))
				return cc
			},
		},
	)
}
