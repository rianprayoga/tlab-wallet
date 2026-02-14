package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
	"tlab-wallet/internal/model"

	"golang.org/x/crypto/bcrypt"
)

func (h *HttpHandler) validateUsername(username string, w http.ResponseWriter) bool {
	match, err := regexp.MatchString("^[a-zA-Z0-9_.]{5,100}$", username)
	if err != nil {
		h.errorJSON(w, fmt.Errorf("unexpected error occured"))
		return false
	}

	if !match {
		h.errorJSON(w, fmt.Errorf("username not meet requierment"), http.StatusBadRequest)
		return false
	}
	return true
}

func (h *HttpHandler) LoginUser(w http.ResponseWriter, r *http.Request) {

	var rq model.LoginUserRequest
	err := h.readJSON(w, r, &rq)
	if err != nil {
		return
	}

	if h.validateUsername(rq.Username, w) == false {
		return
	}

	user, err := h.Db.GetUser(strings.ToLower(rq.Username))
	if err != nil {
		if err == sql.ErrNoRows {
			h.errorJSON(w, fmt.Errorf("username or password is wrong"), http.StatusUnauthorized)
			return
		}

		h.Log.Debugf("error get user: %+v", err)
		h.errorJSON(w, fmt.Errorf("unexpected error occured"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rq.Password))
	if err != nil {
		h.errorJSON(w, fmt.Errorf("username or password is wrong"), http.StatusUnauthorized)
		return
	}

	_, token, err := h.Jwt.Encode(map[string]interface{}{
		"uid": user.UserId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(5 * time.Hour).Unix(),
	})
	if err != nil {
		h.Log.Debugf("error generating token: %+v", err)
		h.errorJSON(w, fmt.Errorf("unexpected error occured"))
		return
	}

	h.writeJson(w, http.StatusCreated, struct {
		AuthToken string `json:"authToken"`
	}{
		AuthToken: token,
	})

}

func (h *HttpHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var rq model.RegisterUserRequest

	err := h.readJSON(w, r, &rq)
	if err != nil {
		return
	}

	if h.validateUsername(rq.Username, w) == false {
		return
	}

	rq.Username = strings.ToLower(rq.Username)
	isExist, err := h.Db.UsernameExist(rq.Username)
	if err != nil {
		h.Log.Debugf("error checking username: %+v", err)
		h.errorJSON(w, fmt.Errorf("unexpected error occured"))
		return
	}

	if isExist {
		h.errorJSON(w, fmt.Errorf("username alread exist"), http.StatusConflict)
		return
	}

	user, err := h.Db.AddUser(rq)
	if err != nil {
		h.Log.Debugf("error adding user: %+v", err)
		h.errorJSON(w, fmt.Errorf("unexpted error occured"))
		return
	}

	err = h.Db.CreateWallet(user.UserId)
	if err != nil {
		h.Log.Debugf("error creating wallet: %+v", err)
		h.errorJSON(w, fmt.Errorf("unexpted error occured"))
		return
	}

	h.writeJson(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("User %s created", rq.Username),
	})

}
