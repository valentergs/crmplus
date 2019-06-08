package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/user/CRMplus/controllers"
	"github.com/user/CRMplus/driver"
)

var db *sql.DB

func main() {

	db := driver.ConnectDB()
	controller := controllers.Controller{}

	// gorilla.mux
	router := mux.NewRouter()
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")
	router.HandleFunc("/usuario/add", controller.UsuarioAdd(db)).Methods("POST")
	//router.HandleFunc("/usuario", middlewares.TokenVerifyMiddleware(controller.UsuarioGetAll(db))).Methods("GET")
	router.HandleFunc("/usuario", controller.UsuarioGetAll(db)).Methods("GET")
	//router.HandleFunc("/usuario/{id}", middlewares.TokenVerifyMiddleware(controller.UsuarioGetOne(db))).Methods("GET")
	router.HandleFunc("/usuario/{id}", controller.UsuarioGetOne(db)).Methods("GET")
	router.HandleFunc("/search/usuario", controller.Search(db)).Methods("GET").Queries("q", "{q}")
	router.HandleFunc("/usuario/delete/{id}", controller.UsuarioDeleteOne(db)).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
		Debug:          true,
	})

	// Insert the middleware
	handler := c.Handler(router)

	log.Println("Listen on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
