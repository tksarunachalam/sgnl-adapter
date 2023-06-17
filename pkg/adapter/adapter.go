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
)

// Adapter implements the framework.Adapter interface to query pages of objects
// from datasources.
type Adapter struct {
	// SCAFFOLDING:
	// Add/remove fields as needed to configure this adapter.

	// Example field.
	Logger *log.Logger
}

// NewAdapter instantiates a new Adapter.
//
// SCAFFOLDING:
// Add/remove parameters as needed to configure this adapter.
func NewAdapter(logger *log.Logger) framework.Adapter[Config] {
	return &Adapter{
		Logger: logger,
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
