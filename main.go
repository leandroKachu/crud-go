package main

import (
	"crud-go/servidor"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/create", servidor.CreateUSer).Methods(http.MethodPost)
	r.HandleFunc("/findUser", servidor.FindUser).Methods(http.MethodGet)
	r.HandleFunc("/findAllUsers", servidor.FindAllUsers).Methods(http.MethodGet)
	r.HandleFunc("/updateUser", servidor.UpdateUser).Methods(http.MethodPost)
	r.HandleFunc("/deletarUser/{id}", servidor.DeletarUser).Methods(http.MethodDelete)

	fmt.Println("Escutando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
