package core

import (
	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/repository"
	"go.uber.org/zap"
)

type Core struct {
	logger     *zap.Logger
	repository repository.Repository
}

func New(logger *zap.Logger, repository repository.Repository) *Core {
	return &Core{logger, repository}
}
