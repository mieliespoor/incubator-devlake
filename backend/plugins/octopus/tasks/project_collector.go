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
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/octopus/models"
)

const RAW_PROJECTS_TABLE = "octopus_api_projects"

var CollectApiProjectsMeta = plugin.SubTaskMeta{
	Name:             "Collect Projects",
	EntryPoint:       CollectApiProjects,
	EnabledByDefault: true,
	Description:      "Collect project data from gitlab api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{},
}

func CollectApiProjects(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*OctopusTaskData)

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: RAW_PROJECTS_TABLE,
			Params: models.OctopusApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.ApplicationName,
			},
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: "/projects/{{.Params.Name}}",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("referesh", "false")
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data []json.RawMessage
			err := api.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}
			return data, nil
		},
	})
	if err != nil {
		return err
	}
	return collector.Execute()
}
