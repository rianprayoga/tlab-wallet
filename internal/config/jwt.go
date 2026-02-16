package config

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func GenerateAuthToken() *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET_KEY")), nil)
	return tokenAuth
}

func Verify(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := VerifyRequest(ja, r, jwtauth.TokenFromHeader)
			ctx = jwtauth.NewContext(ctx, token, err)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

func VerifyRequest(ja *jwtauth.JWTAuth, r *http.Request, findTokenFn func(r *http.Request) string) (jwt.Token, error) {
	var tokenString string

	tokenString = findTokenFn(r)

	if tokenString == "" {
		return nil, jwtauth.ErrNoTokenFound
	}

	return jwtauth.VerifyToken(ja, tokenString)
}

func writeOut(w http.ResponseWriter, message string) {
	res := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}

	out, _ := json.Marshal(res)

	h := w.Header()
	h.Del("Content-Length")
	h.Set("Content-Type", "application/json")
	h.Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(out)
}

func Authenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				writeOut(w, err.Error())
				return
			}

			if token == nil {
				writeOut(w, http.StatusText(http.StatusUnauthorized))
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
