// Copyright Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package simpleprometheusremotewritereceiver

import (
	"errors"
	"time"

	"github.com/jaegertracing/jaeger/pkg/multierror"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/confignet"

	"github.com/signalfx/splunk-otel-collector/internal/receiver/simpleprometheusremotewritereceiver/internal/tools"
)

var _ component.Config = (*Config)(nil)

const (
	typeString = tools.TypeStr
)

type Config struct {
	// TODO the other guy has some really cool stuff in here
	ListenAddr confignet.NetAddr `mapstructure:",squash"`
	ListenPath string            `mapstructure:"path"`
	Timeout    time.Duration     `mapstructure:"timeout"`
	BufferSize int               `mapstructure:"buffer_size"` // Channel buffer size, defaults to blocking each request until processed
}

func (c *Config) Validate() error {
	var errs []error
	if c.ListenAddr.Endpoint == "" {
		errs = append(errs, errors.New("endpoint must not be empty"))
	}
	if c.ListenAddr.Transport == "" {
		errs = append(errs, errors.New("transport must not be empty"))
	}
	if c.Timeout < time.Second {
		errs = append(errs, errors.New("impractically short timeout"))
	}
	if err := componenttest.CheckConfigStruct(c); err != nil {
		errs = append(errs, err)
	}
	if errs != nil {
		return multierror.Wrap(errs)
	}
	return nil
}

//
//// Unmarshal a confmap.Conf into the config struct.
//func (c *Config) Unmarshal(conf *confmap.Conf) error {
//	if err := conf.Unmarshal(c, confmap.WithErrorUnused()); err != nil {
//		return err
//	}
//	ts := conf.Get("timeout")
//	switch timeout := ts.(type) {
//	case int, int8, int16, int32, int64:
//		to := timeout.(int)
//		c.Timeout = time.Second * time.Duration(to)
//	case string:
//		to, err := time.ParseDuration(timeout)
//		if nil != err {
//			return err
//		}
//		c.Timeout = to
//	}
//	return nil
//}
