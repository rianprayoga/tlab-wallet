package pg

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"tlab-wallet/internal/model"
)

var (
	ErrInsufucientBalance   = errors.New("insuficient balance")
	ErrSourceWalletNotFound = errors.New("source wallet not found")
	ErrTargetWalletNotFound = errors.New("target wallet not found")
)

func (pr *PgRepo) Transaction(sender string, receiver string, amount uint) (*model.TransactionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	tx, err := pr.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var enough bool
	if err = tx.
		QueryRowContext(ctx, `SELECT (balance >= $1) FROM wallets WHERE user_id = $2`, amount, sender).
		Scan(&enough); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSourceWalletNotFound
		}
		return nil, err
	}

	if !enough {
		return nil, ErrInsufucientBalance
	}

	var targetWalletId string
	if err = tx.
		QueryRowContext(ctx, `UPDATE wallets SET balance = balance + $1 where user_id = $2 RETURNING wallet_id`, amount, receiver).
		Scan(&targetWalletId); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTargetWalletNotFound
		}
		return nil, err
	}

	var w model.Wallet
	if err = tx.
		QueryRowContext(ctx, `UPDATE wallets SET balance = balance - $1 where user_id = $2 RETURNING wallet_id, balance`, amount, sender).
		Scan(&w.WalletId, &w.Balance); err != nil {
		return nil, err
	}

	var res model.TransactionResponse
	if err = tx.
		QueryRowContext(
			ctx,
			`insert into transactions(amount, sender_id, receiver_id, status) 
			values($1, $2, $3, 'SUCCESS') returning transaction_id, status, amount`,
			amount,
			sender,
			receiver).
		Scan(
			&res.TransactionId,
			&res.Status,
			&res.Amount,
		); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &res, nil

}

func (pr *PgRepo) GetTransactions(userId string, size int, page int) (*model.Page[model.TrxHistoryResponse], error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	var total int
	err := pr.DB.QueryRowContext(
		ctx,
		`SELECT count(transaction_id) FROM transactions WHERE sender_id = $1`,
		userId).Scan(&total)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &model.Page[model.TrxHistoryResponse]{
				Data:    make([]model.TrxHistoryResponse, 0),
				HasNext: false,
			}, nil
		}

		return nil, err
	}

	maxPage := math.Ceil(float64(total)/float64(size)) - 1
	if float64(page) > maxPage {

		return &model.Page[model.TrxHistoryResponse]{
			Data:    make([]model.TrxHistoryResponse, 0),
			HasNext: false,
		}, nil
	}

	rows, err := pr.DB.QueryContext(
		ctx,
		`SELECT transaction_id, receiver_id, amount, status, created_at 
		FROM transactions WHERE sender_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userId,
		size,
		size*page,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trxs []model.TrxHistoryResponse
	for rows.Next() {
		var trx model.TrxHistoryResponse
		if err := rows.Scan(
			&trx.TransactionId,
			&trx.Receiver,
			&trx.Amount,
			&trx.Status,
			&trx.CreatedAt,
		); err != nil {
			return nil, err
		}

		trxs = append(trxs, trx)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &model.Page[model.TrxHistoryResponse]{
		Data:    trxs,
		HasNext: page < int(maxPage),
	}, nil
}
