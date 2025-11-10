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
	"context"
	"fmt"
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/octopus/models"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/mitchellh/mapstructure"
)

type OctopusTestConnResponse struct {
	shared.ApiBody
	Connection *models.OctopusConn
}

func testConnection(ctx context.Context, connection models.OctopusConn) (*plugin.ApiResourceOutput, errors.Error) {
	// validate
	if vld != nil {
		e := vld.StructExcept(connection, "BasicAuth", "AccessToken")
		if e != nil {
			return nil, errors.Convert(e)
		}
	}
	// test connection
	apiClient, err := api.NewApiClientFromConnection(ctx, basicRes, &connection)
	if err != nil {
		return nil, err
	}
	// serverInfo checking
	res, err := apiClient.Get("serverstatus/system-info", nil, nil)
	if err != nil {
		return nil, err
	}

	body := OctopusTestConnResponse{}
	body.Success = false
	body.Message = "failed"
	switch res.StatusCode {
	case 200: // right StatusCode
		valid := &models.OctopusSystemInfo{}
		err = api.UnmarshalResponse(res, valid)
		if err != nil {
			return nil, err
		}
		body.Success = true
		body.Message = "success"
		connection = connection.Sanitize()
		body.Connection = &connection

		return &plugin.ApiResourceOutput{Body: body, Status: 200}, nil
	case 401: // error secretKey or nonceStr
		return &plugin.ApiResourceOutput{Body: false, Status: http.StatusUnauthorized}, nil
	default: // unknown what happen , back to user
		return &plugin.ApiResourceOutput{Body: res.Body, Status: res.StatusCode}, nil
	}
}

// TestConnection test Octopus Deploy connection
// @Summary test Octopus Deploy connection
// @Description Test Octopus Deploy Connection
// @Tags plugins/octopus
// @Param body body models.OctopusConnection true "json body"
// @Success 200  {object} OctopusTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/octopus/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.OctopusConn
	e := mapstructure.Decode(input.Body, &connection)
	if e != nil {
		return nil, errors.Convert(e)
	}
	// test connection
	result, err := testConnection(context.TODO(), connection)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// TestExistingConnection test octopus deploy connection options
// @Summary test octopus deploy connection
// @Description Test octopus deploy Connection
// @Tags plugins/octopus
// @Param connectionId path int true "connection ID"
// @Success 200  {object} OctopusTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/octopus/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return nil, errors.Convert(err)
	}
	// test connection
	testConnectionResult, testConnectionErr := testConnection(context.TODO(), connection.OctopusConn)
	if testConnectionErr != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, testConnectionErr)
	}
	if testConnectionResult.Status != http.StatusOK {
		errMsg := fmt.Sprintf("Test connection fail, unexpected status code: %d", testConnectionResult.Status)
		return nil, plugin.WrapTestConnectionErrResp(basicRes, errors.Default.New(errMsg))
	}
	return testConnectionResult, nil
}

// PostConnections create octopus deploy connection
// @Summary create octopus deploy connection
// @Description Create octopus deploy connection
// @Tags plugins/octopus
// @Param body body models.OctopusConnection true "json body"
// @Success 200  {object} models.OctopusConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/octopus/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Post(input)
}

// PatchConnection patch octopus deploy connection
// @Summary patch octopus deploy connection
// @Description Patch octopus deploy connection
// @Tags plugins/octopus
// @Param body body models.OctopusConnection true "json body"
// @Param connectionId path int false "connection ID"
// @Success 200  {object} models.OctopusConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/octopus/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Patch(input)
}

// DeleteConnection delete a octopus connection
// @Summary delete a octopus connection
// @Description Delete a octopus connection
// @Tags plugins/octopus
// @Param connectionId path int false "connection ID"
// @Success 200  {object} models.OctopusConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} srvhelper.DsRefs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/octopus/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Delete(input)
}

// ListConnections get all octopus connections
// @Summary get all octopus connections
// @Description Get all octopus connections
// @Tags plugins/octopus
// @Success 200  {object} []models.OctopusConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/octopus/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// GetConnection get octopus connection detail
// @Summary get octopus connection detail
// @Description Get octopus connection detail
// @Tags plugins/octopus
// @Param connectionId path int false "connection ID"
// @Success 200  {object} models.OctopusConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/octopus/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}
