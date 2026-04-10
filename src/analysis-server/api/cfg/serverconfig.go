package cfg

import (
	"encoding/json"
	"fmt"
	"os"
)

type ServerConfFac struct {
	Path *string
}

type ServerConf struct {
	Port        int
	SynDuration int
	BaseUrl     string
}

func (c *ServerConf) CheckValid() error {
	if c.Port <= 0 {
		return fmt.Errorf("ServerConf need Port")
	}
	if c.SynDuration < 1 {
		return fmt.Errorf("ServerConf need SynDuration")
	}
	if len(c.BaseUrl) == 0 {
		return fmt.Errorf("ServerConf need BaseUrl")
	}
	return nil
}

func (fac ServerConfFac) ParseConfig() (*ServerConf, error) {
	data, err := os.ReadFile(*fac.Path)
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
