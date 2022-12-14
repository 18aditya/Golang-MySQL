package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
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
	vars := mux.Vars(r)
	key := vars["id"]
	var (
		response model.Response
	)
	// key := r.Header["Authorization"][0]
	db := config.Connect()
	defer db.Close()
	sql_query := fmt.Sprint("Select * from Users Where IdUsers = ", key)

	rows := db.QueryRow(sql_query)
	err := rows.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.CreatedAt)
	if err != nil && err == sql.ErrNoRows {
		response.Status = 404
		response.Message = "user not found"
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
