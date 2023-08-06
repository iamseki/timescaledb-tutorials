package repository

import (
	"fmt"

	"go.uber.org/zap"
)

type MinerRevenue struct {
	Time   string  `db:"time" json:"time"`
	Fees   float64 `db:"fees" json:"fees"`
	Reward float64 `db:"reward" json:"reward"`
}

func (r *PGSQLRepository) MinerReveueByHour(days int) ([]MinerRevenue, error) {
	revenues := []MinerRevenue{}

	err := r.db.Select(&revenues, fmt.Sprintf(`
	WITH coinbase AS (
		SELECT 
			block_id, 
			output_total as coinbase_tx
		FROM transactions t 
		WHERE t.is_coinbase IS TRUE AND time > NOW() - INTERVAL '%v days'
	)
	SELECT
		bucket as "time",
		AVG(block_fee_sat)*0.00000001 AS "fees",
		FIRST((c.coinbase_tx - block_fee_sat), bucket)*0.00000001 AS "reward"
	FROM one_hour_blocks b
	INNER JOIN coinbase c on c.block_id = b.block_id
	GROUP BY bucket
	ORDER BY bucket;
		`, days))
	if err != nil {
		r.logger.Error("Failed to read MinerRevenueByHour", zap.Error(err))
		return nil, err
	}

	return revenues, nil
}
