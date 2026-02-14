package pg

import (
	"context"
	"tlab-wallet/internal/model"
)

func (pr *PgRepo) CreateWallet(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	_, err := pr.DB.ExecContext(
		ctx,
		`insert into wallets(user_id, balance) values($1, $2)`,
		userId,
		0,
	)

	if err != nil {
		return err
	}

	return nil
}

func (pr *PgRepo) TopUpWallet(userId string, balance uint) (*model.Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	tx, err := pr.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var res model.Wallet
	if err = tx.
		QueryRowContext(ctx, `UPDATE wallets SET balance = balance + $2 where user_id = $1 RETURNING wallet_id, balance`, userId, balance).
		Scan(&res.WalletId, &res.Balance); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &res, nil
}

func (pr *PgRepo) GetWallet(userId string) (*model.Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	var res model.Wallet
	err := pr.DB.QueryRowContext(
		ctx,
		`SELECT wallet_id, user_id, balance, created_at FROM wallets WHERE user_id = $1`,
		userId,
	).Scan(
		&res.WalletId,
		&res.UserId,
		&res.Balance,
		&res.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
