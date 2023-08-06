package repository

import (
	"fmt"

	"go.uber.org/zap"
)

type BlockVolume struct {
	Time         string  `db:"time" json:"time"`
	Transactions float64 `db:"transactions" json:"transactions,omitempty"`
	BlockWeight  float64 `db:"block_weight" json:"blockWeight,omitempty"`
	MiningFee    float64 `db:"mining_fee" json:"miningFee"`
}

func (r *PGSQLRepository) BlockVolumeTxXMiningFee(days int) ([]BlockVolume, error) {
	blocks := []BlockVolume{}

	err := r.db.Select(&blocks, fmt.Sprintf(`
	SELECT
		bucket as "time",
		avg(tx_count) AS transactions,
		avg(block_fee_sat)*0.00000001 AS "mining_fee"
 	FROM one_hour_blocks
 	WHERE bucket > now() - INTERVAL '%v day'
 	GROUP BY bucket
 	ORDER BY bucket; 
	`, days))
	if err != nil {
		r.logger.Error("failed to execute BlockVolumeTxXMiningFee", zap.Error(err))
		return nil, err
	}

	return blocks, nil
}

func (r *PGSQLRepository) BlockVolumeXMiningFee(days int) ([]BlockVolume, error) {
	blocks := []BlockVolume{}

	err := r.db.Select(&blocks, fmt.Sprintf(`
	SELECT
		bucket as "time",
		avg(block_weight) AS "block_weight",
		avg(block_fee_sat)*0.00000001 AS "mining_fee"
 	FROM one_hour_blocks
 	WHERE bucket > now() - INTERVAL '%v day'
 	GROUP BY bucket
 	ORDER BY bucket; 
	`, days))
	if err != nil {
		r.logger.Error("failed to execute BlockVolumeXMiningFee", zap.Error(err))
		return nil, err
	}

	return blocks, nil
}
