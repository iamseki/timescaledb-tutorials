package core

import (
	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/repository"
)

type TransactionInput struct {
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

func (c *Core) transactionRepository(t TransactionInput) repository.Transaction {
	return repository.Transaction{
		Time:           t.Time,
		BlockId:        t.BlockId,
		Hash:           t.Hash,
		Size:           t.Size,
		Weight:         t.Weight,
		IsCoinbase:     t.IsCoinbase,
		OutputTotal:    t.OutputTotal,
		OutputTotalUSD: t.OutputTotalUSD,
		Fee:            t.Fee,
		FeeUSD:         t.FeeUSD,
	}
}

func (c *Core) InsertTransaction(t TransactionInput) error {
	transaction := c.transactionRepository(t)
	return c.repository.InsertTransaction(transaction)
}
