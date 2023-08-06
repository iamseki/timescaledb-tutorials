package core

import (
	"errors"

	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/repository"
	"go.uber.org/zap"
)

func (c *Core) ListRecentBlocks(limit int) ([]repository.RecentBlock, error) {
	if limit > 1000 { // silly validation
		c.logger.Error("limit must not be greater than 1000", zap.Int("limit", limit))
		return nil, errors.New("limit must not be greater than 1000")
	}

	res, err := c.repository.ListRecentBlocks(limit)
	if err != nil {
		c.logger.Error("repository error on ListRecentBlocks", zap.Error(err))
		return nil, errors.New("repository error on ListRecentBlocks")
	}

	return res, nil
}
