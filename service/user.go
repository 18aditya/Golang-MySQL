package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"main.go/config"
	"main.go/model"
)

type RealEmailResponse struct {
	Status string `json:"status"`
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {

	var user model.Users
	var arr_user []model.Users

	db := config.Connect()
	defer db.Close()

	rows, err := db.Query("Select * from Users")
	if err != nil {
		log.Print(err)
	} else {

		for rows.Next() {
			if err := rows.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.CreatedAt); err != nil {
				log.Fatal(err.Error())

			} else {
				var post model.Posts
				var arr_post []model.Posts
				rows, err := db.Query("Select Id,Title,Description FROM Posts WHERE Posts.UserId = ? ORDER BY ID", user.Id)
				if err != nil {
					log.Print(err)
				}

				for rows.Next() {
					if err := rows.Scan(&post.Id, &post.Title, &post.Description); err != nil {
						log.Fatal(err.Error())

					} else {
						arr_post = append(arr_post, post)
					}
				}

				user.Posts = arr_post
				arr_user = append(arr_user, user)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(arr_user)
	}
}

func InsertUser(w http.ResponseWriter, r *http.Request) {

	var (
		response model.Response
	)

	db := config.Connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)
	if err != nil {
		response.Status = 500
		response.Message = "Method Not Allowed"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {

		First_name := r.FormValue("first_name")
		Last_name := r.FormValue("last_name")
		Email := r.FormValue("email")

		pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

		url := "https://isitarealemail.com/api/email/validate?email=" + url.QueryEscape(Email)

		if matches := pattern.MatchString(Email); !matches {
			response.Status = 400
			response.Message = "Email format is not valid"

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Add("Bearer", `f8a521d2-7ffb-4cbf-aa75-5de19557a959`)

			if res, err := http.DefaultClient.Do(req); err != nil {
				response.Status = 404
				response.Message = "Email not found"

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				defer res.Body.Close()
			} else {
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)

				var myJson RealEmailResponse
				json.Unmarshal(body, &myJson)

				if myJson.Status != "valid" || err != nil {
					response.Status = 500
					response.Message = "Server cannot be reached"

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
				} else {

					if res, err := db.Exec("INSERT INTO Users (first_name, last_name,email) values (?,?,?)",
						First_name,
						Last_name,
						Email,
					); err != nil {
						me, ok := err.(*mysql.MySQLError)

						if !ok || me.Number == 1062 {
							response.Status = 1062
							response.Message = "Email is Already taken"
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(response)
						}

					} else {
						lastId, err := res.LastInsertId()
						if err != nil {
							response.Status = 500
							response.Message = "There's an Error"
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(response)
						}
						response.Status = 200
						response.Message = fmt.Sprint("Succesfully add: ", lastId)
						log.Print("Insert data to database")

						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(response)
					}
				}
			}
		}
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var response model.Response

	db := config.Connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}
	id := r.FormValue("id")
	First_name := r.FormValue("first_name")
	Last_name := r.FormValue("last_name")
	Email := r.FormValue("email")

	_, err = db.Exec("UPDATE Users set first_name = ?, last_name = ?,email = ? where IdUsers = ? ",
		First_name,
		Last_name,
		Email,
		id,
	)

	if err != nil {
		log.Print(err)
	}

	response.Status = 1
	response.Message = "Success Update Data"
	log.Print("Update data to database")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	var response model.Response
	var user model.Users
	var arr_user []model.Users
	vars := mux.Vars(r)
	id := vars["id"]

	db := config.Connect()
	defer db.Close()

	sql_query := fmt.Sprint("Select * from Users Where IdUsers = ", id)
	rows := db.QueryRow(sql_query)
	err := rows.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.CreatedAt)
	if err != nil && err == sql.ErrNoRows {
		response.Status = 404
		response.Message = "user not found"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	} else {
		var post model.Posts
		var arr_post []model.Posts
		rows, err := db.Query("Select Id,Title,Description FROM Posts WHERE Posts.UserId = ? ORDER BY ID", user.Id)
		if err != nil {
			log.Print(err)
		}

		for rows.Next() {
			if err := rows.Scan(&post.Id, &post.Title, &post.Description); err != nil {
				log.Fatal(err.Error())

			} else {
				arr_post = append(arr_post, post)
			}
		}

		user.Posts = arr_post
		arr_user = append(arr_user, user)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(arr_user)
	}

}

// func DeletePost(w http.ResponseWriter, r *http.Request) {
// 	var response model.Response

// 	db := config.Connect()
// 	defer db.Close()

// 	err := r.ParseMultipartForm(4096)
// 	if err != nil {
// 		panic(err)
// 	}

// 	id := r.FormValue("id")

// 	_, err = db.Exec("DELETE from posts where idposts = ? ",
// 		id,
// 	)

// 	if err != nil {
// 		log.Print(err)
// 	}

// 	response.Status = 1
// 	response.Message = "Success Delete Data"
// 	log.Print("Delete data to database")

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)

// }
