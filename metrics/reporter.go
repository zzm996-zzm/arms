package metrics

import (
	"encoding/json"
	"fmt"
	"time"
)

type Reporter interface {
	StartRepeatedReport(durationInSeconds int64)
}

type ConsoleReporter struct {
	metricsStorage MetricsStorage
}

func NewConsoleReporter(storage MetricsStorage) *ConsoleReporter {
	return &ConsoleReporter{metricsStorage: storage}
}
func (c *ConsoleReporter) StartRepeatedReport(durationInSeconds int64) {
	//从数据库拉取数据
	durationInMillis := durationInSeconds * 1000
	endTimeInMillis := time.Now().UnixNano() / 1e6
	startTimeInMillis := endTimeInMillis - durationInMillis

	requestInfos := c.metricsStorage.GetRequestInfosToMap(startTimeInMillis, endTimeInMillis)
	stats := make(map[string]*RequestStat)
	for apiName, v := range requestInfos {
		requestStat := Aggregator(v, durationInMillis)
		stats[apiName] = requestStat
	}

	fmt.Println("Time Span: [", startTimeInMillis, ", ", endTimeInMillis, "]")
	res, err := json.Marshal(stats)
	if err != nil {
		fmt.Println("json error:", err)
	}

	fmt.Println(string(res))
}
