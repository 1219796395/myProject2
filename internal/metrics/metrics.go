package metrics

/*
import (
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// answer stats
	_metricAnswerRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "biz",
		Subsystem: "answer",
		Name:      "code_total",
		Help:      "The total number of answers",
	}, []string{"surveyId", "code"}) // code: 1 valid; 2 invalid;

	// reward stats
	_metricRewardRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "biz",
		Subsystem: "reward",
		Name:      "code_total",
		Help:      "The total number of rewards",
	}, []string{"surveyId", "code"}) // code: 1 success; 2 fail;
)

func init() {
	prometheus.MustRegister(_metricAnswerRequests, _metricRewardRequests)
}

var (
	AnswerCounter metrics.Counter
	RewardCounter metrics.Counter
)

func init() {
	AnswerCounter = prom.NewCounter(_metricAnswerRequests)
	RewardCounter = prom.NewCounter(_metricRewardRequests)
}

// // example
// func Example(surveyId int64, code string) {
// 	AnswerCounter.With(strconv.Itoa(int(surveyId)), code).Inc()
// }
*/
