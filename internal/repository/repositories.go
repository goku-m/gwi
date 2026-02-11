package repository

import "github.com/goku-m/gwi/internal/server"

type Repositories struct {
	Farmer *FarmerRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Farmer: NewFarmerRepository(s),
	}
}
