package routes

import (
	controllers "new_user_auth_prac/controller"

	"github.com/gorilla/mux"
)

func AuthRouterHandler(r *mux.Router) {
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/signup", controllers.Signup).Methods("POST")
}