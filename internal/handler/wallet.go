package handler

import (
	"fmt"
	"net/http"
	"tlab-wallet/internal/model"
)

func (h *HttpHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	userId := GetUidFromToken(r)

	res, err := h.Db.GetWallet(userId)
	if err != nil {
		h.Log.Debugf("error get wallet: %+v", err)
		h.errorJSON(w, fmt.Errorf("unexpected error occured"))
	}

	h.writeJson(w, http.StatusOK, model.WalletResponse{
		WalletId: res.WalletId,
		Balance:  res.Balance,
	})

}

func (h *HttpHandler) TopUpWallet(w http.ResponseWriter, r *http.Request) {
	userId := GetUidFromToken(r)

	var rq model.TopUpRequest
	err := h.readJSON(w, r, &rq)
	if err != nil {
		return
	}

	wallet, err := h.Db.TopUpWallet(userId, rq.Balance)
	if err != nil {
		h.Log.Debugf("error topup wallet: %+v", err)
		h.errorJSON(w, fmt.Errorf("unexpected error occured"))
		return
	}

	h.writeJson(w, http.StatusOK, model.WalletResponse{
		WalletId: wallet.WalletId,
		Balance:  wallet.Balance,
	})

}
