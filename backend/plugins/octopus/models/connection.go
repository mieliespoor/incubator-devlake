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

package models

import (
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

// OctopusConn holds the essential information to connect to the Octopus API
type OctopusConn struct {
	api.RestConnection `mapstructure:",squash"`
	api.AccessToken    `mapstructure:",squash"`
}

var _ plugin.ApiConnection = (*OctopusConnection)(nil)

func (conn *OctopusConn) Sanitize() OctopusConn {
	conn.Token = utils.SanitizeString(conn.Token)
	return *conn
}

type OctopusConnection struct {
	api.BaseConnection `mapstructure:",squash"`
	OctopusConn        `mapstructure:",squash"`
}

func (connection OctopusConnection) Sanitize() OctopusConnection {
	connection.OctopusConn = connection.OctopusConn.Sanitize()
	return connection
}

// TableName returns the table name for OctopusConnection
func (OctopusConnection) TableName() string {
	return "_tool_octopus_connections"
}
