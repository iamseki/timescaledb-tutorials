package core

import (
	"errors"

	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/repository"
	"go.uber.org/zap"
)

func (c *Core) MinersRewardByHour(days int) ([]repository.MiningReward, error) {
	if days == 0 {
		days = 5
		c.logger.Info("Using default last days on MinersRewardByHour", zap.Int("days", days))
	}

	res, err := c.repository.MinersRewardByHour(days)
	if err != nil {
		c.logger.Error("repository error on MinersRewardByHour", zap.Error(err))
		return nil, errors.New("repository error on MinersRewardByHour")
	}

	return res, nil
}

func (c *Core) MiningBlockWeightXFeeByHour(days int) ([]repository.MiningFee, error) {
	if days == 0 {
		days = 5
		c.logger.Info("Using default last days on MiningBlockWeightXFeeByHour", zap.Int("days", days))
	}

	res, err := c.repository.MiningBlockWeightXFeeByHour(days)
	if err != nil {
		c.logger.Error("repository error on MiningBlockWeightXFeeByHour", zap.Error(err))
		return nil, errors.New("repository error on MiningBlockWeightXFeeByHour")
	}

	return res, nil
}

func (c *Core) MiningRevenueByHour(days int) ([]repository.MiningRevenue, error) {
	if days == 0 {
		days = 5
		c.logger.Info("Using default last days on MiningBlockWeightXFeeByHour", zap.Int("days", days))
	}

	res, err := c.repository.MiningRevenueByHour(days)
	if err != nil {
		c.logger.Error("repository error on MiningRevenueByHour", zap.Error(err))
		return nil, errors.New("repository error on MiningRevenueByHour")
	}

	return res, nil
}
