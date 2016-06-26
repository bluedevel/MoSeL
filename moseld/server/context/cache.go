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
	"time"
	"github.com/bluedevel/mosel/api"
	"sync"
	"errors"
)

type dataCache struct {
	points map[string][]dataPoint

	m      sync.Mutex
}

type dataPoint struct {
	Time time.Time
	Info api.NodeInfo
}

func NewDataCache() *dataCache {
	c := &dataCache{}
	c.points = make(map[string][]dataPoint)
	return c
}

func (cache *dataCache) Add(node string, t time.Time, info api.NodeInfo) {

	var arr []dataPoint

	cache.m.Lock()

	if _, ok := cache.points[node]; !ok {
		arr = make([]dataPoint, 0)
	} else {
		arr = cache.points[node]
	}

	arr = append(arr, dataPoint{
		Time: t.Round(time.Second),
		Info: info,
	})

	cache.points[node] = arr
	cache.m.Unlock()
}

func (cache *dataCache) Get(node string, t time.Time) (api.NodeInfo, error) {

	points, err := cache.GetAll(node)

	if err != nil {
		return api.NodeInfo{}, err
	}

	for _, p := range points {
		if p.Time.Unix() == t.Unix() {
			return p.Info, nil
		}
	}

	return api.NodeInfo{}, errors.New("No datapoint found for time " + t.String())
}

func (cache *dataCache) GetSince(node string, t time.Time) ([]dataPoint, error) {

	points, err := cache.GetAll(node)

	if err != nil {
		return nil, err
	}

	result := make([]dataPoint, 0)

	for _, p := range points {
		if p.Time.Unix() > t.Unix() {
			result = append(result, p)
		}
	}

	return result, nil
}

func (cache *dataCache) GetAll(node string) ([]dataPoint, error) {
	cache.m.Lock()

	points, ok := cache.points[node]

	if !ok {
		cache.m.Unlock()
		return nil, errors.New("No node with name " + node)
	}

	cache.m.Unlock()
	return points, nil
}

func (cache *dataCache) GetNodes() []string {
	nodes := make([]string, len(cache.points))

	i := 0
	for k := range cache.points {
		nodes[i] = k
		i++
	}

	return nodes
}