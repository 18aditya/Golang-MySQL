package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"main.go/config"
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

	var (
		response model.Response
	)

	db := config.Connect()
	defer db.Close()

	if r.Header["Authorization"] != nil {
		rows, err := db.Query("Select Id from Users Where IdUsers = 1 ")
		if err != nil || rows == nil {
			response.Status = 404
			response.Message = "user not found"

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
