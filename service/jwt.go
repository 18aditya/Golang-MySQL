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

func CreateJWT(key string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["id"] = key

	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return tokenStr, nil
}

func GetJWT(w http.ResponseWriter, r *http.Request) {
	var user model.Users
	// var arr_user []model.Users

	var (
		response model.Response
	)
	key := r.Header["Authorization"][0]
	db := config.Connect()
	defer db.Close()
	sql := fmt.Sprintf("Select * from Users Where IdUsers = %s", key)

	if r.Header["Authorization"] != nil {
		rows, err := db.Query(sql)
		if null := rows.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.CreatedAt); null == nil || err != nil {
			response.Status = 404
			response.Message = "user not found"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

		} else {
			for rows.Next() {
				if err := rows.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.CreatedAt); err != nil {
					response.Status = 404
					response.Message = fmt.Sprintf("%s", err)

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)

				} else {
					token, err := CreateJWT(key)
					if err != nil {
						response.Status = 403
						response.Message = "Cannot create token"

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
	} else {
		response.Status = 402
		response.Message = "Insert user Id"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

}
