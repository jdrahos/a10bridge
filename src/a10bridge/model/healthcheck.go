package model

//HealthCheck data sctructure for health check setup
type HealthCheck struct {
	Name                      string
	Endpoint                  string
	Port                      int
	ExpectCode                string
	RetryCount                int
	RequiredConsecutivePasses int
	Interval                  int
	Timeout                   int
}
