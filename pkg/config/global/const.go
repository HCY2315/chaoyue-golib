package global

// key的数据类型
type (
	GlobalKey      string // 全局配置key类型
	MediaKey       string // 媒体配置key类型
	ApolloKey      string // Apollo key类型
	ConfigKey      string // Config Key类型
	ServiceReqType int    // 服务请求类型
)

// CDN厂商名字
const (
	CDNProvider_AliYun  string = "ali" // 阿里
	CDNProvider_Tencent string = "tx"  // 腾讯
	CDNProvider_WangSu  string = "ws"  // 网宿
	CDNProvider_JD      string = "jd"  // 京东
)

// config/global下的配置key
const (
	GlobalKeyName                = ""
	GlobalKeyDatabaseBase        = GlobalKeyName + "database_base"
	GlobalKeyDatabaseBaseSlave   = GlobalKeyName + "database_base_slave"
	GlobalKeyDatabase            = GlobalKeyName + "database"
	GlobalKeyDatabaseSlave       = GlobalKeyName + "database_slave"
	GlobalKeyDatabaseReplay      = GlobalKeyName + "database_replay"       // replay专用数据库
	GlobalKeyDatabaseReplaySlave = GlobalKeyName + "database_replay_slave" // replay专用数据库
	GlobalKeyImRedis             = GlobalKeyName + "im_redis"
	GlobalKeyKafka               = GlobalKeyName + "kafka"
	GlobalKeyKafkaData           = GlobalKeyName + "kafka_data"
	GlobalKeyKafkaStream         = GlobalKeyName + "kafka_stream"
	GlobalKeyKafkaES             = GlobalKeyName + "kafka_es" // kafka => es
	GlobalKeyMongodb             = GlobalKeyName + "mongodb"
	GlobalKeyMongodbSlave        = GlobalKeyName + "mongodb_slave"
	GlobalKeyRedis               = GlobalKeyName + "redis"
	GlobalKeyRedisSlave          = GlobalKeyName + "redis_slave"
	GlobalKeyStatsd              = GlobalKeyName + "statsd"
	GlobalKeyFusingType          = GlobalKeyName + "fusing_type"
	GlobalKeyLiveAccount         = GlobalKeyName + "live_account"
	GlobalKeyH5Site              = GlobalKeyName + "h5_site"
	GlobalKeyRegion              = GlobalKeyName + "region"
	GlobalKeyInternalIps         = GlobalKeyName + "internal_ip_info"     // 内网IP列表
	GlobalKeyAlarmUrl            = GlobalKeyName + "alarm_url"            // 报警地址
	GlobalKeyInfluxDb            = GlobalKeyName + "influxdb"             // Influxdb地址
	GlobalKeyRocketMQ            = GlobalKeyName + "rocketmq"             // rocketmq地址 【deprecated】
	GlobalKeyRedisUserStatus     = GlobalKeyName + "redis_user_status"    // 存储用户session'状态redis
	GlobalKeyLivePlatform        = GlobalKeyName + "live_platform"        // 直播平台地址
	GlobalKeyKafkaAsyncQueue     = GlobalKeyName + "kafka_async_queue"    // kafka异步队列的信息：brokers,topics,groups
	GlobalKeyTraceAgentAddr      = GlobalKeyName + "trace_agent_addr"     // jaeger agent 地址
	GlobalKeyCDN                 = GlobalKeyName + "cdn"                  // cdn 配置
	GlobalKeySSOIdService        = GlobalKeyName + "ssoid_service"        // ssoid与service节点信息配置
	GlobalKeySSOIdServiceSwitch  = GlobalKeyName + "ssoid_service_switch" // ssoid与service节点信息配置开关

	// media
	GlobalKey_MediaCDNLive = GlobalKeyName + "media_cdn_live"
	GlobalKey_MediaCDNVod  = GlobalKeyName + "media_cdn_vod"
	//global_media_cdn_last_modify_time

)
