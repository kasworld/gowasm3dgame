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

	"github.com/kasworld/gowasm3dgame/game/stage"
	"github.com/kasworld/gowasm3dgame/lib/w3dlog"
)

type Manager struct {
	log      *w3dlog.LogBase
	mutex    sync.RWMutex `prettystring:"hide"`
	id2stage map[string]*stage.Stage
}

func New(log *w3dlog.LogBase) *Manager {
	man := &Manager{
		log:      log,
		id2stage: make(map[string]*stage.Stage),
	}
	return man
}

func (man *Manager) Count() int {
	return len(man.id2stage)
}

func (man *Manager) GetAny() *stage.Stage {
	man.mutex.RLock()
	defer man.mutex.RUnlock()
	for _, v := range man.id2stage {
		return v
	}
	return nil
}

func (man *Manager) GetList() []*stage.Stage {
	man.mutex.RLock()
	defer man.mutex.RUnlock()
	rtn := make([]*stage.Stage, len(man.id2stage))
	i := 0
	for _, v := range man.id2stage {
		rtn[i] = v
		i++
	}
	return rtn
}

func (man *Manager) GetByUUID(uuid string) *stage.Stage {
	man.mutex.RLock()
	defer man.mutex.RUnlock()
	return man.id2stage[uuid]
}

func (man *Manager) Add(stg *stage.Stage) *stage.Stage {
	man.mutex.Lock()
	defer man.mutex.Unlock()
	old := man.id2stage[stg.UUID]
	man.id2stage[stg.UUID] = stg
	return old
}

func (man *Manager) Del(stg *stage.Stage) *stage.Stage {
	man.mutex.Lock()
	defer man.mutex.Unlock()
	old := man.id2stage[stg.UUID]
	delete(man.id2stage, stg.UUID)
	return old
}
