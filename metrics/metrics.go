package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// CheckTotal 巡检总次数
	CheckTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "healthcheck_total",
		Help: "巡检执行总次数",
	}, []string{"url"})

	// CheckErrors 错误次数
	CheckErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "healthcheck_errors_total",
		Help: "巡检失败总次数",
	}, []string{"url"})

	// CheckDuration 响应时间（Histogram，支持分位数查询）
	CheckDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "healthcheck_duration_ms",
		Help:    "接口响应时间（毫秒）",
		Buckets: []float64{50, 100, 200, 500, 1000, 2000, 5000},
	}, []string{"url"})

	// CheckStatusCode HTTP 状态码分布
	CheckStatusCode = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "healthcheck_status_code_total",
		Help: "各 HTTP 状态码出现次数",
	}, []string{"url", "status_code"})
)
