// Copyright 2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
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
	"sort"
	"strconv"
	"strings"
	"sync"
)

type stageI interface {
	GetUUID() string
}

type stageList []stageI

func (stgl stageList) Len() int { return len(stgl) }
func (stgl stageList) Swap(i, j int) {
	stgl[i], stgl[j] = stgl[j], stgl[i]
}
func (stgl stageList) Less(i, j int) bool {
	ao1 := stgl[i]
	ao2 := stgl[j]
	return ao1.GetUUID() > ao2.GetUUID()
}

type logI interface {
	Error(format string, v ...interface{})
}

type Manager struct {
	log      logI
	mutex    sync.RWMutex `prettystring:"hide"`
	id2stage map[string]stageI
}

func New(log logI) *Manager {
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

// GetStageByStageToEnter return
// random stage if empty
// stage by number if parse as int
// stage by uuid
func (man *Manager) GetStageByStageToEnter(StageToEnter string) stageI {
	StageToEnter = strings.TrimSpace(StageToEnter)
	if StageToEnter == "" {
		return man.GetAny()
	}
	if i64, err := strconv.ParseInt(StageToEnter, 0, 64); err == nil {
		stgList := man.GetList()
		sort.Sort(stageList(stgList))
		i := int(i64) % len(stgList)
		return stgList[i]
	}
	return man.GetByUUID(StageToEnter)
}
