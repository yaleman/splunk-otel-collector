// Copyright Splunk, Inc.
// Copyright The OpenTelemetry Authors
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

// Program otelcol is the OpenTelemetry Collector that collects stats
// and traces and exports to a configured backend.
package main

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/envprovider"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
)

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	uris := os.Args[1:]

	envProvider := envprovider.New()
	fileProvider := fileprovider.New()

	providers := map[string]confmap.Provider{
		"file":                      fileProvider,
		"this.is.a.file.uri.scheme": fileProvider,
		"this.is.an.env.uri.scheme": envProvider,
		"env":                       envProvider,
	}

	resolver, err := confmap.NewResolver(confmap.ResolverSettings{
		URIs:       uris,
		Providers:  providers,
		Converters: nil,
	})
	panicIfNotNil(err)
	config, err := resolver.Resolve(context.Background())
	panicIfNotNil(err)
	stringMap := config.ToStringMap()
	panicIfNotNil(err)
	fmt.Printf("config as map[string]any: %v\n", stringMap)
	resolver.Shutdown(context.Background())
}
