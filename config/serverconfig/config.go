// Copyright 2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)

package serverconfig

import (
	"fmt"
	"path/filepath"

	"github.com/kasworld/gowasm3dgame/lib/w3dlog"
	"github.com/kasworld/prettystring"
)

type Config struct {
	LogLevel              w3dlog.LL_Type `default:"31" argname:""`
	SplitLogLevel         w3dlog.LL_Type `default:"0" argname:""`
	BaseLogDir            string         `default:""  argname:""`
	ServerDataFolder      string         `default:"./serverdata" argname:""`
	ClientDataFolder      string         `default:"./clientdata" argname:""`
	ServicePort           string         `default:"34101"  argname:""`
	AdminPort             string         `default:"34201"  argname:""`
	WebAdminID            string         `default:"root" argname:""`
	WebAdminPass          string         `default:"password" argname:"" prettystring:"hidevalue"`
	ServiceHostBase       string         `default:"http://localhost" argname:""`
	ConcurrentConnections int            `default:"1000" argname:""`

	ActTurnPerSec float64 `default:"30.0" argname:""`
}

func (config Config) MakeLogDir() string {
	rstr := filepath.Join(
		config.BaseLogDir,
		"w3dserver.logfiles",
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config Config) MakePIDFileFullpath() string {
	rstr := filepath.Join(
		config.BaseLogDir,
		"w3dserver.pid",
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config Config) MakeOutfileFullpath() string {
	rstr := "w3dserver.out"
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config Config) StringForm() string {
	return prettystring.PrettyString(config, 4)
}
