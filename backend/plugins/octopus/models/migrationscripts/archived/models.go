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

package archived

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type OctopusConnection struct {
	Name                    string `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	archived.RestConnection `mapstructure:",squash"`
	archived.AccessToken    `mapstructure:",squash"`
	archived.Model
}

// TableName returns the table name for OctopusConnection
func (OctopusConnection) TableName() string {
	return "_tool_octopus_connections"
}

type OctopusProject struct {
	common.Scope   `mapstructure:",squash" gorm:"embedded"`
	Id             string    `gorm:"primaryKey;type:VARCHAR(255) NOT NULL"`
	ConnectionId   uint64    `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Name           string    `gorm:"type:VARCHAR(255)"`
	Description    string    `gorm:"type:TEXT"`
	LastModifiedBy string    `gorm:"type:VARCHAR(255)"`
	LastModifiedOn time.Time `gorm:"type:TIMESTAMP;DEFAULT:CURRENT_TIMESTAMP"`
	Slug           string    `gorm:"type:VARCHAR(255)"`
	SpaceId        string    `gorm:"type:VARCHAR(255)"`
	archived.NoPKModel
}

func (OctopusProject) TableName() string {
	return "_tool_octopus_projects"
}

type OctopusScopeConfig struct {
	archived.ScopeConfig `mapstructure:",squash" json:",inline" gorm:"embedded"`
	Name                 string   `gorm:"type:varchar(255);index"`
	Entities             []string `gorm:"type:json;serializer:json"`
	ID                   uint64   `gorm:"primaryKey;autoIncrement"`
	ConnectionId         uint64   `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ProjectId            string   `gorm:"primaryKey;type:VARCHAR(255) NOT NULL"`
	EnvNamePattern       string   `gorm:"type:varchar(255)"`
	// add more fields if necessary
}

func (OctopusScopeConfig) TableName() string {
	return "_tool_octopus_scope_configs"
}
