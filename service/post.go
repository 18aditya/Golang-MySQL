package service

import (
	"encoding/json"
	"log"
	"net/http"

	"main.go/config"
	"main.go/model"
)

func GetAllPost(w http.ResponseWriter, r *http.Request) {
	var post model.Posts
	var arr_post []model.Posts

	db := config.Connect()
	defer db.Close()

	rows, err := db.Query("Select Id,Title,Description from Posts")
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(arr_post)

}

func InsertsPost(w http.ResponseWriter, r *http.Request) {

	var response model.Response

	db := config.Connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	UserId := r.FormValue("userid")

	_, err = db.Exec("INSERT INTO Posts (Title, Description,UserId) values (?,?,?)",
		title,
		description,
		UserId,
	)

	if err != nil {
		log.Print(err)
	}

	response.Status = 1
	response.Message = "Success Add"
	log.Print("Insert data to database")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var response model.Response

	db := config.Connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	id := r.FormValue("id")
	title := r.FormValue("title")
	description := r.FormValue("description")

	_, err = db.Exec("UPDATE Posts set Title = ?, Description = ? where IdPosts = ? ",
		title,
		description,
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

func DeletePost(w http.ResponseWriter, r *http.Request) {
	var response model.Response

	db := config.Connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	id := r.FormValue("id")

	_, err = db.Exec("DELETE from Posts where IdPosts = ? ",
		id,
	)

	if err != nil {
		log.Print(err)
	}

	response.Status = 1
	response.Message = "Success Delete Data"
	log.Print("Delete data to database")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
