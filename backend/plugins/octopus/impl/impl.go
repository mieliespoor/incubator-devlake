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

package impl

import (
	"fmt"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/subtaskmeta/sorter"
	"github.com/apache/incubator-devlake/plugins/octopus/api"
	"github.com/apache/incubator-devlake/plugins/octopus/models"
	"github.com/apache/incubator-devlake/plugins/octopus/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/octopus/tasks"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginModel
	plugin.PluginMigration
	plugin.PluginSource
	plugin.DataSourcePluginBlueprintV200
	plugin.CloseablePluginTask
} = (*Octopus)(nil)

type Octopus struct{}

var sortedSubtaskMetas []plugin.SubTaskMeta

func init() {
	var err error
	// check subtask meta loop and gen subtask list when plugin init
	sortedSubtaskMetas, err = sorter.NewDependencySorter(tasks.SubTaskMetaList).Sort()
	if err != nil {
		panic(err)
	}
}

func (p Octopus) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)
	return nil
}

func (p Octopus) Connection() dal.Tabler {
	return &models.OctopusConnection{}
}

func (p Octopus) Scope() plugin.ToolLayerScope {
	return &models.OctopusProject{}
}

func (p Octopus) ScopeConfig() dal.Tabler {
	return &models.OctopusScopeConfig{}
}

func (p Octopus) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes)
}

func (p Octopus) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.OctopusConnection{},
		&models.OctopusDeployment{},
		&models.OctopusEnvironment{},
		&models.OctopusProject{},
		&models.OctopusScopeConfig{},
	}
}

func (p Octopus) Description() string {
	return "To collect data from Octopus Deploy"
}

func (p Octopus) Name() string {
	return "octopus"
}

func (p Octopus) SubTaskMetas() []plugin.SubTaskMeta {
	return sortedSubtaskMetas
}

func (p Octopus) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	logger := taskCtx.GetLogger()
	logger.Debug("%v", options)
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("connectionId is invalid")
	}
	connection := &models.OctopusConnection{}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		p.Name(),
	)
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "connection not found")
	}

	apiClient, err := tasks.NewOctopusApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}

	// Add logic to enrich options and fetch scope config if needed

	taskData := tasks.OctopusTaskData{
		Options:   op,
		ApiClient: apiClient,
	}

	return &taskData, nil
}

func (p Octopus) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/octopus"
}

func (p Octopus) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Octopus) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
			"GET":    api.GetConnection,
		},
		"connections/:connectionId/test": {
			"POST": api.TestExistingConnection,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":    api.GetScope,
			"PATCH":  api.PatchScope,
			"DELETE": api.DeleteScope,
		},
		"connections/:connectionId/scopes/:scopeId/latest-sync-state": {
			"GET": api.GetScopeLatestSyncState,
		},
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/search-remote-scopes": {
			"GET": api.SearchRemoteScopes,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScopes,
		},
		"connections/:connectionId/scope-configs": {
			"POST": api.CreateScopeConfig,
			"GET":  api.GetScopeConfigList,
		},
		"connections/:connectionId/scope-configs/:scopeConfigId": {
			"PATCH":  api.PatchScopeConfig,
			"GET":    api.GetScopeConfig,
			"DELETE": api.DeleteScopeConfig,
		},
		"connections/:connectionId/proxy/rest/*path": {
			"GET": api.Proxy,
		},
		"scope-config/:scopeConfigId/projects": {
			"GET": api.GetProjectsByScopeConfig,
		},
	}
}

func (p Octopus) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.OctopusTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
