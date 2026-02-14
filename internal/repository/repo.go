package repository

import "tlab-wallet/internal/model"

type Repo interface {
	UsernameExist(username string) (bool, error)
	AddUser(rq model.RegisterUserRequest) (*model.User, error)
	GetUser(username string) (*model.User, error)
	GetUserById(userId string) (*model.User, error)

	CreateWallet(userId string) error
}
