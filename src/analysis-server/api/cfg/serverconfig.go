package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ServerConfFac struct {
	Path *string
}

type ServerConf struct {
	Port    int
	Cores   int
	BaseUrl string
}

func (c *ServerConf) CheckValid() error {
	if c.Port <= 0 {
		return fmt.Errorf("ServerConf need Port")
	}
	if c.Cores < 1 {
		return fmt.Errorf("ServerConf need Cores")
	}
	if len(c.BaseUrl) == 0 {
		return fmt.Errorf("ServerConf need BaseUrl")
	}
	return nil
}

func (fac ServerConfFac) ParseConfig() (*ServerConf, error) {
	data, err := ioutil.ReadFile(*fac.Path)
	if err != nil {
		return nil, err
	}
	config := &ServerConf{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
