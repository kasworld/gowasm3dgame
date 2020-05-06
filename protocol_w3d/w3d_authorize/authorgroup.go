// Copyright 2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package w3d_authorize

import (
	"fmt"

	"github.com/kasworld/gowasm3dgame/config/authdata"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idcmd"
)

func AddAdminKey(key string) error {
	var err error
	if _, exist := authdata.Authkey2Admin[key]; exist {
		err = fmt.Errorf("key %v exist, overwright", key)
	}
	authdata.Authkey2Admin[key] = [2][]string{
		[]string{"Login", "Admin"}, []string{"DelAfterLogin"},
	}
	return err
}

var allAuthorizationSet = map[string]*AuthorizedCmds{
	"PreLogin": NewByCmdIDList([]w3d_idcmd.CommandID{
		w3d_idcmd.Login,
	}),

	"DelAfterLogin": NewByCmdIDList([]w3d_idcmd.CommandID{
		w3d_idcmd.Login,
	}),

	"Login": NewByCmdIDList([]w3d_idcmd.CommandID{
		w3d_idcmd.Heartbeat,
		w3d_idcmd.Chat,
	}),
	"Admin": NewByCmdIDList([]w3d_idcmd.CommandID{}),
}

func NewPreLoginAuthorCmdIDList() *AuthorizedCmds {
	return allAuthorizationSet["PreLogin"].Duplicate()
}

func (acicl *AuthorizedCmds) UpdateByAuthKey(key string) error {
	ag, exist := authdata.Authkey2Admin[key]
	if !exist {
		ag = [2][]string{[]string{"Login"}, []string{"DelAfterLogin"}}
	}
	// process include
	for _, authgroupname := range ag[0] {
		cmdidList := allAuthorizationSet[authgroupname]
		if cmdidList == nil {
			return fmt.Errorf("Can't Found authgroup %v", authgroupname)
		}
		acicl.Union(cmdidList)
	}
	// process exclude
	for _, authgroupname := range ag[1] {
		cmdidList := allAuthorizationSet[authgroupname]
		if cmdidList == nil {
			return fmt.Errorf("Can't Found authgroup %v", authgroupname)
		}
		acicl.SubIntersection(cmdidList)
	}
	return nil
}
