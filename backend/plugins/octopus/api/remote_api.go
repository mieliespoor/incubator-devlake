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
	"fmt"
	"net/url"
	"strconv"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/octopus/models"
)

type OctopusRemotePagination struct {
	Skip int `json:"skip"`
	Take int `json:"take"`
}

type OctopusProjectResourceCollection struct {
	Items               []models.OctopusProject
	ItemsPerPage        int               `json:"ItemsPerPage"`
	NumberOfPages       int               `json:"NumberOfPages"`
	IsStaleResultsCache bool              `json:"IsStaleResultsCache"`
	TotalResults        int               `json:"TotalResults"`
	PageNumber          int               `json:"PageNumber"`
	LastPageNumber      int               `json:"LastPageNumber"`
	Links               map[string]string `json:"Links"`
}

func queryOctopusProjects(
	apiClient plugin.ApiClient,
	keyword string,
	page OctopusRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.OctopusProject],
	nextPage *OctopusRemotePagination,
	err errors.Error,
) {

	var allProjects []models.OctopusProject
	pge := 0

	for {
		projects, hasmore, err := getProjectsPerPage(apiClient, "", pge)
		if err != nil {
			return nil, nil, err
		}
		allProjects = append(allProjects, projects...)
		if !hasmore {
			break
		}
		pge++
	}

	children = make([]dsmodels.DsRemoteApiScopeListEntry[models.OctopusProject], 0)
	for _, project := range allProjects {
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.OctopusProject]{
			Scope: project,
			Id:    project.Id,
			Name:  project.Name,
		})
	}

	return
}

func getProjectsPerPage(apiClient plugin.ApiClient, spaceId string, page int) ([]models.OctopusProject, bool, errors.Error) {
	params := url.Values{}
	params.Add("skip", fmt.Sprintf("%d", strconv.Itoa(page*100)))
	params.Add("take", "100")

	endpoint := fmt.Sprintf("/api/spaces/%s/projects", spaceId)

	res, err := apiClient.Get(endpoint, params, nil)
	if err != nil {
		return nil, false, err
	}

	resBody := OctopusProjectResourceCollection{}

	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return nil, false, err
	}

	hasmore := resBody.PageNumber < resBody.NumberOfPages-1

	return resBody.Items, hasmore, nil
}

func listOctopusRemoteScopes(
	connection *models.OctopusConnection,
	apiClient plugin.ApiClient,
	spaceId string,
	page OctopusRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.OctopusProject],
	nextPage *OctopusRemotePagination,
	err errors.Error,
) {
	if page.Take == 0 {
		page.Take = 100
	}
	if page.Skip < 0 {
		page.Skip = 0
	}
	options := map[string]string{
		"skip": fmt.Sprintf("%d", page.Skip),
		"take": fmt.Sprintf("%d", page.Take),
	}
	if spaceId != "" {
		options["spaceId"] = spaceId
	}
	var projects []models.OctopusProject
	// res, err := apiClient.Get("projects/search", url.Values{
	// 	"skip":  {fmt.Sprintf("%v", page.Skip)},
	// 	"take": {fmt.Sprintf("%v", page.Take)},
	// 	"partialName":  {keyword},
	// }, nil)
	apiClient.Get("projects", nil, nil)
	err = apiClient.GetPaged("projects", options, &projects)
	if err != nil {
		return nil, nil, err
	}
	for _, project := range projects {
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.OctopusProject]{
			// Scope: project,
			Type: api.RAS_ENTRY_TYPE_SCOPE,
			Id:   project.Id,
			Name: project.Name,
			Data: &project,
		})
	}

	// for _, project := range resBody.Components {
	// 	children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.SonarqubeProject]{
	// 		Type:     api.RAS_ENTRY_TYPE_SCOPE,
	// 		Id:       fmt.Sprintf("%v", project.ProjectKey),
	// 		ParentId: nil,
	// 		Name:     project.Name,
	// 		FullName: project.Name,
	// 		Data:     project.ConvertApiScope(),
	// 	})
	// }

	// if resBody.Paging.Total > resBody.Paging.PageIndex*resBody.Paging.PageSize {
	// 	nextPage = &SonarqubeRemotePagination{
	// 		Page:     resBody.Paging.PageIndex + 1,
	// 		PageSize: resBody.Paging.PageSize,
	// 	}
	// }

	if len(projects) == page.Take {
		nextPage = &OctopusRemotePagination{
			Skip: page.Skip + page.Take,
			Take: page.Take,
		}
	}
	return children, nextPage, nil
}

// RemoteScopes list all available scopes on the remote server
// @Summary list all available scopes on the remote server
// @Description list all available scopes on the remote server
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.JenkinsJob]
// @Tags plugins/jenkins
// @Router /plugins/jenkins/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// @Summary Remote server API proxy
// @Description Forward API requests to the specified remote server
// @Param connectionId path int true "connection ID"
// @Param path path string true "path to a API endpoint"
// @Router /plugins/jenkins/connections/{connectionId}/proxy/{path} [GET]
// @Tags plugins/jenkins
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
