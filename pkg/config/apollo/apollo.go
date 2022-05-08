package apollo

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"git.cestong.com.cn/cecf/cecf-golib/pkg/utils"
	"github.com/pkg/errors"

	"github.com/spf13/viper"

	"git.cestong.com.cn/cecf/cecf-golib/pkg/config"

	"git.cestong.com.cn/cecf/cecf-golib/pkg/config/apollo/agollo"
)

const (
	// ConfigSuffixYaml ...
	ConfigSuffixYaml = ".yaml"

	// ConfigSuffixYml ...
	ConfigSuffixYml = ".yml"

	// ConfigSuffixJSON ...
	ConfigSuffixJSON = ".json"

	// ConfigSuffixToml ...
	ConfigSuffixToml = ".toml"

	// ConfigSuffixProperties ...
	ConfigSuffixProperties = ".properties"

	// ConfigFormatYaml ...
	ConfigFormatYaml = "yaml"
)

const (
	// 数据源前缀
	DataSourcePre = "jiagou.ds"
)

var (
	// ConfigFormats ...
	ConfigFormats = []string{
		ConfigSuffixYaml,
		ConfigSuffixYml,
		ConfigSuffixJSON,
		ConfigSuffixToml,
	}
)

// Apollo ...
type Apollo struct {
	conf        *agollo.Conf
	changeEvent <-chan *agollo.ChangeEvent
	stopChan    chan struct{}
}

// StartWithFile 从文件中初始化Apollo对象 默认支持 ConfigFormats 中的文件格式
func StartWithFile(fileName string) (apollo *Apollo, err error) {

	// 初始化viper
	viper.Reset()
	viper.SetConfigFile(fileName)
	suffix := path.Ext(fileName)
	if !Contains(ConfigFormats, suffix) {
		err = errors.New("file format is not supported !!!")
		return
	}
	suffix = strings.TrimPrefix(suffix, ".")
	viper.SetConfigType(suffix)
	viper.ReadInConfig()

	// 读取apollo的配置
	var ret config.ApolloConfig
	err = viper.UnmarshalKey("apollo", &ret)
	if err != nil {
		return
	}

	// 通过配置对象开始agollo的连接初始化
	apollo, err = StartWithConfig(&ret)
	return
}

// StartWithConfig 从配置中初始化Apollo对象
func StartWithConfig(config *config.ApolloConfig) (apollo *Apollo, err error) {
	// 初始化配置
	conf := &agollo.Conf{
		AppID:          config.AppID,
		Cluster:        config.Cluster,
		NameSpaceNames: config.Namespaces,
		CacheDir:       config.CacheDir,
		MetaAddr:       config.MetaAddr,
	}

	// 通过配置对象开始agollo的连接初始化
	err = agollo.StartWithConf(conf)
	if err != nil {
		return
	}

	apollo = &Apollo{
		conf:     conf,
		stopChan: make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-apollo.stopChan:
				return
			}
		}
	}()

	return
}

func GetSystemConfigByKeys(nodeMap map[string]interface{}) (err error) {
	for key, value := range nodeMap {
		err = viper.UnmarshalKey(key, value)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("key=%v", key))
			return
		}
	}

	return
}

// WatchUpdate watch配置变化
func (a *Apollo) WatchUpdate() <-chan *agollo.ChangeEvent {
	a.changeEvent = agollo.WatchUpdate()
	return a.changeEvent
}

// Reload 重新加载配置
func (a *Apollo) Reload() (err error) {

	var (
		buf = bytes.NewBufferString("")
		m   = &Map{
			m: map[string]interface{}{},
		}
		rawContent = ""
	)

	// 读取所有的namespace的内容
	for _, namespace := range a.conf.NameSpaceNames {
		suffix := path.Ext(namespace)
		switch suffix {
		case ConfigSuffixYaml, ConfigSuffixYml:
			// 过滤私有配置
			// 加换行符以避免配置文件中最后一行不是空行而报错'load config from apollo failed with err: error converting YAML to JSON: yaml: line 13: mapping values are not allowed in this context'
			rawContent = agollo.GetNameSpaceContent(namespace, "") + "\n"
		case ConfigSuffixProperties:
			namespace = strings.TrimSuffix(namespace, suffix)
			allKeys := agollo.GetAllKeys(namespace)
			// 默认过滤公共配置 properties格式的文本内容 只能读取公共配置，结构为key = json 结构
			for _, key := range allKeys {
				value := agollo.GetStringValueWithNameSpace(namespace, key, "")
				m.Set(key, value)
				// 校验json串
				err = checkJsonValid([]byte(value))
				if err != nil {
					m.Set(key, value)
				} else {
					var valueMap map[string]interface{}
					err = utils.Unmarshal([]byte(value), &valueMap)
					if err != nil {
						return errors.Wrap(err, fmt.Sprintf("utils.Unmarshal key = %v,value = %v", key, value))
					}
					m.Set(key, valueMap)
				}
			}
			if len(m.GetMap()) == 0 {
				continue
			}
			rawContent, err = jsonMapToYaml(m.GetMap())
			if err != nil {
				return errors.Wrap(err, "jsonMapToYaml")
			}
		}
		// 将apollo读取到的内容写入buf
		_, err = buf.WriteString(rawContent)
		if err != nil {
			return err
		}
	}

	// 校验格式
	viper.SetConfigType(ConfigFormatYaml)
	err = checkFormat(buf.Bytes())
	if err != nil {
		return err
	}

	// 从buf中读取配置信息
	err = viper.ReadConfig(buf)
	if err != nil {
		return err
	}

	return
}

// Stop 停止
func (a *Apollo) Stop() (err error) {
	a.stopChan <- struct{}{}
	return agollo.Stop()
}

// On 注册配置变更回调函数，名为global的namespace不需要后缀，其它例如名为config的需要后缀，比如config.yaml
func (a *Apollo) On(namespaceWithSuffix, key string, callback agollo.Observer) {
	agollo.On(namespaceWithSuffix, key, callback)
}

// 获取数据源namespace
func (a *Apollo) GetDataSourceNamespace(dataSourcePre string) (namespace string) {
	if len(dataSourcePre) == 0 {
		dataSourcePre = DataSourcePre
	}
	for _, n := range a.conf.NameSpaceNames {
		if strings.Contains(n, DataSourcePre) {
			suffix := path.Ext(n)
			namespace = strings.TrimSuffix(n, suffix)
			return
		}
	}
	return
}
