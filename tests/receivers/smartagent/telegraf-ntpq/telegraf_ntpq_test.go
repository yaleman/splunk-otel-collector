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
	"path/filepath"
	"testing"

	"github.com/signalfx/splunk-otel-collector/tests/testutils"
)

func TestTelegrafNTPQProvidesAllMetrics(t *testing.T) {
	//tarball := "splunk-otel-collector_0.74.0_amd64.tar.gz"
	tarball := "splunk-otel-collector_0.78.1_amd64.tar.gz"

	testutils.AssertAllMetricsReceived(
		t, "all.yaml", "", nil,
		[]testutils.CollectorBuilder{func(collector testutils.Collector) testutils.Collector {
			cc := collector.(*testutils.CollectorContainer)
			cc.Container = cc.Container.WithContext(
				filepath.Join(".", "testdata", "collector-with-ntpq"),
			).WithBuildArgs(map[string]*string{
				"SPLUNK_OTEL_COLLECTOR_TARBALL": &tarball,
			})
			return cc.WithArgs()
		}})
}
