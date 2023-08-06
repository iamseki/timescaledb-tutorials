package repository

import (
	"fmt"

	"go.uber.org/zap"
)

type TxVolume struct {
	Time              string  `db:"time" json:"time"`
	TransactionVolume int     `db:"transaction_volume" json:"transactionVolume"`
	Fees              float64 `db:"fees" json:"fees,omitempty"`
	BTCUSDRate        float64 `db:"btc_usd_rate" json:"BTCUSDRate"`
}

func (r *PGSQLRepository) TxVolumeXFeesByHour(days int) ([]TxVolume, error) {
	txVolumeByHour := []TxVolume{}

	err := r.db.Select(&txVolumeByHour, fmt.Sprintf(`
	SELECT
		bucket as "time",
		tx_count as "transaction_volume",
		average(stats_fee_sat) as fees
	FROM one_hour_transactions
	WHERE bucket > NOW() - INTERVAL '%v days'
	ORDER BY bucket;
	`, days))
	if err != nil {
		r.logger.Error("Failed to read TxVolumeXFeesByHour", zap.Error(err))
		return nil, err
	}

	return txVolumeByHour, nil
}

func (r *PGSQLRepository) TxVolumeXUSDByHour(days int) ([]TxVolume, error) {
	txVolumeByHour := []TxVolume{}

	err := r.db.Select(&txVolumeByHour, fmt.Sprintf(`
	SELECT 
		bucket as "time",
		tx_count as "transaction_volume",
		total_fee_usd / (total_fee_sat * 0.00000001) as "btc_usd_rate"
	FROM one_hour_transactions
	WHERE bucket > now() - interval '%v days'
	ORDER BY bucket;
	`, days))
	if err != nil {
		r.logger.Error("Failed to read TxVolumeXUSDByHour", zap.Error(err))
		return nil, err
	}

	return txVolumeByHour, nil
}
