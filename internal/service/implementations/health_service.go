package implementations

import (
	"context"

	"beautyessentials.com/internal/repository/interfaces"
	serviceInterfaces "beautyessentials.com/internal/service/interfaces"
)

// HealthServiceImpl implements the HealthService interface
type HealthServiceImpl struct {
	healthRepo interfaces.HealthRepository
}

// NewHealthService creates a new instance of HealthServiceImpl
func NewHealthService(healthRepo interfaces.HealthRepository) serviceInterfaces.HealthService {
	return &HealthServiceImpl{
		healthRepo: healthRepo,
	}
}

// CheckHealth checks if the service is healthy
func (s *HealthServiceImpl) CheckHealth(ctx context.Context) (bool, error) {
	// Call the repository layer
	return s.healthRepo.CheckHealth(ctx)
}