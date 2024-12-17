package routes

import (
	controllers "new_user_auth_prac/controller"

	"github.com/gorilla/mux"
)

func UserRouterHandler(r *mux.Router) {
	r.HandleFunc("/getUser/{userId}", controllers.GetUser).Methods("GET")
	r.HandleFunc("/getUsers", controllers.GetUsers).Methods("GET")
	r.HandleFunc("/updateUser/{userId}", controllers.UpdateUser).Methods("POST")
	r.HandleFunc("/deleteUser/{userId}", controllers.DeleteUser).Methods("DELETE")
}