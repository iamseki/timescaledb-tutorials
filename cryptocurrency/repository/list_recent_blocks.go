package repository

import (
	"fmt"

	"go.uber.org/zap"
)

type RecentBlock struct {
	BlockID          int     `db:"block_id" json:"blockID"`
	TransactionCount int     `db:"transaction_count" json:"transactionCount"`
	BlockWeight      int     `db:"block_weight" json:"blockWeight"`
	BlockValueUSD    float64 `db:"block_value_usd" json:"blockValueUSD"`
}

func (r *PGSQLRepository) ListRecentBlocks(limit int) ([]RecentBlock, error) {
	recentBlocks := []RecentBlock{}

	err := r.db.Select(&recentBlocks, fmt.Sprintf(`
	WITH recent_blocks as (
		SELECT block_id FROM transactions
		WHERE is_coinbase IS TRUE
		ORDER BY time DESC
		LIMIT %v
	)
	SELECT 
		t.block_id,
		COUNT(*) as transaction_count,
		SUM(t.weight) as block_weight,
		SUM(t.output_total_usd) as block_value_usd
	FROM transactions t
	INNER JOIN recent_blocks rb on rb.block_id = t.block_id
	WHERE t.is_coinbase IS NOT TRUE
	GROUP BY t.block_id
	ORDER BY t.block_id asc;
	`, limit))
	if err != nil {
		r.logger.Error("Error on list recent blocks query", zap.Error(err))
		return nil, err
	}

	return recentBlocks, nil
}
