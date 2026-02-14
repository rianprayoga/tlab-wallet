package handler

import (
	"errors"
	"fmt"
	"net/http"
	"tlab-wallet/internal/model"
	"tlab-wallet/internal/repository/pg"
)

func (h *HttpHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userId := GetUidFromToken(r)

	var rq model.TransferRequest
	err := h.readJSON(w, r, &rq)
	if err != nil {
		return
	}

	if rq.Receiver == userId {
		h.errorJSON(w, fmt.Errorf("receiver can't be the same with sender"), http.StatusBadRequest)
		return
	}

	wallet, err := h.Db.Transaction(userId, rq.Receiver, rq.Balance)
	if err != nil {

		if errors.Is(err, pg.ErrInsufucientBalance) || errors.Is(err, pg.ErrSourceWalletNotFound) || errors.Is(err, pg.ErrTargetWalletNotFound) {
			h.errorJSON(w, err, http.StatusBadRequest)
			return
		}

		h.Log.Debugf("error transaction: %+v", err)
		h.errorJSON(w, fmt.Errorf("unexpected error occured"))
		return
	}

	h.writeJson(w, http.StatusOK, model.WalletResponse{
		WalletId: wallet.WalletId,
		Balance:  wallet.Balance,
	})

}
