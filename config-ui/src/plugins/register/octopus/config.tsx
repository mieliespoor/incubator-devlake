/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import { ExternalLink } from '@/components';
import { DOC_URL } from '@/release';
import { IPluginConfig } from '@/types';

import Icon from './assets/icon.svg?react';

export const OctopusDeployConfig: IPluginConfig = {
  plugin: 'octopus',
  name: 'Octopus Deploy',
  icon: ({ color }) => <Icon fill={color} />,
  sort: 10,
  isBeta: true,
  connection: {
    docLink: '',
    initialValues: {
      endpoint: 'https://',
    },
    fields: [
      'name',
      {
        key: 'endpoint',
        subLabel: 'Provide the Octopus Deploy server API endpoint. E.g. https://octopus.example.com/api/',
      },
      {
        key: 'token',
        label: 'X-Octopus-ApiKey',
        subLabel: (
          <>
            Provide your Octopus Deploy API key for authentication.{' '}
            <ExternalLink link="https://octopus.com/docs/octopus-rest-api/how-to-create-an-api-key">
              Learn how to generate an api key
            </ExternalLink>
          </>
        ),
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses 3,000 requests/hour for data collection for Octopus Deploy. You can adjust the collection speed by setting your desired rate limit.',
        learnMore: '',
        externalInfo: 'Octopus Deploy does not specify a maximum rate limit value.',
        defaultValue: 3000,
      },
    ],
  },
  dataScope: {
    title: 'Applications',
    millerColumn: {
      columnCount: 2,
      firstColumnTitle: 'Projects',
    },
  },
  scopeConfig: {
    entities: ['CICD'],
    transformation: {
      component: 'OctopusDeployTransformation',
      envNamePattern: '(?i)prod(.*)',
      deploymentPattern: '',
      productionPattern: '',
    },
  },
};
