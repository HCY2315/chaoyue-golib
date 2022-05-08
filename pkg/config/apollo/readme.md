# apollo

## 使用介绍

### 读取配置
config.toml配置内容如下：
```toml
[cluster]
    domain = "10"
    id = "24001"

[apollo]
    app_id = "apigateway"
    cluster = "default"
    namespaces = ["config.yaml","global.properties"]
    meta_addr = "http://10.155.20.189:8180"
    cache_dir = "."
```

### apollo读取配置

```go
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
    // 检测config,yaml 命名空间的hystrix.ask_app_server的值得变化
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
```

### viper读取配置

```go
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

```

# viper

## 介绍

 viper 是一个配置解决方案

特性：

- 支持 JSON/TOML/YAML/HCL/envfile/Java properties 等多种格式的配置文件；
- 可以设置监听配置文件的修改，修改时自动加载新的配置；
- 从环境变量、命令行选项和`io.Reader`中读取配置；
- 从远程配置系统中读取和监听修改，如 etcd/Consul；
- 代码逻辑中显示设置键值。

## 读取键

config.yaml文件内容如下

```yaml
mysql:
  ip: 127.0.0.1
  port: "4802"
  user: ccp_dev
  password: ccp_dev
  dbname: ccp_apollo

list:
  - 1
  - 2
  - 3
  - 4
```

从配置文件初始化viper

```go
func InitViperConfig() error {
   viper.SetConfigName("config")
   viper.SetConfigType("yaml")
   viper.AddConfigPath(".")
   viper.SetDefault("redis.port", 6381)
   err := viper.ReadInConfig()
   if err != nil {
      log.Fatalf("read config failed: %v", err)
   }
   return err
}
```

viper test 案例

```go
func TestInitViperConfig(t *testing.T) {
   var err error
   InitViperConfig()

   // GetType系列方法可以返回指定类型的值

   // 获取当个key
   fmt.Println("mysql.port is Set: ", viper.IsSet("mysql.port"))
   // mysql.port is Set:  true

   fmt.Println("mysql.port: ", viper.Get("mysql.port"))
   // mysql.port:  4802
   fmt.Println("redis.port: ",viper.Get("redis.port"))
   // redis.port:  6381

   // UnmarshalKey
   mysql := &Mysql{}
   err = viper.UnmarshalKey("mysql",mysql)
   fmt.Println("mysql: ",mysql)
   fmt.Println("mysql: ",viper.Get("mysql"))
   fmt.Println("mysql: ",viper.GetStringMap("mysql"))
   fmt.Println("mysql: ",viper.GetStringMapString("mysql"))
   // mysql:  &{10.155.10.48 4802 ccp_dev ccp_dev ccp_apollo}
	
   fmt.Println("list: ",viper.GetStringSlice("list"))
   // list: [1 2 3 4]


   // 获取所有的配置
   fmt.Println("list: ",viper.AllSettings())
   // map[list:[1 2 3 4] mysql:map[dbname:ccp_apollo ip:10.155.10.48 password:ccp_dev port:4802 user:ccp_dev] redis:map[port:6381]]

   _ = err


}
```



## 小工具

### json toml yaml互转 

https://toolkit.site/zh/format.html

https://onlineyamltools.com/

### json yaml 转struct插件

jsontogo


