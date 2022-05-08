// Package config 配置相关代码。本地和远程配置的读取和监听
package config

import (
	"git.cestong.com.cn/cecf/cecf-golib/pkg/errors"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

func CfgFromFile(cfgPath string, value interface{}) error {
	content, errRead := ioutil.ReadFile(cfgPath)
	if errRead != nil {
		return errors.Wrap(errRead, "read file [%s]", cfgPath)
	}
	if err := yaml.Unmarshal(content, value); err != nil {
		return errors.Wrap(err, "yaml unmarshal")
	}
	return nil
}
