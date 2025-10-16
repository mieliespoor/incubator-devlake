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

import "fmt"

type OctopusDeployment struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ConnectionId  uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL" json:"connectionId"`
	DeploymentId  string `gorm:"primaryKey;type:VARCHAR(255) NOT NULL" json:"deploymentId"`
	ProjectId     string `gorm:"type:VARCHAR(255)" json:"projectId"`
	EnvironmentId string `gorm:"type:VARCHAR(255)" json:"environmentId"`
	ReleaseId     string `gorm:"type:VARCHAR(255)" json:"releaseId"`
	DeployedAt    string `gorm:"type:VARCHAR(255)" json:"deployedAt"`
	DeployedBy    string `gorm:"type:VARCHAR(255)" json:"deployedBy"`
	TenantId      string `gorm:"type:VARCHAR(255)" json:"tenantId"`
	TenantName    string `gorm:"type:VARCHAR(255)" json:"tenantName"`
	// add more fields if necessary
}

func (OctopusDeployment) TableName() string {
	return "octopus_deployments"
}

func (d OctopusDeployment) GetDeploymentKey() string {
	return fmt.Sprintf("%d-%s", d.ConnectionId, d.DeploymentId)
}
