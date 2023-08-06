package core

import (
	"errors"

	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/repository"
	"go.uber.org/zap"
)

func (c *Core) MinersRevenueByHour(days int) ([]repository.MinerRevenue, error) {
	if days == 0 {
		days = 5
		c.logger.Info("Using default last days on MinersRevenueByHour", zap.Int("days", days))
	}

	res, err := c.repository.MinerReveueByHour(days)
	if err != nil {
		c.logger.Error("repository error on MinerReveueByHour", zap.Error(err))
		return nil, errors.New("repository error on MinerReveueByHour")
	}

	return res, nil
}
