package promethus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// monitorComponent 监控指标类型.
type monitorComponentType string

const (
	// 计数器.
	counter monitorComponentType = "counter"
	// 摘要.
	summary monitorComponentType = "summary"
	// 仪表盘.
	gauge monitorComponentType = "gauge"

	// 定时器触发记录总数.
	timerExecTotalCnt        = "timer_exec_total_cnt"
	timerExecTotalCntSummary = "定时器触发记录总数"

	// 定时器触发延时
	timerDelayCnt        = "timer_delay_cnt"
	timerDelayCntSummary = "定时器触发延时"

	// 处于激活态的定时器总数.
	timerEnabledCnt        = "timer_enabled_cnt"
	timerEnabledCntSummary = "激活态定时器总数"

	// 未触发定时器数量.
	timerUnexecedCnt        = "timer_unexeced_cnt"
	timerUnexecedCntSummary = "未按时执行的定时器数量"

	reportName        = "_name"
	reportType        = "_type"
	timerApp   string = "xtimer"

	// 通用标签.
	label = "label"
	timer = "timer"
)

// Reporter 监控上报服务.
type Reporter struct {
	timerExecRecorder     *prometheus.CounterVec
	timeDelayRecorder     prometheus.ObserverVec
	timerEnabledRecorder  *prometheus.GaugeVec
	timerUnexecedRecorder *prometheus.GaugeVec
}

var reporter = newReporter()

// GetReporter 获取单例上报服务.
func GetReporter() *Reporter {
	return reporter
}

// newReporter 监控上报服务构造器.
func newReporter() *Reporter {
	r := Reporter{
		// 定时器触发记录.
		timerExecRecorder: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: timerExecTotalCnt,
			Help: timerExecTotalCntSummary,
		}, []string{
			timerApp,
			reportName,
			reportType,
		}).MustCurryWith(prometheus.Labels{reportName: timerExecTotalCntSummary,
			reportType: string(counter)}),

		// 定时器延时记录.
		timeDelayRecorder: promauto.NewSummaryVec(prometheus.SummaryOpts{
			Name:       timerDelayCnt,
			Help:       timerDelayCntSummary,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 0.999: 0.0001, 0.9999: 0.00001},
		}, []string{
			timerApp,
			reportName,
			reportType,
		}).MustCurryWith(prometheus.Labels{reportName: timerDelayCntSummary,
			reportType: string(summary)}),

		// 处于激活态的定时器总数.
		timerEnabledRecorder: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: timerEnabledCnt,
			Help: timerEnabledCntSummary,
		}, []string{
			label,
			reportName,
			reportType,
		}).MustCurryWith(prometheus.Labels{reportName: timerEnabledCntSummary,
			reportType: string(gauge)}),

		// 未触发定时器数量.
		timerUnexecedRecorder: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: timerUnexecedCnt,
			Help: timerUnexecedCntSummary,
		}, []string{
			label,
			reportName,
			reportType,
		}).MustCurryWith(prometheus.Labels{reportName: timerUnexecedCntSummary,
			reportType: string(gauge)}),
	}

	// prometheus.MustRegister(r.timerExecRecorder, r.timeDelayRecorder, r.timerEnabledRecorder, r.timerUnexecedRecorder)
	return &r
}

func (r *Reporter) ReportExecRecord(app string) {
	r.timerExecRecorder.WithLabelValues(app).Inc()
}

func (r *Reporter) ReportTimerDelayRecord(app string, cost float64) {
	r.timeDelayRecorder.WithLabelValues(app).Observe(cost)
}

func (r *Reporter) ReportTimerEnabledRecord(total float64) {
	r.timerEnabledRecorder.WithLabelValues(timer).Set(total)
}

func (r *Reporter) ReportTimerUnexecedRecord(total float64) {
	r.timerUnexecedRecorder.WithLabelValues(timer).Set(total)
}
