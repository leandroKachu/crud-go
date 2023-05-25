package main

import (
	"crud-go/servidor"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// db.AutoMigrate(&campaign.Campaign{}, &campaign.Contact{})

	r := mux.NewRouter()
	r.HandleFunc("/criar", servidor.CriarUsuario).Methods(http.MethodPost)
	r.HandleFunc("/buscarUsuario", servidor.BuscarUsuario).Methods(http.MethodGet)
	r.HandleFunc("/buscarUsuarios", servidor.BuscarUsuarios).Methods(http.MethodGet)
	r.HandleFunc("/atualizarUser", servidor.AtualizarUser).Methods(http.MethodPost)
	r.HandleFunc("/deletarUser/{id}", servidor.DeletarUser).Methods(http.MethodDelete)

	fmt.Println("Escutando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
