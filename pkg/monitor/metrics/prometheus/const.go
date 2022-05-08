/*
 * @File   : const
 * @Author : huangbin
 *
 * @Created on 2020/8/27 6:41 下午
 * @Project : prometheus
 * @Software: GoLand
 * @Description  :
 */

package prometheus

//go:generate enumer -type=MetricConstType -trimprefix=MetricConstType

type MetricConstType int

const (
	CSTDefaultRouter = "/metrics"
)

const (
	CSTMetricNone         MetricConstType = iota // 无效
	CSTMetricGauge                               // 仪表盘
	CSTMetricGaugeVec                            // 仪表盘
	CSTMetricCounter                             // 计数器
	CSTMetricCounterVec                          // 计数器
	CSTMetricHistogram                           // 直方图
	CSTMetricHistogramVec                        // 直方图
	CSTMetricSummary                             // 摘要
	CSTMetricSummaryVec                          // 摘要
	CSTMetricMax                                 // 摘要
)
