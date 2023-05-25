package servidor

import (
	"crud-go/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

var data struct {
	ID int `json:"id"`
}

func CreateUSer(w http.ResponseWriter, r *http.Request) {
	//ioutil.ReadAll(r.body) vai ler o que esta dentro do body, no caso um objeto json com  meus valores
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisiscao"))
		return
	}

	var user user
	if err := json.Unmarshal(body, &user); err != nil {
		w.Write([]byte("erro ao conectar"))
		return
	}

	db, err := db.Connection().DB()
	if err != nil {
		log.Panic("not connected")
		fmt.Println(err.Error())
	}

	defer db.Close()

	query := "INSERT INTO usuarios (nome, email) VALUES ($1, $2) RETURNING id"
	var userID int

	err = db.QueryRow(query, user.Nome, user.Email).Scan(&userID)
	// idInsercao, err := insercao.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte("Error on creating user"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("user with id %d", userID)))
}

func FindUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println(&data)
	err := json.NewDecoder(r.Body).Decode(&data)
	fmt.Println(&data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := db.Connection().DB()
	if err != nil {
		log.Panic("not connected")
		fmt.Println(err.Error())
	}

	defer db.Close()
	var user user

	err = db.QueryRow("SELECT id, nome, email FROM usuarios WHERE id = $1", data.ID).Scan(&user.ID, &user.Nome, &user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(user)
	responseString := string(response)
	fmt.Println(responseString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func FindAllUsers(w http.ResponseWriter, r *http.Request) {
	db, err := db.Connection().DB()
	if err != nil {
		log.Panic("not connected")
		fmt.Println(err.Error())
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, nome, email FROM usuarios")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	users := []user{}

	for rows.Next() {
		var user user
		err := rows.Scan(&user.ID, &user.Nome, &user.Email)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Converta a lista de usuários para JSON
	response, err := json.Marshal(users)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var UpdateUser struct {
		ID    int    `json:"id"`
		Nome  string `json:"nome"`
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&UpdateUser)

	fmt.Println(reflect.TypeOf(UpdateUser)) // int

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := db.Connection().DB()
	if err != nil {
		log.Panic("not connected")
		fmt.Println(err.Error())
	}

	defer db.Close()

	query := "UPDATE usuarios SET nome = $1, email = $2 WHERE id = $3"
	_, err = db.Exec(query, UpdateUser.Nome, UpdateUser.Email, UpdateUser.ID) // Substitua "novoValor" e "id" pelos valores corretos

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Updated user")
	}

	response, err := json.Marshal(UpdateUser)
	responseString := string(response)
	fmt.Println(responseString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func DeletarUser(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	id, err := strconv.ParseInt(param["id"], 10, 32)

	if err != nil {
		fmt.Println("not possible to convert")
	}

	db, err := db.Connection().DB()

	if err != nil {
		log.Panic("not connected")
		fmt.Println(err.Error())

	}

	defer db.Close()

	query := "DELETE FROM usuarios WHERE id = $1"

	_, err = db.Exec(query, id)

	if err != nil {
		fmt.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("finalized fi"))

}
