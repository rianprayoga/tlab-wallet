package model

type Wallet struct {
	WalletId  string
	UserId    string
	Balance   uint
	CreatedAt string
}

type WalletResponse struct {
	WalletId string `json:"walletId"`
	Balance  uint   `json:"balance" `
}
