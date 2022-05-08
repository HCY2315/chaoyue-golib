package apollo

import (
	"log"
	"testing"

	"github.com/HCY2315/chaoyue-golib/pkg/config"
	"github.com/HCY2315/chaoyue-golib/pkg/config/apollo/agollo"

	"github.com/spf13/viper"
)

type Mysql struct {
	IP       string `mapstructure:"ip"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
}

type Addr struct {
	Addrs []string `mapstructure:"addrs"`
}

func TestApollo_StartWithFile(t *testing.T) {

	// 从配置文件中读取
	//apollo, err := StartWithFile("./config.json")
	//apollo, err := StartWithFile("./config.toml")
	apollo, err := StartWithFile("./config.toml")

	panicErr(err)
	err = apollo.Reload()
	panicErr(err)
	defer apollo.Stop()
	go func() {
		for a := range apollo.WatchUpdate() {
			UpdateConfig(a, apollo)
		}
	}()

	// keys
	t.Log("kafka_async_queue = ", viper.Get("kafka_async_queue"))
	t.Log("zookeeper = ", viper.Get("zookeeper"))
	t.Log("viper.AllKeys() = ", viper.AllKeys())
	t.Log("viper.Get(log) = ", viper.Get("log"))
	t.Log("viper.Get(database) = ", viper.Get("database"))
	t.Log("viper.Get(ss) = ", viper.Get("ss"))
	t.Log("master.config.port", viper.Get("master.config.port"))
	t.Log("master.config.password", viper.Get("master.config.password"))

	m, s, slaveEnable := GetDataSourceDb()
	t.Logf("master = %+v,slave = %+v,slaveEnable = %v", m, s, slaveEnable)

	// 非json串
	t.Log("viper.Get(alarm_url) = ", viper.Get("alarm_url"))

	// zookeeper
	var zkAddr *Addr
	//err = json.Unmarshal([]byte(viper.GetString("zookeeper")),&zkAddr)
	err = viper.UnmarshalKey("zookeeper", &zkAddr)
	panicErr(err)
	t.Log("zookeeper = ", zkAddr)

	// mysql
	var db *Mysql
	err = viper.UnmarshalKey("database", &db)
	panicErr(err)
	t.Log("db = ", db)
	select {}
}

// 更新配置
func UpdateConfig(a *agollo.ChangeEvent, apollo *Apollo) {
	err := apollo.Reload()
	panicErr(err)
	log.Println(a.Namespace, a.Changes)
	log.Println(apollo.GetDataSourceNamespace(""))
	if a.Namespace == apollo.GetDataSourceNamespace("") {
		log.Println(GetDataSourceDb())
	} else {
		log.Println("WatchUpdate later viper.Get(version) = ", viper.Get("log.level"))
	}
}

func TestApollo_StartWithConfig(t *testing.T) {

	// 从配置中读取
	apollo, err := StartWithConfig(&config.ApolloConfig{
		AppID:      "blackwordserver",
		Cluster:    "default",
		Namespaces: []string{"config.yaml", "global.properties"},
		MetaAddr:   "http://10.155.19.204:8180",
		CacheDir:   ".",
	})

	panicErr(err)
	err = apollo.Reload()
	panicErr(err)
	defer apollo.Stop()
	go func() {
		for range apollo.WatchUpdate() {
			err = apollo.Reload()
			panicErr(err)
			t.Log("WatchUpdate later viper.Get(version) = ", viper.Get("version"))
		}
	}()
	//for _,key := range viper.AllKeys() {
	//	t.Log("key =  ",key,"value = ",viper.Get(key))
	//}

	// keys
	t.Log("kafka_async_queue = ", viper.Get("kafka_async_queue"))
	t.Log("zookeeper = ", viper.Get("zookeeper"))
	t.Log("viper.AllKeys() = ", viper.AllKeys())
	t.Log("viper.Get(log) = ", viper.Get("log"))
	t.Log("viper.Get(database) = ", viper.Get("database"))

	// 非json串
	t.Log("viper.Get(alarm_url) = ", viper.Get("alarm_url"))

	// zookeeper
	var zkAddr *Addr
	err = viper.UnmarshalKey("zookeeper", &zkAddr)
	panicErr(err)
	t.Log("zookeeper = ", zkAddr)

	// mysql
	var db *Mysql
	err = viper.UnmarshalKey("database", &db)
	panicErr(err)
	t.Log("db = ", db)
	for {

	}
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
