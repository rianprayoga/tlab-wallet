package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"io"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type JSONResponse struct {
	Message string `json:"message"`
}

func GetUidFromToken(r *http.Request) string {
	_, claims, _ := jwtauth.FromContext(r.Context())
	userId := claims["uid"].(string)

	return userId
}

func (h *HttpHandler) writeJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {

	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (h *HttpHandler) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)
	if err != nil {
		h.errorJSON(w, fmt.Errorf("unexpected error occurred"))
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		h.errorJSON(w, fmt.Errorf("invalid body"), http.StatusBadRequest)
		return err
	}

	err = validate.Struct(data)
	var validationErrors validator.ValidationErrors
	if err == nil {
		return nil
	}

	if !errors.As(err, &validationErrors) {
		h.errorJSON(w, fmt.Errorf("unexpected error occurred"))
		return err
	}

	ve := validationErrors[0] // get the 1st error

	h.errorJSON(w, fmt.Errorf("%s", strings.ToLower(ve.Field())+" does not match the requirement"), http.StatusBadRequest)
	return err

}

func (app *HttpHandler) errorJSON(w http.ResponseWriter, err error, status ...int) {

	statusCode := http.StatusInternalServerError
	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := JSONResponse{
		Message: err.Error(),
	}

	app.writeJson(w, statusCode, payload)
}
