package metrics

type MetricsStorage interface {
	SaveRequestInfo(requestInfo *RequestInfo)
	GetRequestInfos(apiName string, startTimeInMillis, endTimeInMillis int64) []*RequestInfo
	GetRequestInfosToMap(startTimeInMillis, endTimeInMillis int64) map[string][]*RequestInfo
}

type RedisMetricsStorage struct {
}

func NewRedisMetricsStorage() *RedisMetricsStorage {
	return &RedisMetricsStorage{}
}

func (s *RedisMetricsStorage) SaveRequestInfo(requestInfo *RequestInfo) {
	return
}
func (s *RedisMetricsStorage) GetRequestInfos(apiName string, startTimeInMillis, endTimeInMillis int64) []*RequestInfo {
	return nil
}
func (s *RedisMetricsStorage) GetRequestInfosToMap(startTimeInMillis, endTimeInMillis int64) map[string][]*RequestInfo {
	return map[string][]*RequestInfo{
		"register": []*RequestInfo{
			&RequestInfo{
				"register",
				123.0,
				10234.0,
			},
			&RequestInfo{
				"register",
				223.0,
				11234.0,
			},
			&RequestInfo{
				"register",
				323.0,
				12334.0,
			},
		},
		"login": []*RequestInfo{
			&RequestInfo{
				"login",
				23.0,
				12434.0,
			},
			&RequestInfo{
				"login",
				1223.0,
				14234.0,
			},
		},
	}
}
