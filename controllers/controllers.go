package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/user/CRMplus/models"
	"github.com/user/CRMplus/utils"
	"golang.org/x/crypto/bcrypt"
)

//Controller will be exported
type Controller struct{}

//Signup will be exported
func (c Controller) Signup(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var usuario models.Usuario

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

		utils.ResponseJSON(w, usuario)

	}
}

// Login will be exported
func (c Controller) Login(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var usuario models.Usuario
		var jwt models.JWT
		var error models.Error

		json.NewDecoder(r.Body).Decode(&usuario)

		// Essa é a senha passada pelo usuário quando enviar o request
		senha := usuario.Senha

		// Verificar se usuário existe no DB
		row := db.QueryRow("select * from usuario where email=$1;", usuario.Email)
		err := row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Senha, &usuario.Email, &usuario.Celular, &usuario.Superuser, &usuario.Ativo, &usuario.Departamento)
		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "Usuário inexistente"
				utils.RespondWithError(w, http.StatusBadRequest, error)
				return
			} else {
				log.Fatal(err)
			}
		}

		token, err := utils.GenerateToken(usuario)
		if err != nil {
			log.Fatal(err)
		}

		// usuario.Senha é o resultado do DB query
		hashedPassword := usuario.Senha

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(senha))

		if err != nil {
			error.Message = "Senha inválida"
			utils.RespondWithError(w, http.StatusUnauthorized, error)
			return
		}
		w.WriteHeader(http.StatusOK)

		jwt.Token = token
		utils.ResponseJSON(w, jwt)

	}
}
