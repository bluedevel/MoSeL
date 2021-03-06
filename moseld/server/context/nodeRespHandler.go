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
package context

import (
	"github.com/bluedevel/mosel/api"
)

type nodeRespHandler struct {
	cache           *dataCache
	dataPersistence DataPersistence
}

func NewNodeRespHandler(cache *dataCache, dataPersistence DataPersistence) (*nodeRespHandler, error) {
	return &nodeRespHandler{
		cache:           cache,
		dataPersistence: dataPersistence,
	}, nil
}

func (handler nodeRespHandler) handleNodeResp(node string, resp api.NodeResponse) {
	//log.Println(resp)
	handler.cache.Add(node, resp.Time, resp.NodeInfo)

	if handler.dataPersistence != nil {
		handler.dataPersistence.Add(node, resp.Time, resp.NodeInfo)
	}
	//log.Println(resp)
}
