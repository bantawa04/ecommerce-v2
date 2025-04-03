package interfaces

import "context"

// HealthRepository defines the interface for health check repository operations
type HealthRepository interface {
	CheckHealth(ctx context.Context) (bool, error)
}