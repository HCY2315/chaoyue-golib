package config

import (
	"github.com/HCY2315/chaoyue-golib/pkg/errors"
	"github.com/HCY2315/chaoyue-golib/pkg/utils"
)

type AppConfig struct {
	ServiceID     string `json:"serviceID"`
	RunModeString string
	runMode       utils.RunMode
	HostName      string
	HostIP        string
}

func (c AppConfig) RunMode() utils.RunMode {
	return c.runMode
}

// REF
func (c *AppConfig) ValidateAndLoad() error {
	var err error
	c.runMode, err = utils.ParseRunMode(c.RunModeString)
	if err != nil {
		return errors.Wrap(err, "解析RunMode:%s", c.RunModeString)
	}
	return nil
}
