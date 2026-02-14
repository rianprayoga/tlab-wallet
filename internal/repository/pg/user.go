package pg

import (
	"context"
	"database/sql"
	"time"
	"tlab-wallet/internal/model"
)

type PgRepo struct {
	DB *sql.DB
}

const DbTimeout = 3 * time.Second

func (pr *PgRepo) UsernameExist(username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	var id string
	err := pr.DB.QueryRowContext(
		ctx,
		`SELECT user_id FROM users WHERE username = $1`,
		username,
	).Scan(
		&id,
	)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return true, err
	}

	return true, nil

}

func (pr *PgRepo) AddUser(rq model.RegisterUserRequest) (*model.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	var res model.User
	err := pr.DB.QueryRowContext(
		ctx,
		`insert into users(username, password) values($1,crypt($2,gen_salt('bf',10))) returning user_id, username`,
		rq.Username,
		rq.Password,
	).Scan(
		&res.UserId,
		&res.Username,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (pr *PgRepo) GetUser(username string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	var res model.User

	err := pr.DB.QueryRowContext(
		ctx,
		`SELECT user_id, username, password FROM users WHERE username = $1`,
		username,
	).Scan(
		&res.UserId,
		&res.Username,
		&res.Password,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (pr *PgRepo) GetUserById(userId string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	var res model.User

	err := pr.DB.QueryRowContext(
		ctx,
		`SELECT user_id, username FROM users WHERE user_id = $1`,
		userId,
	).Scan(
		&res.UserId,
		&res.Username,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
