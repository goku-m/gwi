package service

import (
	"github.com/goku-m/gwi/internal/repository"
	"github.com/goku-m/gwi/internal/server"
)

type Services struct {
	Auth   *AuthService
	Farmer *FarmerService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	// s.Job.SetAuthService(authService)

	// awsClient, err := aws.NewAWS(s)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create AWS client: %w", err)
	// }

	return &Services{
		Auth:   authService,
		Farmer: NewFarmerService(s, repos.Farmer),
	}, nil
}
