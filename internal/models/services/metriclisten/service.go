package metriclisten

type AgentMetric struct {
	MetricType  string
	MetricName  string
	MetricValue float64
}

type AgentMetrics []AgentMetric

type Service interface {
	GetAgentMetrics() AgentMetrics
}
