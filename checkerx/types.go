package checkerx

// HealthChecker represent service check operations
type HealthChecker interface {
	// CheckHealth
	CheckHealth() error

	// Stop health check loop
	Stop() error

	// Wait until service is ready
	Wait() error
}
