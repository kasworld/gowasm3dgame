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

package stage2d

import (
	"fmt"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_connbytemanager"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
)

func (stg *Stage) String() string {
	return fmt.Sprintf("Team(%v)", len(stg.Teams))
}

func (stg *Stage) GetUUID() string {
	return stg.UUID
}

func (stg *Stage) GetConnManager() *w3d_connbytemanager.Manager {
	return stg.Conns
}

func (stg *Stage) ToPacket_StatsInfo() *w3d_obj.RspStatsInfo_data {
	rtn := &w3d_obj.RspStatsInfo_data{}
	for _, bt := range stg.Teams {
		teamStats := w3d_obj.TeamStat{
			UUID:     bt.UUID,
			Alive:    bt.IsAlive,
			AP:       int(bt.ActPoint),
			Score:    int(bt.Score),
			Kill:     bt.Kill,
			Death:    bt.Death,
			Color24:  uint32(bt.Color24),
			ActStats: bt.ActStats,
		}
		rtn.Stats = append(rtn.Stats, teamStats)
	}
	return rtn
}
