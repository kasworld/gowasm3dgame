// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)

package serverconfig

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/kasworld/gowasm3dgame/lib/w3dlog"
	"github.com/kasworld/prettystring"
	"gopkg.in/ini.v1"
)

type Config struct {
	LogLevel              w3dlog.LL_Type `argname:""`
	SplitLogLevel         w3dlog.LL_Type `argname:""`
	BaseLogDir            string         `default:"/tmp/"  argname:""`
	ServerDataFolder      string         `default:"./serverdata" argname:""`
	ClientDataFolder      string         `default:"./www" argname:""`
	ServicePort           string         `default:":34101"  argname:""`
	AdminPort             string         `default:":34201"  argname:""`
	ConcurrentConnections int            `default:"1000" argname:""`
	ActTurnPerSec         float64        `default:"30.0" argname:""`
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

func LoadIni(urlpath string, config interface{}) error {
	datas, err := loadData(urlpath)
	if err != nil {
		return err
	}
	f, err := ini.Load(datas)
	if err != nil {
		return err
	}
	if err := f.MapTo(config); err != nil {
		return err
	}
	return nil
}

func loadData(urlpath string) ([]byte, error) {
	var fd io.Reader
	u, err := url.Parse(urlpath)
	if err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		resp, err := http.Get(urlpath)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		fd = resp.Body
	} else {
		ffd, err := os.Open(urlpath)
		if err != nil {
			return nil, err
		}
		defer ffd.Close()
		fd = ffd
	}
	return ioutil.ReadAll(fd)
}

func SaveIni(filename string, config interface{}) error {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	return generateToIni(config, fd)
}

func generateToIni(config interface{}, w io.Writer) error {
	cfg := ini.Empty()
	ini.ReflectFrom(cfg, config)
	_, err := cfg.WriteToIndent(w, "\t")
	return err
}
