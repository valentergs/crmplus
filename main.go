package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//Usuario is an exportable type
type Usuario struct {
	Nome         string `json:"nome"`
	Sobrenome    string `json:"sobrenome"`
	Senha        string `json:"senha"`
	Email        string `json:"email"`
	Celular      string `json:"celular"`
	Superuser    bool   `json:"superuser"`
	Ativo        bool   `json:"ativo"`
	Departamento string `json:"departamento"`
}

//Error is an exportable type
type Error struct {
	Message string `json:"message"`
}

var db *sql.DB

func main() {

	var err error

	db, err = sql.Open("postgres", "user=rodrigovalente password=Gustavo2012 host=localhost port=5432 dbname=CRMplus sslmode=disable")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/signup", signup).Methods("POST")

	log.Println("Listen on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func respondWithError(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
	return
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func signup(w http.ResponseWriter, r *http.Request) {

	var usuario Usuario

	json.NewDecoder(r.Body).Decode(&usuario)

	// Hash usuario.Senha
	hash, err := bcrypt.GenerateFromPassword([]byte(usuario.Senha), 10)
	if err != nil {
		log.Fatal(err)
	}

	// Assign hash to usuario.Senha
	usuario.Senha = string(hash)

	sqlstmt := `insert into usuario (nome, sobrenome, senha, email, celular, superuser, ativo, departamento) values ($1,$2,$3,$4,$5,$6,$7,$8);`
	_, err = db.Exec(sqlstmt, usuario.Nome, usuario.Sobrenome, usuario.Senha, usuario.Email, usuario.Celular, usuario.Superuser, usuario.Ativo, usuario.Departamento)
	if err != nil {
		panic(err)
	}

	// Hide output of usuario.Senha
	usuario.Senha = "********"

	w.Header().Set("Content-Type", "application/json")

	responseJSON(w, usuario)
}
