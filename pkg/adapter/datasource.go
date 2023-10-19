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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

const (
	// SCAFFOLDING:
	// Update the set of valid entity types supported by this adapter.

	Users  string = "users"
	Groups string = "groups"
)

// Entity contains entity specific information, such as the entity's unique ID attribute and the
// endpoint to query that entity.
type Entity struct {
	// SCAFFOLDING:
	// Add or remove fields as needed. This should be used to store entity specific information
	// such as the entity's unique ID attribute name and the endpoint to query that entity.

	// uniqueIDAttrExternalID is the external ID of the entity's uniqueId attribute.
	uniqueIDAttrExternalID string
}

var (
	// SCAFFOLDING:
	// Using the consts defined above, update the set of valid entity types supported by this adapter.

	// ValidEntityExternalIDs is a map of valid external IDs of entities that can be queried.
	// The map value is the Entity struct which contains the unique ID attribute.
	ValidEntityExternalIDs = map[string]Entity{
		Users: {
			uniqueIDAttrExternalID: "user_id",
		},
		Groups: {
			uniqueIDAttrExternalID: "group_id",
		},
	}
)

// Datasource directly implements a Client interface to allow querying
// an external datasource.
type Datasource struct {
	Client *http.Client
}

// NewClient returns a Client to query the datasource.
func NewClient(timeout int) Client {
	return &Datasource{
		Client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	var req *http.Request

	// SCAFFOLDING:
	// Populate the request with the appropriate path, headers, and query parameters to query the
	// datasource.
	url := fmt.Sprintf("%s/api/%s", request.BaseURL, request.EntityExternalID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, &framework.Error{
			Message: "Failed to create HTTP request to datasource.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Timeout API calls that take longer than 5 seconds
	apiCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req = req.WithContext(apiCtx)

	// Add headers to the request, if any.
	// req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Authorization", "Bearer Token")

	res, err := d.Client.Do(req)
	if err != nil {
		return nil, &framework.Error{
			Message: "Failed to send request to datasource.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	response := &Response{
		StatusCode:       res.StatusCode,
		RetryAfterHeader: res.Header.Get("Retry-After"),
	}

	if res.StatusCode != http.StatusOK {
		return response, nil
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: "Failed to read response body.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	objects, nextCursor, parseErr := ParseResponse(body, request.EntityExternalID, request.PageSize)
	if parseErr != nil {
		return nil, parseErr
	}

	response.Objects = objects
	response.NextCursor = nextCursor

	return response, nil
}

func ParseResponse(
	body []byte, entityExternalID string, pageSize int64,
) (objects []map[string]any, nextCursor string, err *framework.Error) {
	var data map[string]any

	unmarshalErr := json.Unmarshal(body, &data)
	if unmarshalErr != nil {
		return nil, "", &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// SCAFFOLDING:
	// Replace `response` with the field name in the datasource response that contains the
	// list of objects.
	rawData, found := data["response"]
	if !found {
		return nil, "", &framework.Error{
			Message: "Field missing in the datasource response: response.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	rawObjects, ok := rawData.([]any)
	if !ok {
		return nil, "", &framework.Error{
			Message: fmt.Sprintf(
				"Entity %s field exists in the datasource response but field value is not a list of objects: %T.",
				entityExternalID,
				rawData,
			),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	parsedObjects, parserErr := parseObjects(rawObjects, entityExternalID)
	if parserErr != nil {
		return nil, "", parserErr
	}

	// SCAFFOLDING:
	// Populate nextCursor with the cursor returned from the datasource, if present.
	nextCursor = ""

	return parsedObjects, nextCursor, nil
}

// parseObjects parses []any into []map[string]any. If any object in the slice is not a map[string]any,
// a framework.Error is returned.
func parseObjects(objects []any, entityExternalId string) ([]map[string]any, *framework.Error) {
	parsedObjects := make([]map[string]any, 0, len(objects))

	for _, object := range objects {
		parsedObject, ok := object.(map[string]any)
		if !ok {
			return nil, &framework.Error{
				Message: fmt.Sprintf(
					"An object in Entity: %s could not be parsed. Expected: map[string]any. Got: %T.",
					entityExternalId,
					object,
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		parsedObjects = append(parsedObjects, parsedObject)
	}

	return parsedObjects, nil
}
