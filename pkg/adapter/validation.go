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
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

const (
	ErrMsgInvalidAddress            = "Provided datasource address is not an https:// URL"
	ErrMsgInvalidAuth               = "Provided datasource auth is missing required basic credentials"
	ErrMsgInvalidEntityExternalId   = "Provided entity external ID is invalid"
	ErrMsgMissingUniqueIdAttribute  = "Requested entity attributes are missing unique ID attribute"
	ErrMsgChildEntitiesNotSupported = "Requested entity does not support child entities"
	ErrMsgInvalidOrderedFalse       = "Ordered must be true"
	ErrMsgInvalidPageSizeFmt        = "Provided page size (%d) exceeds maximum (%d)"
)

const (
	// MaxPageSize is the maximum page size allowed in a GetPage request.
	//
	// SCAFFOLDING:
	// Update this limit to match the limit of the datasource.
	MaxPageSize = 100

	// UniqueIdAttribute is the name of the attribute containing the unique ID of
	// each returned object for the requested entity.
	//
	// SCAFFOLDING:
	// Update this to match the name of the unique ID attribute in the
	// requested entity.
	UniqueIdAttribute = "id"
)

var (
	// ValidEntityExternalIds is the set of valid external IDs of entities that
	// can be queried.
	//
	// SCAFFOLDING:
	// Update this set to match the set of entities that can be queried from
	// the datasource.
	ValidEntityExternalIds = map[string]struct{}{
		"User":  {},
		"Group": {},
	}
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if err := request.Config.Validate(ctx); err != nil {
		return &framework.Error{
			Message: err.Error(),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// SCAFFOLDING:
	// Modify this validation to match the format of the datasource's
	// address.
	if !strings.HasPrefix(request.Address, "https://") {
		return &framework.Error{
			Message: ErrMsgInvalidAddress,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// SCAFFOLDING:
	// Modify this validation to match the authn mechanism(s) supported by the
	// datasource.
	if request.Auth == nil || request.Auth.Basic == nil {
		return &framework.Error{
			Message: ErrMsgInvalidAuth,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if _, found := ValidEntityExternalIds[request.Entity.ExternalId]; !found {
		return &framework.Error{
			Message: ErrMsgInvalidEntityExternalId,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Validate that at least the unique ID attribute for the requested entity
	// is requested.
	var uniqueIdAttributeFound bool
	for _, attribute := range request.Entity.Attributes {
		if attribute.ExternalId == UniqueIdAttribute {
			uniqueIdAttributeFound = true
			break
		}
	}

	if !uniqueIdAttributeFound {
		return &framework.Error{
			Message: ErrMsgMissingUniqueIdAttribute,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Validate that no child entities are requested.
	//
	// SCAFFOLDING:
	// Modify this validation if the entity contains child entities.
	if len(request.Entity.ChildEntities) > 0 {
		return &framework.Error{
			Message: ErrMsgChildEntitiesNotSupported,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// SCAFFOLDING:
	// If the datasource doesn't support sorting results by unique ID
	// attribute for the requested entity, check instead that Ordered is set to
	// false.
	if !request.Ordered {
		return &framework.Error{
			Message: ErrMsgInvalidOrderedFalse,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if request.PageSize > MaxPageSize {
		return &framework.Error{
			Message: fmt.Sprintf(ErrMsgInvalidPageSizeFmt, request.PageSize, MaxPageSize),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	return nil
}
