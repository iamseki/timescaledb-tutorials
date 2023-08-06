package repository

import (
	"fmt"

	"go.uber.org/zap"
)

func (r *PGSQLRepository) ListRecentBlocks(limit int) ([]RecentBlock, error) {
	recentBlocks := []RecentBlock{}

	err := r.db.Select(&recentBlocks, fmt.Sprintf(`
	SELECT time, hash, block_id, fee_usd  FROM transactions
	WHERE is_coinbase IS NOT TRUE
	ORDER BY time DESC
	LIMIT %v;	
	`, limit))
	if err != nil {
		r.logger.Error("Error on list recent blocks query", zap.Error(err))
		return nil, err
	}

	return recentBlocks, nil
}
