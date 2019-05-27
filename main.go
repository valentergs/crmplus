package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/user/CRMplus/controllers"
	"github.com/user/CRMplus/driver"
	"github.com/user/CRMplus/middlewares"
)

var db *sql.DB

func main() {

	db := driver.ConnectDB()
	controller := controllers.Controller{}

	// gorilla.mux
	router := mux.NewRouter()
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")
	router.HandleFunc("/usuario/add", controller.UsuarioAdd(db)).Methods("POST")
	router.HandleFunc("/usuario", middlewares.TokenVerifyMiddleware(controller.UsuarioGetAll(db))).Methods("GET")
	router.HandleFunc("/usuario/{id}", middlewares.TokenVerifyMiddleware(controller.UsuarioGetOne(db))).Methods("GET")
	router.HandleFunc("/usuario/delete/{id}", middlewares.TokenVerifyMiddleware(controller.UsuarioDeleteOne(db))).Methods("DELETE")
	router.HandleFunc("/api/cliente", middlewares.TokenVerifyMiddleware(controller.ClienteAPI(db))).Methods("GET")

	log.Println("Listen on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
