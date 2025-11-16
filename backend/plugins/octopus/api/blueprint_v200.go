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
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/octopus/models"
)

func MakePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connectionId uint64,
	bpScopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	// load connection, scope and scopeConfig from the db
	connection, err := dsHelper.ConnSrv.FindByPk(connectionId)
	if err != nil {
		return nil, nil, err
	}
	scopeDetails, err := dsHelper.ScopeApi.MapScopeDetails(connectionId, bpScopes)
	if err != nil {
		return nil, nil, err
	}

	sc, err := makeScopeV200(connectionId, scopeDetails)
	if err != nil {
		return nil, nil, err
	}

	pp, err := makePipelinePlanV200(subtaskMetas, connection, scopeDetails)
	if err != nil {
		return nil, nil, err
	}

	return pp, sc, nil
}

func makeScopeV200(
	connectionId uint64,
	scopeDetails []*srvhelper.ScopeDetail[models.OctopusProject, models.OctopusScopeConfig],
) ([]plugin.Scope, errors.Error) {
	sc := make([]plugin.Scope, 0, 3*len(scopeDetails))

	for _, scope := range scopeDetails {
		octopusProject, scopeConfig := scope.Scope, scope.ScopeConfig
		id := didgen.NewDomainIdGenerator(&models.OctopusProject{}).Generate(connectionId, octopusProject.Id)

		// add cicd_scope to scopes
		if utils.StringsContains(scopeConfig.Entities, plugin.DOMAIN_TYPE_CICD) {
			scopeCICD := devops.NewCicdScope(id, octopusProject.Slug)
			sc = append(sc, scopeCICD)
		}
	}

	return sc, nil
}

func makePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connection *models.OctopusConnection,
	scopeDetails []*srvhelper.ScopeDetail[models.OctopusProject, models.OctopusScopeConfig],
) (coreModels.PipelinePlan, errors.Error) {
	plans := make(coreModels.PipelinePlan, 0, len(scopeDetails))
	for i, scope := range scopeDetails {
		octopusProject, scopeConfig := scope.Scope, scope.ScopeConfig

		stage := plans[i]
		if stage == nil {
			stage = coreModels.PipelineStage{}
		}

		// construct subtasks
		task, err := helper.MakePipelinePlanTask(
			pluginName,
			subtaskMetas,
			scopeConfig.Entities,
			map[string]interface{}{
				"connectionId": connection.ID,
				"projectId":    octopusProject.Id,
				"Name":         octopusProject.Name,
			})

		if err != nil {
			return nil, err
		}
		stage = append(stage, task)
		plans[i] = stage
	}

	return plans, nil
}
