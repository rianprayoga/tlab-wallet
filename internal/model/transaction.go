package model

type TransactionResponse struct {
	TransactionId string `json:"transactionId"`
	Amount        uint   `json:"amount"`
	Status        string `json:"status"`
}

type Page[T any] struct {
	Data    []T  `json:"data"`
	HasNext bool `json:"hasNext"`
}

type TrxHistoryResponse struct {
	TransactionId string `json:"transactionId"`
	Receiver      string `json:"receiverId"`
	Amount        uint   `json:"amount"`
	Status        string `json:"status"`
	CreatedAt     string `json:"createdAt"`
}
