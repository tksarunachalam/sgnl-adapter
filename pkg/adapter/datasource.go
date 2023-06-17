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
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"github.com/sgnl-ai/adapter-template/pkg/example_datasource"
)

// SCAFFOLDING:
// Update the set of error messages.
const (
	ErrMsgExampleDatasourceErrorFmt                = "Example datasource returned an error: %v"
	ErrMsgExampleDatasourceStatusCodeFmt           = "Example datasource returned unexpected status code: %d"
	ErrMsgExampleDatasourceInvalidAttributeTypeFmt = "Example datasource returned an attribute with an incompatible type: %s"
)

// RequestPageFromDatasource requests a page of objects from a datasource.
func (a *Adapter) RequestPageFromDatasource(ctx context.Context, request *framework.Request[Config]) framework.Response {
	// SCAFFOLDING:
	// Modify the implementation of this method to perform a query to your
	// real datasource. This example implementation query an in-memory
	// example datasource that returns JSON objects.

	exampleRequest := &example_datasource.Request{
		URL:      fmt.Sprintf("%s/%s/%s", request.Address, request.Config.DatasourceVersion, request.Entity.ExternalId),
		Username: request.Auth.Basic.Username,
		Password: request.Auth.Basic.Password,
		PageSize: request.PageSize,
		Cursor:   request.Cursor,
	}

	a.Logger.Printf("Querying example datasource at URL %s", exampleRequest.URL)

	exampleResponse, err := a.ExampleClient.GetPage(ctx, exampleRequest)

	if err != nil {
		a.Logger.Printf("Example datasource query failed: %v", err)

		return framework.NewGetPageResponseError(
			&framework.Error{
				Message: fmt.Sprintf(ErrMsgExampleDatasourceErrorFmt, err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
			},
		)
	}

	adapterErr := web.HTTPError(exampleResponse.StatusCode, exampleResponse.RetryAfterHeader)
	if adapterErr != nil {
		a.Logger.Printf("Example datasource query returned failure status code %d", exampleResponse.StatusCode)

		return framework.NewGetPageResponseError(adapterErr)
	}

	page := &framework.Page{
		NextCursor: exampleResponse.NextCursor,
	}

	page.Objects, err = web.ConvertJSONObjectList(&request.Entity, exampleResponse.Objects)
	if err != nil {
		a.Logger.Printf("Failed to parse JSON objects returned by the datasource: %v", err)

		return framework.NewGetPageResponseError(
			&framework.Error{
				Message: fmt.Sprintf(ErrMsgExampleDatasourceInvalidAttributeTypeFmt, err.Error()),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
			},
		)
	}

	return framework.NewGetPageResponseSuccess(page)
}
