package handler

import (
	"fmt"
	"net/http"
	"tlab-wallet/internal/model"

	"github.com/go-chi/jwtauth/v5"
)

func (h *HttpHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	userId := claims["uid"].(string)
	res, err := h.Db.GetWallet(userId)
	if err != nil {
		if err != nil {
			h.Log.Debugf("error get wallet: %+v", err)
			h.errorJSON(w, fmt.Errorf("unexpected error occured"))
		}
	}

	h.writeJson(w, http.StatusOK, model.WalletResponse{
		WalletId: res.WalletId,
		Balance:  res.Balance,
	})

}
