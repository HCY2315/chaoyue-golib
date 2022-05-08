/*
 * @File   : struct
 * @Author : huangbin
 *
 * @Created on 2020/8/28 10:18 上午
 * @Project : prometheus
 * @Software: GoLand
 * @Description  :
 */

package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

//prometheus
//////////////////////////////////////////////////////////////////////////

type MetricCache struct {
	Enable bool          // 是否定时清除指标缓存
	Point  time.Time     // 上次清除时间点
	Period time.Duration // 清除周期,如果不设置，则默认按照一小时清除一次
}

type MetricNode struct {
	MetricType   MetricConstType
	MetricName   string
	Metric       prometheus.Collector
	ClearCache   MetricCache
	HaveRegister bool
}

// consul
//////////////////////////////////////////////////////////////////////////

type Node struct {
	Id      string          `json:"ID"`      // 服务ID
	Name    string          `json:"Name"`    // 服务名字
	Address string          `json:"Address"` // 服务注册到consul到IP，用于服务发现
	Port    int             `json:"Port"`    // 服务注册到consul的port，用于服务发现
	Tags    []string        `json:"Tags"`    // 服务的tag，自定义，可以根据这个tag来区分同一个服务名的服务
	Checks  HealthCheckInfo `json:"Check"`   // 健康检查
}
type HealthCheckInfo struct {
	DeregisterCriticalServiceAfter string `json:"DeregisterCriticalServiceAfter"` // 服务停止后多长时间后销毁，例：60m
	Http                           string `json:"HTTP"`                           // 指定健康检查的URL，调用后只要返回20X，consul都认为是健康的
	Interval                       string `json:"Interval"`                       // 健康检查间隔时间，每隔10s，调用一次上面的URL
}
