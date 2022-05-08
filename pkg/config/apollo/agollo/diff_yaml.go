package agollo

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// ProcessYamlDiffUnderCertainKeys 只支持single doc的yaml
func ProcessYamlDiffUnderCertainKeys(observersDelegate map[string]Observer, namespaceWithSuffix string, oldYamlStr string, newYamlStr string) error {
	oldYamlViper, err := addYamlToViper(oldYamlStr)
	if err != nil {
		return err
	}

	newYamlViper, err := addYamlToViper(newYamlStr)
	if err != nil {
		return err
	}

	for k, v := range observersDelegate {
		go func(key string, observer Observer) {
			strList := strings.SplitN(key, namespaceToKeySep, 2)
			registeredNamespace := strList[0]
			if namespaceWithSuffix != registeredNamespace {
				return
			}
			viperKey := strList[1]
			if viperKey == "" {
				return
			}
			oldValue := oldYamlViper.Get(viperKey)
			newValue := newYamlViper.Get(viperKey)
			if oldValue == nil && newValue == nil {
				return
			}
			if reflect.DeepEqual(oldValue, newValue) {
				return
			}
			observer(viperKey, oldValue, newValue)
		}(k, v)
	}
	return nil
}

func addYamlToViper(yamlStr string) (*viper.Viper, error) {
	buf := bytes.NewBufferString(yamlStr)
	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(buf)
	if err != nil {
		return nil, err
	}
	return v, nil
}
