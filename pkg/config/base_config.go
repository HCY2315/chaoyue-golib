package config

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
)

func ParseConfigFile(file string, out interface{}) error {
	_, err := toml.DecodeFile(file, out)
	if err != nil {
		return err
	}
	return nil
}

type ApolloConfig struct {
	AppID      string   `mapstructure:"app_id" toml:"app_id" json:"app_id"`
	Cluster    string   `mapstructure:"cluster" toml:"cluster" json:"cluster"`
	Namespaces []string `mapstructure:"namespaces" toml:"namespaces" json:"namespaces"`
	MetaAddr   string   `mapstructure:"meta_addr" toml:"meta_addr" json:"meta_addr"`
	CacheDir   string   `mapstructure:"cache_dir" toml:"cache_dir" json:"cache_dir"`
}

type GlobalConfig struct {
	Domain string   `toml:"Domain"`
	ID     string   `toml:"ID"`
	ZKHost []string `toml:"ZKHost"`
}

type ModuleConfig struct {
	Cluster *GlobalConfig
	Log     *LogConfig
}

func (mc *ModuleConfig) Validate() (err error) {
	if mc.Cluster == nil {
		err = fmt.Errorf("ModuleConfig.Cluster field required")
		return
	}
	return nil
}

type BaseConfig struct {
	Db         *DbConfig
	NoSql      *NoSqlConfig
	NoSqlSlave *NoSqlConfig
	ImNoSql    *NoSqlConfig
	LocalNoSql *NoSqlConfig
	Log        *LogConfig
	Zk         *ZkConfig
	Statsd     *StatsdConfig
}

type DbConfig struct {
	Host        string `mapstructure:"ip" json:"ip"`
	Port        string `mapstructure:"port" json:"port"`
	User        string `mapstructure:"user" json:"user"`
	Password    string `mapstructure:"password" json:"password"`
	DbName      string `mapstructure:"dbname" json:"dbname"`
	MaxLifetime time.Duration
	MaxIdleConn int
	MaxConn     int
	ShowSql     bool
}

// NoSqlConfig redis配置项，对应global.redis
type NoSqlConfig struct {
	Addr         string   `mapstructure:"conn" json:"conn"`             // "ip:port", 优先判断ClusterAddr，若ClusterAddr为空，再连接Addr
	Password     string   `mapstructure:"password" json:"password"`     // Password
	DbNum        int      `mapstructure:"idx,string" json:"idx,string"` // DbNum
	PoolSize     int      // 默认DefaultPoolSize，超过MaxPoolSize，则置为MaxPoolSize
	ClusterAddr  []string // ["ip:port"]，优先判断ClusterAddr，若为空，再判断SentinelAddr，若为空，则连接Addr
	SpecialSize  bool     // 特殊池大小，true: 可自己设定PoolSize
	SentinelAddr []string // ["ip:port"]，优先判断ClusterAddr，若为空，再判断SentinelAddr，若为空，则连接Addr
	MasterName   string   // sentinel masterName
}

// ESConfig elasticsearch配置项，对应global.elasticsearch
type ESConfig struct {
	Addrs []string `mapstructure:"addrs" json:"addrs"`
}

type StatsdConfig struct {
	Addrs     []string `mapstructure:"addrs" json:"addrs"` // ["ip:port"]
	BaseAddrs []string `mapstructure:"base" json:"base"`
}

type LogConfig struct {
	FileName    string `mapstructure:"file_name" toml:"file_name"`
	Console     bool   `mapstructure:"console" toml:"console"`
	Level       string `mapstructure:"level" toml:"level"`
	MaxFileSize int    `mapstructure:"max_file_size" toml:"max_file_size"`
	MaxDays     int    `mapstructure:"max_days" toml:"max_days"`
	Compress    bool   `mapstructure:"compress" toml:"compress"`
	Options     []zap.Option
}

// HystrixConfig 断路器hystrix配置
type HystrixConfig struct {
	Enable                 bool `mapstructure:"enable" toml:"enable"`
	Timeout                int  `mapstructure:"timeout" toml:"timeout"`
	MaxConcurrentRequests  int  `mapstructure:"max_concurrent_requests" toml:"max_concurrent_requests"`
	RequestVolumeThreshold int  `mapstructure:"request_volume_threshold" toml:"request_volume_threshold"`
	SleepWindow            int  `mapstructure:"sleep_window" toml:"sleep_window"`
	ErrorPercentThreshold  int  `mapstructure:"error_percent_threshold" toml:"error_percent_threshold"`
}

type ZkConfigAuth struct {
	Scheme   string
	User     string
	Password string
}

type ZkConfig struct {
	Addr    []string `mapstructure:"addrs"`
	Timeout int
	Auth    *ZkConfigAuth
}
