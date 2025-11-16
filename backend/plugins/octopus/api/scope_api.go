/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

// PutScopes create or update Octopus applications
// @Summary create or update octopus applications
// @Description Create or update Octopus application scopes
// @Tags plugins/octopus
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scope body models.OctopusApplication true "json"
// @Success 200  {object} []models.OctopusApplication
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/octopus/connections/{connectionId}/scopes [PUT]
func PutScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.PutMultiple(input)
}

// UpdateScope patch to octopus application
// @Summary patch to octopus application
// @Description Patch Octopus application scope
// @Tags plugins/octopus
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scopeId path string false "application name"
// @Param scope body models.OctopusApplication true "json"
// @Success 200  {object} models.OctopusApplication
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/octopus/connections/{connectionId}/scopes/{scopeId} [PATCH]
func PatchScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Patch(input)
}

// GetScopeList get Octopus applications
// @Summary get Octopus applications
// @Description Get Octopus applications
// @Tags plugins/octopus
// @Param connectionId path int false "connection ID"
// @Param searchTerm query string false "search term for scope name"
// @Param blueprints query bool false "also return blueprints using these scopes as part of the payload"
// @Success 200  {object} []models.OctopusApplication
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/octopus/connections/{connectionId}/scopes [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetPage(input)
}

// GetScope get one Octopus application
// @Summary get one Octopus application
// @Description Get one Octopus application
// @Tags plugins/octopus
// @Param connectionId path int false "connection ID"
// @Param scopeId path string false "application name"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page number, default 1"
// @Success 200  {object} models.OctopusApplication
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/octopus/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetScopeDetail(input)
}

// DeleteScope delete plugin data associated with the scope and optionally the scope itself
// @Summary delete plugin data associated with the scope and optionally the scope itself
// @Description Delete data associated with Octopus application scope
// @Tags plugins/octopus
// @Param connectionId path int true "connection ID"
// @Param scopeId path string true "application name"
// @Param delete_data_only query bool false "Only delete the scope data, not the scope itself"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/octopus/connections/{connectionId}/scopes/{scopeId} [DELETE]
func DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Delete(input)
}
