// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stagemanager

import (
	"sync"

	"github.com/kasworld/gowasm3dgame/lib/w3dlog"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_connbytemanager"
)

type stageI interface {
	GetUUID() string
	GetConnManager() *w3d_connbytemanager.Manager
}

type stageList []stageI

func (aol stageList) Len() int { return len(aol) }
func (aol stageList) Swap(i, j int) {
	aol[i], aol[j] = aol[j], aol[i]
}
func (aol stageList) Less(i, j int) bool {
	ao1 := aol[i]
	ao2 := aol[j]
	return ao1.GetUUID() > ao2.GetUUID()
}

type Manager struct {
	log      *w3dlog.LogBase
	mutex    sync.RWMutex `prettystring:"hide"`
	id2stage map[string]stageI
}

func New(log *w3dlog.LogBase) *Manager {
	man := &Manager{
		log:      log,
		id2stage: make(map[string]stageI),
	}
	return man
}

func (man *Manager) Count() int {
	return len(man.id2stage)
}

func (man *Manager) GetAny() stageI {
	man.mutex.RLock()
	defer man.mutex.RUnlock()
	for _, v := range man.id2stage {
		return v
	}
	return nil
}

func (man *Manager) GetList() []stageI {
	man.mutex.RLock()
	defer man.mutex.RUnlock()
	rtn := make([]stageI, len(man.id2stage))
	i := 0
	for _, v := range man.id2stage {
		rtn[i] = v
		i++
	}
	return rtn
}

func (man *Manager) GetByUUID(uuid string) stageI {
	man.mutex.RLock()
	defer man.mutex.RUnlock()
	return man.id2stage[uuid]
}

func (man *Manager) Add(stg stageI) stageI {
	man.mutex.Lock()
	defer man.mutex.Unlock()
	old := man.id2stage[stg.GetUUID()]
	man.id2stage[stg.GetUUID()] = stg
	return old
}

func (man *Manager) Del(stg stageI) stageI {
	man.mutex.Lock()
	defer man.mutex.Unlock()
	old := man.id2stage[stg.GetUUID()]
	delete(man.id2stage, stg.GetUUID())
	return old
}
