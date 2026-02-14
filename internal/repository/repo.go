package repository

import "tlab-wallet/internal/model"

type Repo interface {
	UsernameExist(username string) (bool, error)
	AddUser(rq model.RegisterUserRequest) (*model.User, error)
	GetUser(username string) (*model.User, error)
	GetUserById(userId string) (*model.User, error)

	CreateWallet(userId string) error
	GetWallet(userId string) (*model.Wallet, error)
	TopUpWallet(userId string, balance uint) (*model.Wallet, error)

	Transaction(sender string, receiver string, amount uint) (*model.TransactionResponse, error)
	GetTransactions(userId string, size int, page int) (*model.Page[model.TrxHistoryResponse], error)
}
