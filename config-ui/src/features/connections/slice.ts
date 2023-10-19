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

import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { flatten } from 'lodash';

import API from '@/api';
import type { ConnectionForm } from '@/api/connection/types';
import { RootState } from '@/app/store';
import { PluginConfig } from '@/plugins';
import { IConnection, IConnectionStatus } from '@/types';

import { transformConnection } from './utils';

const initialState: {
  connections: IConnection[];
} = {
  connections: [],
};

export const init = createAsyncThunk('connections/init', async () => {
  const res = await Promise.all(
    PluginConfig.map(async ({ plugin }) => {
      const connections = await API.connection.list(plugin);
      return connections.map((connection) => transformConnection(plugin, connection));
    }),
  );
  return flatten(res);
});

export const fetchConnections = createAsyncThunk('connections/fetchConnections', async (plugin: string) => {
  const connections = await API.connection.list(plugin);
  return {
    plugin,
    connections: connections.map((connection) => transformConnection(plugin, connection)),
  };
});

export const testConnection = createAsyncThunk(
  'connections/testConnection',
  async ({ unique, plugin, endpoint, proxy, token, username, password, authMethod, secretKey, appId }: IConnection) => {
    const res = await API.connection.test(plugin, {
      endpoint,
      proxy,
      token,
      username,
      password,
      authMethod,
      secretKey,
      appId,
    });

    return {
      unique,
      status: res.success ? IConnectionStatus.ONLINE : IConnectionStatus.OFFLINE,
    };
  },
);

export const addConnection = createAsyncThunk('connections/addConnection', async ({ plugin, ...payload }: any) => {
  const connection = await API.connection.create(plugin, payload);
  return transformConnection(plugin, connection);
});

export const updateConnection = createAsyncThunk('connections/updateConnection', async (payload: ConnectionForm) => {});

export const slice = createSlice({
  name: 'connections',
  initialState,
  reducers: {},
  extraReducers(builder) {
    builder
      .addCase(init.fulfilled, (state, action) => {
        state.connections = action.payload;
      })
      .addCase(fetchConnections.fulfilled, (state, action) => {
        state.connections = state.connections.concat(action.payload.connections);
      })
      .addCase(addConnection.fulfilled, (state, action) => {
        state.connections.push(action.payload);
      })
      .addCase(testConnection.pending, (state, action) => {
        const existingConnection = state.connections.find((cs) => cs.unique === action.meta.arg.unique);
        if (existingConnection) {
          existingConnection.status = IConnectionStatus.TESTING;
        }
      })
      .addCase(testConnection.fulfilled, (state, action) => {
        const existingConnection = state.connections.find((cs) => cs.unique === action.payload.unique);
        if (existingConnection) {
          existingConnection.status = action.payload.status;
        }
      });
  },
});

export const {} = slice.actions;

export default slice.reducer;

export const selectAllConnections = (state: RootState) => state.connections.connections;

export const selectConnections = (state: RootState, plugin: string) =>
  state.connections.connections.filter((connection) => connection.plugin === plugin);

export const selectConnection = (state: RootState, unique: string) =>
  state.connections.connections.find((cs) => cs.unique === unique);
