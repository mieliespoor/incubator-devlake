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

import { useMemo, useState } from 'react';
import { Button, Intent } from '@blueprintjs/core';
import { pick } from 'lodash';

import API from '@/api';
import { useAppDispatch, useAppSelector } from '@/app/hook';
import { ExternalLink, Buttons } from '@/components';
import { selectConnection } from '@/features/connections';
import { PluginConfig, PluginConfigType } from '@/plugins';
import { operator } from '@/utils';

import { Form } from './fields';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId?: ID;
  onSuccess?: (id: ID) => void;
}

export const ConnectionForm = ({ plugin, connectionId, onSuccess }: Props) => {
  const [values, setValues] = useState<any>({});
  const [errors, setErrors] = useState<Record<string, any>>({});
  const [operating, setOperating] = useState(false);

  const dispatch = useAppDispatch();
  const connection = useAppSelector((state) => selectConnection(state, `${plugin}-${connectionId}`));

  const {
    name,
    connection: { docLink, fields, initialValues },
  } = useMemo(() => PluginConfig.find((p) => p.plugin === plugin) as PluginConfigType, [plugin]);

  const disabled = useMemo(() => {
    return Object.values(errors).some((value) => value);
  }, [errors]);

  const handleTest = async () => {
    await operator(
      () =>
        API.connection.test(
          plugin,
          pick(values, [
            'endpoint',
            'token',
            'username',
            'password',
            'proxy',
            'authMethod',
            'appId',
            'secretKey',
            'tenantId',
            'tenantType',
            'dbUrl',
          ]),
        ),
      {
        setOperating,
        formatMessage: () => 'Test Connection Successfully.',
      },
    );
  };

  const handleSave = async () => {
    const [success, res] = await operator(
      () =>
        !connectionId ? API.connection.create(plugin, values) : API.connection.update(plugin, connectionId, values),
      {
        setOperating,
        formatMessage: () => (!connectionId ? 'Create a New Connection Successful.' : 'Update Connection Successful.'),
      },
    );

    if (success) {
      onSuccess?.(res.id);
    }
  };

  return (
    <S.Wrapper>
      <S.Tips>
        If you run into any problems while creating a new connection for {name},{' '}
        <ExternalLink link={docLink}>check out this doc</ExternalLink>.
      </S.Tips>
      <S.Form>
        <Form
          name={name}
          fields={fields}
          initialValues={{ ...initialValues, ...(connection ?? {}) }}
          values={values}
          errors={errors}
          setValues={setValues}
          setErrors={setErrors}
        />
        <Buttons position="bottom" align="right">
          <Button loading={operating} disabled={disabled} outlined text="Test Connection" onClick={handleTest} />
          <Button
            loading={operating}
            disabled={disabled}
            intent={Intent.PRIMARY}
            outlined
            text="Save Connection"
            onClick={handleSave}
          />
        </Buttons>
      </S.Form>
    </S.Wrapper>
  );
};
