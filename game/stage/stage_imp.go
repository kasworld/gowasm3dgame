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

package stage

import (
	"fmt"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_connbytemanager"
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
