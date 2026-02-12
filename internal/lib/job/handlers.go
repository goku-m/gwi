package job

import (
	"github.com/goku-m/gwi/internal/config"
	"github.com/goku-m/gwi/internal/lib/email"
	"github.com/rs/zerolog"
)

var emailClient *email.Client

func (j *JobService) InitHandlers(config *config.Config, logger *zerolog.Logger) {
	emailClient = email.NewClient(config, logger)
}
