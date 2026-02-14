package pg

import "context"

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
