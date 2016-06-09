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
	"log"
)

type nodeRespHandler struct {
	cache *dataCache
}

func NewNodeRespHandler(cache *dataCache) *nodeRespHandler {
	log.Println(cache)
	return &nodeRespHandler{
		cache: cache,
	}
}

func (handler nodeRespHandler) handleNodeResp(node string, resp api.NodeResponse) {
	log.Println(resp)
	handler.cache.Add(node, "test", resp.Time, float64(resp.Test))
}
