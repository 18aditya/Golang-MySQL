package middleware

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
				}
				return []byte(os.Getenv("SECRET")), nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized: " + err.Error()))
			}

			if token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(r.Header["Token"][1]))
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
		}
	})
}
