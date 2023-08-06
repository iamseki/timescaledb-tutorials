package core

import (
	"errors"

	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/repository"
	"go.uber.org/zap"
)

func (c *Core) BlockVolumeTxXMiningFee(days int) ([]repository.BlockVolume, error) {
	if days == 0 {
		days = 5
		c.logger.Info("Using default last days on BlockVolumeTxXMiningFee", zap.Int("days", days))
	}

	res, err := c.repository.BlockVolumeTxXMiningFee(days)
	if err != nil {
		c.logger.Error("repository error on BlockVolumeTxXMiningFee", zap.Error(err))
		return nil, errors.New("repository error on BlockVolumeTxXMiningFee")
	}

	return res, nil
}

func (c *Core) BlockVolumeXMiningFee(days int) ([]repository.BlockVolume, error) {
	if days == 0 {
		days = 5
		c.logger.Info("Using default last days on BlockVolumeXMiningFee", zap.Int("days", days))
	}

	res, err := c.repository.BlockVolumeXMiningFee(days)
	if err != nil {
		c.logger.Error("repository error on BlockVolumeXMiningFee", zap.Error(err))
		return nil, errors.New("repository error on BlockVolumeXMiningFee")
	}

	return res, nil
}
