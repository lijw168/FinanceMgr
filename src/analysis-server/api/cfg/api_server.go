// api server 四种类型
// 1. user 提供volume 的基本功能
// 2. admin 提供集群管理和操作四
// 3. heartbeat
// 4. monitor 提供tsdb监控
package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	"common/config"
)

type ApiServerConf struct {
	LogConf    *config.LogConf   `json:"log"`
	ServerConf *ServerConf       `json:"server"`
	MysqlConf  *config.MysqlConf `json:"mysql"`
}

func IsNilInterface(i interface{}) bool {
	return reflect.ValueOf(i).IsNil()
}

func (a *ApiServerConf) CheckValid() error {
	if a.LogConf == nil {
		return fmt.Errorf("need log config")
	} else if a.ServerConf == nil {
		return fmt.Errorf("need server config")
	} else if a.MysqlConf == nil {
		return fmt.Errorf("need mysql config")
	}

	checkers := map[string]config.ConfigCheck{
		"log":    a.LogConf,
		"server": a.ServerConf,
		"mysql":  a.MysqlConf,
	}

	for k, checker := range checkers {
		if IsNilInterface(checker) {
			continue
		}
		if err := checker.CheckValid(); err != nil {
			return fmt.Errorf("%s err %v", k, err)
		}
	}
	return nil
}

func ParseApiServerConfig(path *string) (*ApiServerConf, error) {
	data, err := ioutil.ReadFile(*path)
	if err != nil {
		return nil, err
	}
	config := &ApiServerConf{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
