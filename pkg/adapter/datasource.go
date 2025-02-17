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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

const (
	// SCAFFOLDING #11 - pkg/adapter/datasource.go: Update the set of valid entity types this adapter supports.
	Users   string = "users"
	Vendors string = "vendors"
	Teams   string = "teams"
)

// Entity contains entity specific information, such as the entity's unique ID attribute and the
// endpoint to query that entity.
type Entity struct {
	// SCAFFOLDING #12 - pkg/adapter/datasource.go: Update Entity fields used to store entity specific information
	// Add or remove fields as needed. This should be used to store entity specific information
	// such as the entity's unique ID attribute name and the endpoint to query that entity.

	// uniqueIDAttrExternalID is the external ID of the entity's uniqueId attribute.
	uniqueIDAttrExternalID string
	endPoint               string
}

// Datasource directly implements a Client interface to allow querying
// an external datasource.
type Datasource struct {
	Client *http.Client
}

type DatasourceResponse struct {
	// SCAFFOLDING #13  - pkg/adapter/datasource.go: Add or remove fields in the response as necessary. This is used to unmarshal the response from the SoR.

	// SCAFFOLDING #14 - pkg/adapter/datasource.go: Update `objects` with field name in the SoR response that contains the list of objects.
	Objects []map[string]any `json:"-"`
	Limit   int              `json:"limit"`
	Offset  int              `json:"offset"`
	More    bool             `json:"more"`
}

var (
	// SCAFFOLDING #15 - pkg/adapter/datasource.go: Update the set of valid entity types supported by this adapter. Used for validation.

	// ValidEntityExternalIDs is a map of valid external IDs of entities that can be queried.
	// The map value is the Entity struct which contains the unique ID attribute.
	ValidEntityExternalIDs = map[string]Entity{
		Users: {
			uniqueIDAttrExternalID: "id",
			endPoint:               Users,
		},
		Vendors: {
			uniqueIDAttrExternalID: "id",
			endPoint:               Vendors,
		},
		Teams: {
			uniqueIDAttrExternalID: "id",
			endPoint:               Teams,
		},
	}
)

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

	// SCAFFOLDING #16 - pkg/adapter/datasource.go: Create the SoR API URL
	// Populate the request with the appropriate path, headers, and query parameters to query the
	// datasource.
	baseUrl, err := url.Parse(fmt.Sprintf("%s/%s", request.BaseURL, ValidEntityExternalIDs[request.EntityExternalID].endPoint))
	if err != nil {
		return nil, &framework.Error{
			Message: "Failed to parse the base URL.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Add query parameters to the URL, if any.
	addQueryParams(baseUrl, request)

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, baseUrl.String(), nil)
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

	// SCAFFOLDING #17 - pkg/adapter/datasource.go: Add any headers required to communicate with the SoR APIs.
	// Add headers to the request, if any.
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	if request.Token != "" {
		req.Header.Add("Authorization", request.Token)
	} else if request.Username != "" && request.Password != "" {
		// Basic Authentication
		auth := request.Username + ":" + request.Password
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	}

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
	// SCAFFOLDING #17-1 - pkg/adapter/datasource.go: To add support for multiple entities that require different parsing functions
	// Add code to call different ParseResponse functions for each entity response.
	objects, nextCursor, parseErr := ParseResponse(body)
	if parseErr != nil {
		return nil, parseErr
	}

	response.Objects = objects
	response.NextCursor = nextCursor

	return response, nil
}

func ParseResponse(body []byte) (objects []map[string]any, nextCursor string, err *framework.Error) {
	var data *DatasourceResponse

	// Unmarshal the response from the datasource.
	// UnmarshalJSON is implemented to handle the response from the datasource.
	unmarshalErr := json.Unmarshal(body, &data)
	if unmarshalErr != nil {
		return nil, "", &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// SCAFFOLDING #18 - pkg/adapter/datasource.go: Add response validations.
	// Add necessary validations to check if the response from the datasource is what is expected.

	// SCAFFOLDING #19 - pkg/adapter/datasource.go: Populate next page information (called cursor in SGNL adapters).
	// Populate nextCursor with the cursor returned from the datasource, if present.
	nextCursor = ""
	if data.More {
		// If there are more pages, next cursor is the offset + limit.
		nextCursor = strconv.Itoa(data.Offset + data.Limit)
	}

	return data.Objects, nextCursor, nil
}

// Custom Unmarshal implementation to handle the response from the datasource
func (d *DatasourceResponse) UnmarshalJSON(data []byte) error {

	// A generic map is used to unmarshal the response first, then the objects are extracted from the map.
	var raw map[string]json.RawMessage

	// Unmarshal into a generic map first
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Check if any of the valid entity external IDs are present in the response
	// Supports unmarshal of different entities in the response.
	// Add more entities as needed in ValidEntityExternalIDs map.
	found := false
	for key := range ValidEntityExternalIDs {
		if value, exists := raw[key]; exists {
			var objects []map[string]any
			if err := json.Unmarshal(value, &objects); err == nil {
				d.Objects = objects
				found = true
				break
			}
		}
	}

	// Check if pagination info is present in the response
	if value, exists := raw["offset"]; exists {
		if err := json.Unmarshal(value, &d.Offset); err != nil {
			return err
		}
	}
	if value, exists := raw["limit"]; exists {
		if err := json.Unmarshal(value, &d.Limit); err != nil {
			return err
		}
	}
	if value, exists := raw["more"]; exists {
		if err := json.Unmarshal(value, &d.More); err != nil {
			return err
		}
	}

	if !found {
		return fmt.Errorf("no valid objects found in JSON")
	}
	return nil
}

func addQueryParams(baseUrl *url.URL, request *Request) {

	query := baseUrl.Query()
	if request.PageSize > 0 {
		query.Add("limit", fmt.Sprintf("%d", request.PageSize))
	}
	if request.Cursor != "" {
		query.Add("offset", request.Cursor)
	}
	if request.Total {
		query.Add("total", "true")
	}
	if request.Query != "" {
		query.Add("query", request.Query)
	}

	baseUrl.RawQuery = query.Encode()

}
