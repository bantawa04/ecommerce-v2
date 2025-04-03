package implementations

import (
	"context"

	"beautyessentials.com/internal/repository/interfaces"
)

// HealthRepositoryImpl implements the HealthRepository interface
type HealthRepositoryImpl struct{}

// NewHealthRepository creates a new instance of HealthRepositoryImpl
func NewHealthRepository() interfaces.HealthRepository {
	return &HealthRepositoryImpl{}
}

// CheckHealth checks if the repository layer is healthy
func (r *HealthRepositoryImpl) CheckHealth(ctx context.Context) (bool, error) {
	// In a real application, this might check database connectivity
	return true, nil
}