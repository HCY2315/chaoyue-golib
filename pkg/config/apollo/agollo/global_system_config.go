package agollo

import (
	"ccp-server/framework/config/global"
	"github.com/spf13/viper"
)

type GlobalSystemConfigKey string

const (
	GlobalKeyDatabase            GlobalSystemConfigKey = "database"
	GlobalKeyDatabaseSlave                             = "database_slave"
	GlobalKeyDatabaseBase                              = "database_base"
	GlobalKeyDatabaseBaseSlave                         = "database_base_slave"
	GlobalKeyDatabaseReplay                            = "database_replay"       // replay专用数据库
	GlobalKeyDatabaseReplaySlave                       = "database_replay_slave" // replay专用数据库
	GlobalKeyImRedis                                   = "im_redis"
	GlobalKeyKafka                                     = "kafka"
	GlobalKeyKafkaData                                 = "kafka_data"
	GlobalKeyKafkaStream                               = "kafka_stream"
	GlobalKeyKafkaES                                   = "kafka_es" // kafka => es
	GlobalKeyMongodb                                   = "mongodb"
	GlobalKeyMongodbSlave                              = "mongodb_slave"
	GlobalKeyRedis                                     = "redis"
	GlobalKeyRedisSlave                                = "redis_slave"
	GlobalKeyStatsd                                    = "statsd"
	GlobalKeyFusingType                                = "fusing_type"
	GlobalKeyLiveAccount                               = "live_account"
	GlobalKeyH5Site                                    = "h5_site"
	GlobalKeyRegion                                    = "region"
	GlobalKeyInternalIps                               = "internal_ip_info"     // 内网IP列表
	GlobalKeyAlarmUrl                                  = "alarm_url"            // 报警地址
	GlobalKeyInfluxDb                                  = "influxdb"             // Influxdb地址
	GlobalKeyRocketMQ                                  = "rocketmq"             // rocketmq地址 【deprecated】
	GlobalKeyRedisUserStatus                           = "redis_user_status"    // 存储用户session'状态redis
	GlobalKeyLivePlatform                              = "live_platform"        // 直播平台地址
	GlobalKeyKafkaAsyncQueue                           = "kafka_async_queue"    // kafka异步队列的信息：brokers,topics,groups
	GlobalKeyTraceAgentAddr                            = "trace_agent_addr"     // jaeger agent 地址
	GlobalKeyCDN                                       = "cdn"                  // cdn 配置
	GlobalKeySSOIdService                              = "ssoid_service"        // ssoid与service节点信息配置
	GlobalKeySSOIdServiceSwitch                        = "ssoid_service_switch" // ssoid与service节点信息配置开关
)

func GetSystemConfig(systemConfigNamespace string) global.SystemConfig {
	var systemConfig global.SystemConfig
	nodeMap := map[GlobalSystemConfigKey]interface{}{
		GlobalKeyDatabase:            &systemConfig.DataBase,
		GlobalKeyDatabaseReplay:      &systemConfig.DataBaseReplay,
		GlobalKeyDatabaseReplaySlave: &systemConfig.DataBaseReplaySlave,
		GlobalKeyDatabaseSlave:       &systemConfig.DataBaseSlave,
		GlobalKeyImRedis:             &systemConfig.IMRedis,
		GlobalKeyKafka:               &systemConfig.Kafka,
		GlobalKeyKafkaData:           &systemConfig.KafkaData,
		GlobalKeyKafkaStream:         &systemConfig.KafkaStream,
		GlobalKeyKafkaES:             &systemConfig.KafkaES,
		GlobalKeyKafkaAsyncQueue:     &systemConfig.KafkaAsyncQueue,
		GlobalKeyMongodb:             &systemConfig.Mongodb,
		GlobalKeyMongodbSlave:        &systemConfig.MongodbSlave,
		GlobalKeyRedis:               &systemConfig.NoSql,
		GlobalKeyRedisSlave:          &systemConfig.NoSqlSlave,
		GlobalKeyRedisUserStatus:     &systemConfig.UserRedis,
		GlobalKeyStatsd:              &systemConfig.Statsd,
		GlobalKeyFusingType:          &systemConfig.GroupFusing,
		GlobalKeyLiveAccount:         &systemConfig.LiveAccount,
		GlobalKeyH5Site:              &systemConfig.H5Site,
	}
	for _, value := range nodeMap {
		if errUnmarshal := viper.UnmarshalKey(systemConfigNamespace, value); errUnmarshal != nil {

		}
	}
	return systemConfig
}
