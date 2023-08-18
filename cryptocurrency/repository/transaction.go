package repository

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Transaction struct {
	Time           string  `db:"time" json:"time"`
	BlockId        string  `db:"block_id" json:"blockID"`
	Hash           string  `db:"hash" json:"hash"`
	Size           int     `db:"size" json:"size"`
	Weight         int     `db:"weight" json:"weight"`
	IsCoinbase     bool    `db:"is_coinbase" json:"isCoinbase"`
	OutputTotal    int64   `db:"output_total" json:"outputTotal"`
	OutputTotalUSD float64 `db:"output_total_usd" json:"outputTotalUSD"`
	Fee            int64   `db:"fee" json:"fee"`
	FeeUSD         float64 `db:"fee_usd" json:"feeUSD"`
}

type TransactionOutboxEvent struct {
	Payload       Transaction `db:"payload"`
	AggregateType string      `db:"aggregatetype"`
	AggregateID   string      `db:"aggregateid"`
}

func (r *PGSQLRepository) InsertTransaction(t Transaction) error {
	trx, err := r.db.Beginx()
	if err != nil {
		r.logger.Error("Database Trx Init failed", zap.Error(err))
		trx.Rollback()
		return err
	}

	insertStmt := `
		INSERT INTO transactions 
		("time",block_id,hash,"size",weight,is_coinbase,output_total,output_total_usd,fee,fee_usd)
		VALUES (:time, :block_id, :hash, :size, :weight, :is_coinbase, :output_total, :output_total_usd, :fee, :fee_usd)
	`
	_, err = trx.NamedExec(insertStmt, t)
	if err != nil {
		r.logger.Error("Insert stmt error", zap.Error(err), zap.Any("transaction", t))
		trx.Rollback()
		return err
	}

	r.logger.Info("Successfully insert transaction", zap.Any("transaction", t))

	event := &TransactionOutboxEvent{
		Payload:       t,
		AggregateType: "INSERTED",
		AggregateID:   uuid.NewString(),
	}

	insertEventStmt := `
		INSERT INTO outbox_events (aggregatetype, aggregateid, payload)
		VALUES (:aggregatetype, :aggregateid, :payload)
	`
	_, err = trx.NamedExec(insertEventStmt, event)
	if err != nil {
		r.logger.Error("Insert Event stmt error", zap.Error(err), zap.Any("event", event))
		trx.Rollback()
		return err
	}

	r.logger.Info("Successfully insert event", zap.Any("transaction", t), zap.Any("event", event))

	err = trx.Commit()
	if err != nil {
		r.logger.Error("commit trx error", zap.Error(err), zap.Any("transaction", t), zap.Any("event", event))
		trx.Rollback()
		return err
	}

	return nil
}
