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
	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST")
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")
	router.HandleFunc("/admin", middlewares.TokenVerifyMiddleware(controller.Admin(db))).Methods("GET")

	log.Println("Listen on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
