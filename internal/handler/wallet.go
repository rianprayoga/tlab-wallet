package handler

import (
	"fmt"
	"net/http"
	"tlab-wallet/internal/model"

	"github.com/go-chi/jwtauth/v5"
)

func getUidFromToken(r *http.Request) string {
	_, claims, _ := jwtauth.FromContext(r.Context())
	userId := claims["uid"].(string)

	return userId
}

func (h *HttpHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	userId := getUidFromToken(r)

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
	userId := getUidFromToken(r)

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
