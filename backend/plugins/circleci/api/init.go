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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/circleci/models"
	"github.com/go-playground/validator/v10"
)

var vld *validator.Validate
var connectionHelper *helper.ConnectionApiHelper
var scopeHelper *helper.ScopeApiHelper[models.CircleciConnection, models.CircleciProject, models.CircleciScopeConfig]
var scHelper *helper.ScopeConfigHelper[models.CircleciScopeConfig, *models.CircleciScopeConfig]
var remoteHelper *helper.RemoteApiHelper[models.CircleciConnection, models.CircleciProject, RemoteProject, helper.NoRemoteGroupResponse]
var basicRes context.BasicRes

func Init(br context.BasicRes, p plugin.PluginMeta) {
	basicRes = br
	vld = validator.New()
	connectionHelper = helper.NewConnectionHelper(
		basicRes,
		vld,
		p.Name(),
	)
	params := &helper.ReflectionParameters{
		ScopeIdFieldName:     "Id",
		ScopeIdColumnName:    "id",
		RawScopeParamName:    "Slug",
		SearchScopeParamName: "name",
	}
	scopeHelper = helper.NewScopeHelper[models.CircleciConnection, models.CircleciProject, models.CircleciScopeConfig](
		basicRes,
		vld,
		connectionHelper,
		helper.NewScopeDatabaseHelperImpl[models.CircleciConnection, models.CircleciProject, models.CircleciScopeConfig](
			basicRes, connectionHelper, params),
		params,
		nil,
	)
	remoteHelper = helper.NewRemoteHelper[models.CircleciConnection, models.CircleciProject, RemoteProject, helper.NoRemoteGroupResponse](
		basicRes,
		vld,
		connectionHelper,
	)
	scHelper = helper.NewScopeConfigHelper[models.CircleciScopeConfig, *models.CircleciScopeConfig](
		basicRes,
		vld,
		p.Name(),
	)
}
