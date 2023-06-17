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

package example_datasource

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

// Datasource directly implements the Client interface to allow querying
// a static in-memory set of JSON objects.
type Datasource struct {
}

// NewClient returns a Client that directly queries the in-memory example
// datasource.
func NewClient() Client {
	return &Datasource{}
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, error) {
	// Parse the entity's external ID from the URL.
	index := strings.LastIndex(request.URL, "/")
	entityExternalId := request.URL[index:]

	entityData, found := Data[entityExternalId]
	if !found {
		return &Response{
			StatusCode: http.StatusNotFound,
		}, nil
	}

	// The cursor is the index of the next object to return.
	startIndex, err := strconv.Atoi(request.Cursor)
	if err != nil {
		return &Response{
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	response := &Response{
		StatusCode: http.StatusOK,
	}

	if startIndex < len(entityData) {
		endIndex := startIndex + int(request.PageSize)
		if endIndex > len(entityData) {
			endIndex = len(entityData)
		}

		response.Objects = entityData[startIndex:endIndex]

		if endIndex < len(entityData) {
			response.NextCursor = strconv.Itoa(endIndex)
		}
	}

	return response, nil
}
