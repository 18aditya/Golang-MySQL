package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"main.go/model"
)

func CreateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour).Unix()

	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return tokenStr, nil
}

func GetJWT(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	var (
		response model.Response
	)
	if err != nil {
		response.Status = 502
		response.Message = "Something is wrong when generating your token"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

	if r.Header["Authorization"] != nil {
		if r.Header["Authorization"][0] != os.Getenv("API_KEY") {

			response.Status = 401
			response.Message = "Api Key is wrong"

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			token, err := CreateJWT()
			if err != nil {
				response.Status = 403
				response.Message = fmt.Sprintf("Cannot create token %s", err)

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			} else {
				response.Status = 200
				response.Message = fmt.Sprintf("%s", token)

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
		}
	}

}
