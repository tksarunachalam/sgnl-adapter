// Copyright 2023 SGNL.ai, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package adapter

import (
	"context"
	"log"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapter-template/pkg/example_datasource"
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from datasources.
type Adapter struct {
	// SCAFFOLDING:
	// Add/remove fields below as needed to configure this adapter.

	// Logger is a standard logger.
	Logger *log.Logger

	// ExampleClient provides access to the example datasource.
	ExampleClient example_datasource.Client
}

// NewAdapter instantiates a new Adapter.
//
// SCAFFOLDING:
// Add/remove parameters as needed to configure this adapter.
func NewAdapter(logger *log.Logger, client example_datasource.Client) framework.Adapter[Config] {
	return &Adapter{
		Logger:        logger,
		ExampleClient: client,
	}
}

// GetPage is called by SGNL's ingestion service to query a page of objects
// from a datasource.
func (a *Adapter) GetPage(ctx context.Context, request *framework.Request[Config]) framework.Response {
	a.Logger.Printf("Received GetPage call")

	if err := a.ValidateGetPageRequest(ctx, request); err != nil {
		return framework.NewGetPageResponseError(err)
	}

	// TODO
	return framework.NewGetPageResponseError(nil)
}
