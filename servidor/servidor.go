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

type usuario struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

var data struct {
	ID int `json:"id"`
}

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	//ioutil.ReadAll(r.body) vai ler o que esta dentro do body, no caso um objeto json com  meus valores
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisiscao"))
		return
	}

	var usuario usuario
	if err := json.Unmarshal(body, &usuario); err != nil {
		w.Write([]byte("erro ao conectar"))
		return
	}

	db, err := db.Connection().DB()
	if err != nil {
		log.Panic("deu ruim no banco irmao")
	}

	defer db.Close()

	statement, err := db.Prepare("INSERT INTO usuarios (nome, email) VALUES ($1, $2) RETURNING id;")
	if err != nil {
		w.Write([]byte("erro ao conectar ao statement"))
		fmt.Println(statement)
		fmt.Println(err)
		return
	}
	defer statement.Close()

	insercao, err := statement.Exec(usuario.Nome, usuario.Email)

	fmt.Println(insercao)
	if err != nil {
		w.Write([]byte("erro ao executar ao statement"))
		return
	}

	// idInsercao, err := insercao.LastInsertId()
	var idInsercao int
	err = statement.QueryRow(usuario.Nome, usuario.Email).Scan(&idInsercao)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("errooooooooo filho da putaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("usuario com id %d", idInsercao)))
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {

	fmt.Println(&data)
	err := json.NewDecoder(r.Body).Decode(&data)
	fmt.Println(&data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := db.Connection().DB()
	if err != nil {
		log.Panic("deu ruim no banco irmao")
	}

	defer db.Close()
	var usuario usuario

	err = db.QueryRow("SELECT id, nome, email FROM usuarios WHERE id = $1", data.ID).Scan(&usuario.ID, &usuario.Nome, &usuario.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(usuario)
	responseString := string(response)
	fmt.Println(responseString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, err := db.Connection().DB()
	if err != nil {
		log.Panic("deu ruim no banco irmao")
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, nome, email FROM usuarios")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	usuarios := []usuario{}

	for rows.Next() {
		var usuario usuario
		err := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Email)
		if err != nil {
			log.Fatal(err)
		}
		usuarios = append(usuarios, usuario)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Converta a lista de usuários para JSON
	response, err := json.Marshal(usuarios)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func AtualizarUser(w http.ResponseWriter, r *http.Request) {
	var user struct {
		ID    int    `json:"id"`
		Nome  string `json:"nome"`
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)

	fmt.Println(reflect.TypeOf(user)) // int

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := db.Connection().DB()
	if err != nil {
		log.Panic("deu ruim no banco irmao")
	}

	defer db.Close()

	query := "UPDATE usuarios SET nome = $1, email = $2 WHERE id = $3"
	_, err = db.Exec(query, user.Nome, user.Email, user.ID) // Substitua "novoValor" e "id" pelos valores corretos

	if err != nil {
		// Trate o erro de atualização
		fmt.Println(err.Error())
	} else {
		// A atualização foi bem-sucedida
		fmt.Println("atualizou ")
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

func DeletarUser(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	id, err := strconv.ParseInt(param["id"], 10, 32)

	if err != nil {
		fmt.Println("deu erro na conversao")
	}

	db, err := db.Connection().DB()

	if err != nil {
		log.Panic("banco nao esta possivel de logar")
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
