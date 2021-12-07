package metrics

import "math"

type RequestStat struct {
	MaxResponseTime float64 `json:"maxResponseTime"`
	MinResponseTime float64 `json:"minResponseTime"`
	AvgResponseTime float64 `json:"avgResponseTime"`
	// p99ResponseTime  float64
	// p999ResponseTime float64
	Count int64 `json:"count"`
	Tps   int64 `json:"tps"`
}

func newRequestStat() *RequestStat {
	return &RequestStat{}
}

func Aggregator(requestInfos []*RequestInfo, durationInMillis int64) *RequestStat {
	minRespTime := math.MaxFloat64
	maxRespTime := math.SmallestNonzeroFloat64
	avgRespTime := -1.0
	// p999RespTime := -1.0
	// p99RespTime := -1.0
	sumRespTime := 0.0	
	count := 0
	tps := 0


	for _, requestInfo := range requestInfos {
		count++
		respTime := requestInfo.GetResponseTime()
		if maxRespTime < respTime {
			maxRespTime = respTime
		}

		if minRespTime > respTime {
			minRespTime = respTime
		}

		sumRespTime += respTime

		if count != 0 {
			avgRespTime = sumRespTime / float64(count)
		}

		tps = count / int(durationInMillis) * 1000
	}

	requestStat := newRequestStat()
	requestStat.Count = int64(count)
	requestStat.Tps = int64(tps)
	requestStat.MaxResponseTime = maxRespTime
	requestStat.MinResponseTime = minRespTime
	requestStat.AvgResponseTime = avgRespTime

	return requestStat
}
