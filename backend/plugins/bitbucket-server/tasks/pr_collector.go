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

package tasks

import (
	"github.com/apache/incubator-devlake/core/errors"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_PULL_REQUEST_TABLE = "bitbucket_server_api_pull_requests"

// this struct should be moved to `bitbucket_api_common.go`

var CollectApiPullRequestsMeta = plugin.SubTaskMeta{
	Name:             "collectApiPullRequests",
	EntryPoint:       CollectApiPullRequests,
	EnabledByDefault: true,
	Required:         false,
	Description:      "Collect PullRequests data from Bitbucket Server api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

func CollectApiPullRequests(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_TABLE)
	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    50,
		UrlTemplate: "projects/{{ .Params.FullName }}/pull-requests",
		Query: GetQueryCreatedAndUpdated(
			`values.id,values.state,values.title,values.description,`+
				`values.fromRef.latestCommit,values.links.self.href,values.author.user.name,values.createdDate,values.updatedDate,`+
				`values.toRef.displayId,values.toRef.latestCommit,values.toRef.repository.slug,`+
				`values.fromRef.repository.slug,`+
				`start,limit,size`,
			collectorWithState),
		GetTotalPages:  GetTotalPagesFromResponse,
		ResponseParser: GetRawMessageFromResponse,
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}
