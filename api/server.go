/*
 * Copyright 2016 Robin Engel
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package api

import "time"

type loginResponse struct {
	moselResponse

	Successful bool
	Key        string
	ValidTo    time.Time
}

func NewLoginResponse() loginResponse {
	return loginResponse{
		moselResponse: newMoselResponse(),
	}
}

type nodeInfoRepsonse struct {
	moselResponse

	Data map[string]NodeInfo
}

func NewNodeInfoResponse() nodeInfoRepsonse {
	return nodeInfoRepsonse{
		moselResponse: newMoselResponse(),
		Data: make(map[string]NodeInfo),
	}
}