package handler

import (
	"fmt"
	"net/http"
	"tlab-wallet/internal/model"

	"github.com/go-chi/jwtauth/v5"
)

func (h *HttpHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	userId := claims["uid"].(string)
	user, err := h.Db.GetUserById(userId)
	if err != nil {
		h.Log.Debugf("error get user: %+v", err)
		h.errorJSON(w, fmt.Errorf("unexpected error occured"))
	}

	h.writeJson(w, http.StatusOK, model.ProfileResponse{
		UserId:   user.UserId,
		Username: user.Username,
	})

}
