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

type TopUpRequest struct {
	Balance uint `json:"balance" validate:"required,gte=0"`
}

type TransferRequest struct {
	Receiver string `json:"receiver" validate:"required,uuid4"`
	Amount   uint   `json:"amount" validate:"required,gte=0"`
}
