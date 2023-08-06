package core

import (
	"errors"

	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/repository"
	"go.uber.org/zap"
)

func (c *Core) TxVolumeXFeesByHour(days int) ([]repository.TxVolume, error) {
	if days == 0 {
		days = 5
		c.logger.Info("Using default last days on TxVolumeXFeesByHour", zap.Int("days", days))
	}

	res, err := c.repository.TxVolumeXFeesByHour(days)
	if err != nil {
		c.logger.Error("repository error on TxVolumeXFeesByHour", zap.Error(err))
		return nil, errors.New("repository error on TxVolumeXFeesByHour")
	}

	return res, nil
}

func (c *Core) TxVolumeXUSDByHour(days int) ([]repository.TxVolume, error) {
	if days == 0 {
		days = 5
		c.logger.Info("Using default last days on TxVolumeXBTCUSDRateByHour", zap.Int("days", days))
	}

	res, err := c.repository.TxVolumeXUSDByHour(days)
	if err != nil {
		c.logger.Error("repository error on TxVolumeXBTCUSDRateByHour", zap.Error(err))
		return nil, errors.New("repository error on TxVolumeXBTCUSDRateByHour")
	}

	return res, nil
}
