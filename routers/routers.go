package routers

import (
	"licenseserver/controllers"

	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	r := mux.NewRouter()
	r.HandleFunc("/v1/regist", controllers.Regist).Methods("POST")
	r.HandleFunc("/v1/updateUserInfo", controllers.UpdateUserInfo).Methods("POST")
	r.HandleFunc("/v1/checkPhoneExist", controllers.CheckPhoneNumberExist).Methods("GET")
	r.HandleFunc("/v1/getLicense", controllers.GetLicense).Methods("GET")

	r.HandleFunc("/v2/newSerial", controllers.NewSerial).Methods("GET")
	r.HandleFunc("/v2/registerSerial", controllers.RegisterSerial).Methods("POST")
	r.HandleFunc("/v2/sync", controllers.Sync).Methods("POST")
	r.HandleFunc("/v2/extend", controllers.Extend).Methods("POST")
	r.HandleFunc("/v2/getAvailableSerial", controllers.GetAvailableSerial).Methods("GET")
	r.HandleFunc("/v2/getSerialByPhone", controllers.GetSerialByPhone).Methods("GET")

	r.HandleFunc("/", controllers.Home).Methods("GET")
	r.HandleFunc("/login", controllers.Login)
	r.HandleFunc("/logout", controllers.Logout).Methods("GET")

	http.Handle("/", r)
}
