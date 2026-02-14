package pg

import (
	"context"
	"database/sql"
	"errors"
	"tlab-wallet/internal/model"
)

var (
	ErrInsufucientBalance   = errors.New("insuficient balance")
	ErrSourceWalletNotFound = errors.New("source wallet not found")
	ErrTargetWalletNotFound = errors.New("target wallet not found")
)

func (pr *PgRepo) Transaction(sender string, receiver string, balance uint) (*model.Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	tx, err := pr.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var enough bool
	if err = tx.
		QueryRowContext(ctx, `SELECT (balance >= $1) FROM wallets WHERE user_id = $2`, balance, sender).
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
		QueryRowContext(ctx, `UPDATE wallets SET balance = balance + $1 where user_id = $2 RETURNING wallet_id`, balance, receiver).
		Scan(&targetWalletId); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTargetWalletNotFound
		}
		return nil, err
	}

	var res model.Wallet
	if err = tx.
		QueryRowContext(ctx, `UPDATE wallets SET balance = balance - $1 where user_id = $2 RETURNING wallet_id, balance`, balance, sender).
		Scan(&res.WalletId, &res.Balance); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &res, nil

}
