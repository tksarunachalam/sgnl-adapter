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
	"fmt"
	"strconv"
	"time"
)

const (
	// ObjectsPerEntity is the number of JSON objects created for each entity.
	ObjectsPerEntity = 1000
)

var (
	// Data is the example datasource's JSON data.
	// Each key is an entity external ID and its associated value is that
	// entity's list of JSON objects.
	//
	// This datasource contains two entities, with external IDs:
	//  - User
	//  - Group
	Data map[string][]map[string]any
)

func init() {
	Data = make(map[string][]map[string]any, 2)

	now := time.Now().UTC().Format(time.RFC3339)

	userData := make([]map[string]any, 0, ObjectsPerEntity)
	for i := 1; i <= ObjectsPerEntity; i++ {
		user := map[string]any{
			"id":          strconv.Itoa(i),
			"displayName": fmt.Sprintf("User #%d", i),
			"email":       fmt.Sprintf("user%d@example.com", i),
			"createdAt":   now,
		}
		userData = append(userData, user)
	}
	Data["User"] = userData

	groupData := make([]map[string]any, 0, ObjectsPerEntity)
	for i := 1; i <= ObjectsPerEntity; i++ {
		group := map[string]any{
			"id":          strconv.Itoa(i),
			"displayName": fmt.Sprintf("Group #%d", i),
			"createdAt":   now,
		}
		groupData = append(groupData, group)
	}
	Data["Group"] = groupData
}
