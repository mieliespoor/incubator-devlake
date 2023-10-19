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

import { Colors } from '@blueprintjs/core';
import styled from 'styled-components';

export const Wrapper = styled.div`
  h2 {
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0;
    margin-bottom: 16px;
    padding: 0;
    font-weight: 600;
    color: ${Colors.GREEN5};

    .bp5-icon {
      margin-right: 8px;
    }
  }

  h5 {
    margin: 8px 0;
  }

  p {
    color: #292b3f;
  }
`;

export const Action = styled.div`
  color: #7497f7;

  span {
    cursor: pointer;
  }

  span + span {
    margin-left: 8px;
  }
`;

export const ApiKey = styled.div`
  display: flex;
  align-items: center;

  & > div {
    max-width: 50%;
    margin-right: 8px;
  }
`;

export const Tips = styled.div`
  margin-top: 8px;
`;
