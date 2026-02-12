package job

import (
	"github.com/goku-m/gwi/internal/config"
	"github.com/rs/zerolog"
)

type JobService struct {
	logger *zerolog.Logger
}

func NewJobService(logger *zerolog.Logger, cfg *config.Config) *JobService {
	_ = cfg
	return &JobService{logger: logger}
}

func (j *JobService) Start() error {
	return nil
}

func (j *JobService) Stop() {}
