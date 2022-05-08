package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"ccp-server/framework/apollo"
	"ccp-server/framework/config"
)

// HystrixConfig ...
type HystrixConfig struct {
	Enable                 bool `mapstructure:"enable"`
	Timeout                int  `mapstructure:"timeout"`
	MaxConcurrentRequests  int  `mapstructure:"max_concurrent_requests"`
	RequestVolumeThreshold int  `mapstructure:"request_volume_threshold"`
	SleepWindow            int  `mapstructure:"sleep_window"`
	ErrorPercentThreshold  int  `mapstructure:"error_percent_threshold"`
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 读取config.toml，并创建apollo实例
	apolloInstance, err := apollo.StartWithFile("./config.toml")
	if err != nil {
		log.Printf("start apollo failed with err: %v", err)
		return
	}

	nodeMeta := &config.GlobalConfig{}

	viper.UnmarshalKey("cluster", nodeMeta)
	log.Printf("cluster1 = %+v", nodeMeta)

	// 从apollo中读取公共配置和私有配置至viper
	err = apolloInstance.Reload()
	if err != nil {
		log.Printf("load config from apollo failed with err: %v", err)
		return
	}
	nodeMeta = &config.GlobalConfig{}
	viper.UnmarshalKey("cluster", nodeMeta)
	log.Printf("cluster2 = %+v", nodeMeta)
	apolloInstance.On("config.yaml", "hystrix.ask_app_server", func(key string, oldValue interface{}, newValue interface{}) {
		hystrixConfig := &HystrixConfig{}
		hystrixConfigDecoder := &mapstructure.DecoderConfig{
			Metadata: nil,
			Result:   hystrixConfig,
			TagName:  "mapstructure",
		}
		decoder, _ := mapstructure.NewDecoder(hystrixConfigDecoder)
		decoder.Decode(newValue)
		log.Printf("newValue = %+v", hystrixConfig)
	})
	apolloInstance.On("config.yaml", "log.console", func(key string, oldValue interface{}, newValue interface{}) {
		if _, ok := newValue.(bool); ok {
			log.Printf("successfully decode")
		} else {
			log.Printf("failed to decode")
		}

	})
	apolloInstance.On("config.yaml", "log.level", func(key string, oldValue interface{}, newValue interface{}) {
		if newLevel, ok := newValue.(string); ok {
			log.Printf("successfully decode and newLevel = %s", newLevel)
		} else {
			log.Printf("failed to decode")
		}

	})
	apolloInstance.On("global", "test_change", func(key string, oldValue interface{}, newValue interface{}) {
		log.Printf("key = %s, oldValue = %v, newValue = %v", key, oldValue, newValue)
	})

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-sigs:
				ticker.Stop()
				goto DONE
			case <-ticker.C:
				hystrixTimeout := viper.GetInt32("hystrix.ask_app_server.timeout")
				log.Printf("hystrix timeout = %d", hystrixTimeout)
				alarmURL := viper.GetString("alarm_url")
				log.Printf("alarm url = %s", alarmURL)
			}
		}
	DONE:
		done <- struct{}{}
	}()

	// -------------获取单项配置--------------
	// 获取全局配置
	redisAddr := viper.GetString("redis.conn")
	redisPasswd := viper.GetString("redis.password")
	log.Printf("redis addr = %s", redisAddr)
	log.Printf("redis passwd = %s", redisPasswd)
	// 获取apigateway配置
	isHystrixEnabled := viper.GetBool("hystrix.ask_app_server.enable")
	hystrixTimeout := viper.GetInt32("hystrix.ask_app_server.timeout")
	log.Printf("hystrix enable = %t", isHystrixEnabled)
	log.Printf("hystrix timeout = %d", hystrixTimeout)

	// -------------将配置项反序列化至结构体--------------
	// 获取日志
	cfg := &config.LogConfig{}
	err = viper.UnmarshalKey("log", cfg)
	if err != nil {
		log.Printf("unmarshal log failed with error: %v", err)
	} else {
		log.Printf("log config is %+v", cfg)
	}

	// unmarshal全局配置的database项
	dbConfig := &config.DbConfig{}
	viper.UnmarshalKey("database", dbConfig)
	log.Printf("dbConfig = %+v", dbConfig)
	// unmarshal全局配置的redis项
	redisConfig := &config.NoSqlConfig{}
	viper.UnmarshalKey("redis", redisConfig)
	log.Printf("redisConfig = %+v", redisConfig)
	// unmarshal全局配置的es项
	esConfig := &config.ESConfig{}
	viper.UnmarshalKey("elasticsearch", esConfig)
	log.Printf("esConfig = %+v", esConfig)
	// unmarshal全局配置的kafka项
	// TODO
	// unmarshal statsd
	statsdConfig := &config.StatsdConfig{}
	viper.UnmarshalKey("statsd", statsdConfig)
	log.Printf("statsdConfig = %+v", statsdConfig)

	// 告警地址
	alarmURL := viper.GetString("alarm_url")
	log.Printf("alarmURL = %s", alarmURL)

	// unmarshal全局配置的zookeeper项
	zkConfig := &config.ZkConfig{}
	viper.UnmarshalKey("zookeeper", zkConfig)
	log.Printf("zkConfig = %+v", zkConfig)

	// unmarshal私有配置hystrix
	hystrixConfigMap := map[string]*HystrixConfig{}
	err = viper.UnmarshalKey("hystrix", &hystrixConfigMap)
	if err != nil {
		log.Printf("unmarshal hystrix failed with err: %v", err)
	}
	for k, v := range hystrixConfigMap {
		log.Printf("add hystrix %+v for %s", v, k)
	}

	<-done
}
