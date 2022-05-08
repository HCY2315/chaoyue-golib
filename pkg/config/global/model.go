package global

// mysql
type DataBase struct {
	Ip       string `mapstructure:"ip" json:"ip"`
	Port     string `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	PassWord string `mapstructure:"password" json:"password"`
	DbName   string `mapstructure:"dbname" json:"dbname"`
}

// redis
type NoSql struct {
	Host     string        `mapstructure:"conn" json:"conn"`
	PassWord string        `mapstructure:"password" json:"password"`
	DbNum    string        `mapstructure:"idx" json:"idx"`
	Sentinel NoSqlSentinel `mapstructure:"sentinel" json:"sentinel"` // 哨兵
}

// Sentinel
type NoSqlSentinel struct {
	Addrs      []string `mapstructure:"addrs" json:"addrs"`
	MasterName string   `mapstructure:"master_name" json:"master_name"`
}

// kafka
type Kafka struct {
	Brokers []string `mapstructure:"brokers" json:"brokers"`
	GroupId string   `mapstructure:"groupId" json:"groupId"`
	Topic   string   `mapstructure:"topic" json:"topic"`
}

// KafkaAsyncQueue
type KafkaAsyncQueue struct {
	Brokers map[string][]string `mapstructure:"brokers" json:"brokers"`
	Groups  map[string]string   `mapstructure:"groups" json:"groups"`
	Topics  map[string]string   `mapstructure:"topics" json:"topics"`
}

// grafna | influxdb
type Statsd struct {
	Addrs      []string `mapstructure:"addrs" json:"addrs"`
	Base       []string `mapstructure:"base" json:"base"`             // 主要用于基础服务的metrics上报
	Statistics []string `mapstructure:"statistics" json:"statistics"` // 用于radar和notification的metrics的上报
}

// 熔断
type Fusing struct {
	Group  string `mapstructure:"group" json:"group"`   // 群组限流等级
	Member string `mapstructure:"member" json:"member"` // 成员限流等级
}

// 内网Ip映射信息
type InternalIpInfos struct {
	Ip      string   `mapstructure:"ip" json:"ip"`           // Ip列表
	Domains []string `mapstructure:"domains" json:"domains"` // 逻辑域列表
}

// 直播账户
type LiveAccount struct {
	AppID    string `mapstructure:"appId" json:"appId"`
	AppToken string `mapstructure:"appToken" json:"appToken"`
	URL      string `mapstructure:"url" json:"url"`
}

// 全部配置
type SystemConfig struct {
	DataBase            DataBase        `mapstructure:"database"`
	DataBaseReplay      DataBase        `mapstructure:"database_replay"`
	DataBaseReplaySlave DataBase        `mapstructure:"database_replay_slave"`
	DataBaseSlave       DataBase        `mapstructure:"database_slave"`
	Mongodb             DataBase        `mapstructure:"mongodb"`
	MongodbSlave        DataBase        `mapstructure:"mongodb_slave"`
	IMNoSql             NoSql           `mapstructure:"im_redis"`
	NoSqlSlave          NoSql           `mapstructure:"redis_slave"`
	UserRedis           NoSql           `mapstructure:"redis_user_status"`
	IMRedis             NoSql           `mapstructure:"im_redis"`
	NoSql               NoSql           `mapstructure:"redis"`
	Kafka               Kafka           `mapstructure:"kafka"`
	KafkaData           Kafka           `mapstructure:"kafka_data"`
	KafkaStream         Kafka           `mapstructure:"kafka_stream"`
	KafkaES             Kafka           `mapstructure:"kafka_es"`
	KafkaAsyncQueue     KafkaAsyncQueue `mapstructure:"kafka_async_queue"`
	Statsd              Statsd          `mapstructure:"statsd"`
	GroupFusing         Fusing          `mapstructure:"fusing_type"`
	LiveAccount         LiveAccount     `mapstructure:"live_account"`
	H5Site              H5Site          `mapstructure:"h5_site"`
}

// H5互动课件
type H5Site struct {
	Sites                   []string `mapstructure:"sites" json:"sites"`
	HostUrl                 string   `mapstructure:"hostUrl" json:"hostUrl" `
	ControllerPCName        string   `mapstructure:"controllerPCName" json:"controllerPCName"`
	ControllerPCVersion     string   `mapstructure:"controllerPCVersion" json:"controllerPCVersion"`
	ControllerMobileName    string   `mapstructure:"controllerMobileName" json:"controllerMobileName"`
	ControllerMobileVersion string   `mapstructure:"controllerMobileVersion" json:"controllerMobileVersion"`
}

// CDN配置
type CDNConfig struct {
	SecretID  string `mapstructure:"secret_id" json:"secret_id"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key"`
	Endpoint  string `mapstructure:"endpoint" json:"endpoint"`
	Region    string `mapstructure:"region" json:"region"`
	Procedure string `mapstructure:"procedure" json:"procedure"`
	SubAppID  int    `mapstructure:"sub_app_id" json:"sub_app_id"`
}

// 直播
type LiveCDNConfig struct {
	CDNProviderName    string `mapstructure:"cdn_provider_name" json:"cdn_provider_name"`       //	CDN厂商名字
	EnableLive         bool   `mapstructure:"enable_live" json:"enable_live"`                   // 	使能节点
	PushAuthKey        string `mapstructure:"push_auth_key" json:"push_auth_key"`               //	推流密钥
	PullAuthKey        string `mapstructure:"pull_auth_key" json:"pull_auth_key"`               //	拉流密钥
	PullMediaType      string `mapstructure:"pull_media_type" json:"live_media_type"`           //	拉流类型
	LiveExpiredHour    int    `mapstructure:"live_expired_hour" json:"live_expired_hour"`       //	鉴权超时时间
	PullDomain         string `mapstructure:"pull_domain" json:"pull_domain"`                   //  拉流域名
	PushDomain         string `mapstructure:"push_domain" json:"push_domain"`                   //  推流域名
	ReplayOffline      bool   `mapstructure:"replay_offline" json:"replay_offline"`             //  回放离线模式开关
	ReplayOnline       bool   `mapstructure:"replay_online" json:"replay_online"`               //  回放在线模式开关
	ReplayOnlineMuxer  string `mapstructure:"replay_online_muxer" json:"replay_online_muxer"`   //  回放在线格式
	ReplayOfflineMuxer string `mapstructure:"replay_offline_muxer" json:"replay_offline_muxer"` //  回放离线格式
	Weight             int    `mapstructure:"weight" json:"weight"`                             //  厂商CDN权重
}

// 点播
type VodCDNConfig struct {
	SecretID        string `mapstructure:"secret_id" json:"secret_id"`                 // Ak
	SecretKey       string `mapstructure:"secret_key" json:"secret_key"`               // Sk
	Endpoint        string `mapstructure:"endpoint" json:"endpoint"`                   // CDN地址
	Region          string `mapstructure:"region" json:"region"`                       // CDN区域
	Procedure       string `mapstructure:"procedure" json:"procedure"`                 // 工作流名称
	SubAppID        string `mapstructure:"sub_app_id" json:"sub_app_id"`               // 子应用ID
	CDNProviderName string `mapstructure:"cdn_provider_name" json:"cdn_provider_name"` // CDN厂商名字
	Weight          int    `mapstructure:"weight" json:"weight"`                       // 厂商CDN权重
}

// 点播和直播中的CDN权重
type CDNWeight struct {
	Ali int `mapstructure:"ali" json:"ali"` // 阿里
	Tx  int `mapstructure:"tx" json:"tx"`   // 腾讯
	Ws  int `mapstructure:"ws" json:"ws"`   // 网宿
	Jd  int `mapstructure:"jd" json:"jd"`   // 京东
}

// ssoid对应服务信息
type SSOIdService struct {
	Enable bool               `mapstructure:"enable" json:"enable"`
	Ans    []SSoIdServiceInfo `mapstructure:"ans" json:"ans"`
	Mcs    []SSoIdServiceInfo `mapstructure:"mcs" json:"mcs"`
}

type SSoIdServiceInfo struct {
	SSOIds    []int  `mapstructure:"sso_ids" json:"ssoids"`
	ServiceID string `mapstructure:"service_id" json:"service_id"`
}
