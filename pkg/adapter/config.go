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
	"errors"
)

var (
	ErrConfigIsNil                    = errors.New("request contains no config")
	ErrConfigDatasourceVersionIsEmpty = errors.New("config datasourceVersion is not set")
)

// Config is the optional configuration passed in each GetPage calls to the
// adapter.
type Config struct {
	// SCAFFOLDING:
	// Add/remove fields as needed.
	// Every field MUST have a `json` tag.

	// Example config field.
	DatasourceVersion string `json:"datasourceVersion,omitempty"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	// SCAFFOLDING:
	// Update the checks below to validate the fields in Config.

	switch {
	case c == nil:
		return ErrConfigIsNil
	case c.DatasourceVersion == "":
		return ErrConfigDatasourceVersionIsEmpty
	default:
		return nil
	}
}
