package repository

import (
	"fmt"

	"go.uber.org/zap"
)

type MiningReward struct {
	Time   string  `db:"time" json:"time"`
	Fees   float64 `db:"fees" json:"fees"`
	Reward float64 `db:"reward" json:"reward"`
}

func (r *PGSQLRepository) MinersRewardByHour(days int) ([]MiningReward, error) {
	revenues := []MiningReward{}

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
		r.logger.Error("Failed to read MinerRevenueBTCByHour", zap.Error(err))
		return nil, err
	}

	return revenues, nil
}

type MiningFee struct {
	Time        string  `db:"time" json:"time"`
	BlockWeight float64 `db:"block_weight" json:"blockWeight"`
	MiningFee   float64 `db:"mining_fee" json:"miningFee"`
}

func (r *PGSQLRepository) MiningBlockWeightXFeeByHour(days int) ([]MiningFee, error) {
	mining := []MiningFee{}

	err := r.db.Select(&mining, fmt.Sprintf(`
	WITH stats AS (
		SELECT
				bucket,
				stats_agg(block_weight, block_fee_sat) AS block_stats --- https://docs.timescale.com/api/latest/hyperfunctions/statistical-and-regression-analysis/stats_agg-two-variables/
		FROM one_hour_blocks
		WHERE bucket > NOW() - INTERVAL '%v days'
		GROUP BY bucket
 )
 	SELECT
			bucket as "time",
			average_y(rolling(block_stats) OVER (ORDER BY bucket RANGE '12 hours' PRECEDING)) AS "block_weight",
			average_x(rolling(block_stats) OVER (ORDER BY bucket RANGE '12 hours' PRECEDING))*0.00000001 AS "mining_fee"
 	FROM stats
 	ORDER BY bucket;
		`, days))
	if err != nil {
		r.logger.Error("Failed to read MiningBlockWeightXFeeByHour", zap.Error(err))
		return nil, err
	}

	return mining, nil
}

type MiningRevenue struct {
	Time         string  `db:"time" json:"time"`
	RevenueInBTC float64 `db:"revenue_in_btc" json:"revenueInBTC"`
	RevenueInUSD float64 `db:"revenue_in_usd" json:"revenueInUSD"`
}

func (r *PGSQLRepository) MiningRevenueByHour(days int) ([]MiningRevenue, error) {
	revenues := []MiningRevenue{}

	err := r.db.Select(&revenues, fmt.Sprintf(`
	SELECT
	bucket as "time",
	average_y(rolling(stats_miner_revenue) OVER (ORDER BY bucket RANGE '12 hours' PRECEDING))*0.00000001 AS "revenue_in_btc",
	average_x(rolling(stats_miner_revenue) OVER (ORDER BY bucket RANGE '12 hours' PRECEDING)) AS "revenue_in_usd"
	FROM one_hour_coinbase
	WHERE bucket > NOW() - INTERVAL '%v days'
	ORDER BY 1;
	`, days))
	if err != nil {
		r.logger.Error("Failed to read MiningRevenueByHour", zap.Error(err))
		return nil, err
	}

	return revenues, nil
}
