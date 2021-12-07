package metrics

type RequestInfo struct {
	apiName      string
	responseTime float64
	timestamp    float64
}

func NewRequestInfo(apiName string, responseTime, timestamp float64) *RequestInfo {
	return &RequestInfo{}
}

func (info *RequestInfo) GetApiName() string {
	return info.apiName
}

func (info *RequestInfo) GetResponseTime() float64 {
	return info.responseTime
}

func (info *RequestInfo) GetTimestamp() float64 {
	return info.timestamp
}

type MetricsCollector struct {
	metricsStorage MetricsStorage
}

func (m *MetricsCollector) recordRequest(requestInfo *RequestInfo) {
	if requestInfo == nil || requestInfo.GetApiName() == "" {
		return
	}

	m.metricsStorage.SaveRequestInfo(requestInfo)
}
