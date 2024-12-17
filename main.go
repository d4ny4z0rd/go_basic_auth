package main

import (
	"fmt"
	"log"
	"net/http"
	controllers "new_user_auth_prac/controller"
	"new_user_auth_prac/database"
	"new_user_auth_prac/routes"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err!=nil {
		log.Fatalf("Error loading the env file: %v \n", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router := mux.NewRouter()
	authRouter, userRouter := router.PathPrefix("/auth").Subrouter(), router.PathPrefix("/user").Subrouter()
	routes.AuthRouterHandler(authRouter)
	routes.UserRouterHandler(userRouter)

	db := database.InitDB()
	defer db.Close()

	database.CreateUsersTable(db)

	controllers.InitiateDBConnection(db)

	fmt.Printf("Server is running on port: %s \n", port)
	http.ListenAndServe(":"+port, router)
}

