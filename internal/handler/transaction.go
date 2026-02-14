package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tlab-wallet/internal/model"
	"tlab-wallet/internal/repository/pg"
)

type query struct {
	Page int
	Size int
}

func (h *HttpHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userId := GetUidFromToken(r)

	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")

	q := query{
		Page: 0,
		Size: 5,
	}

	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			h.errorJSON(w, fmt.Errorf("invalid query param page"), http.StatusBadRequest)
			return
		}
		q.Page = pageInt
	}

	if size != "" {
		sizeInt, err := strconv.Atoi(size)
		if err != nil {
			h.errorJSON(w, fmt.Errorf("invalid query param size"), http.StatusBadRequest)
			return
		}
		q.Size = sizeInt
	}

	res, err := h.Db.GetTransactions(userId, q.Size, q.Page)
	if err != nil {
		h.Log.Debugf("error transaction: %+v", err)
		h.errorJSON(w, fmt.Errorf("unexpected error occured"))
		return
	}

	h.writeJson(w, http.StatusOK, res)

}

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

	res, err := h.Db.Transaction(userId, rq.Receiver, rq.Amount)
	if err != nil {

		if errors.Is(err, pg.ErrInsufucientBalance) || errors.Is(err, pg.ErrSourceWalletNotFound) || errors.Is(err, pg.ErrTargetWalletNotFound) {
			h.errorJSON(w, err, http.StatusBadRequest)
			return
		}

		h.Log.Debugf("error transaction: %+v", err)
		h.errorJSON(w, fmt.Errorf("unexpected error occured"))
		return
	}

	h.writeJson(w, http.StatusOK, res)

}
