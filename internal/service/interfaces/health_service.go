package interfaces

import "context"

// HealthService defines the interface for health check service operations
type HealthService interface {
	CheckHealth(ctx context.Context) (bool, error)
}