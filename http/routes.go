package http


import (

    "github.com/gorilla/mux"
)


func RegisterRoutes (r *mux.Router){


	r.HandleFunc("/start/{service}", startHandler).Methods("GET")

} 