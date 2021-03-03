package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/httpWebServer/m/handler"
	pg "github.com/httpWebServer/m/postgres"
	logger "github.com/ipfs/go-log/v2"
)

type PgManager struct {
	db *sql.DB
}

var log = logger.Logger("http/server")

func (PMgr *PgManager) signUp(w http.ResponseWriter, r *http.Request) {
	log.Infof("signup requested")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to get body read", http.StatusInternalServerError)
		log.Errorf("Failed to get body %s", err.Error())
		return
	}
	r.Body.Close()
	var bodyElement map[string]string
	err = json.Unmarshal(body, &bodyElement)
	if err != nil {
		http.Error(w, "Failed unmarshal request body", http.StatusInternalServerError)
		log.Errorf("Failed to unmarshal request body %s", err.Error())
		return
	}
	email, ok := bodyElement["email"]
	if !ok {
		http.Error(w, "Email not found", http.StatusBadRequest)
		log.Errorf("Email not found in body")
		return
	}
	name, ok := bodyElement["name"]
	if !ok {
		http.Error(w, "Name not found", http.StatusBadRequest)
		log.Errorf("Name not found in body")
		return
	}
	password, ok := bodyElement["password"]
	if !ok {
		http.Error(w, "password not found", http.StatusBadRequest)
		log.Errorf("Password not found in body")
		return
	}
	err, id := handler.SignUp(email, name, password, PMgr.db)
	if err != nil {
		http.Error(w, "Failed to signup user", http.StatusInternalServerError)
		log.Errorf("Failed to signup user %s")
		return
	}
	byteData, err := json.Marshal(id)
	if err != nil {
		http.Error(w, "Failed marshal data", http.StatusInternalServerError)
		log.Errorf("Failed to marshal data %s", err.Error())
		return
	}

	w.Header().Set("content-type", "application/json")
	_, err = w.Write(byteData)
	if err != nil {
		http.Error(w, "Failed write response", http.StatusNotAcceptable)
		log.Errorf("Failed to write response to api %s", err.Error())
		return
	}
}

func (PMgr *PgManager) login(w http.ResponseWriter, r *http.Request) {
	log.Infof("login requested")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusFound)
		log.Errorf("Failed to read body %s", err.Error())
		return
	}
	var bodyElement map[string]string
	err = json.Unmarshal(body, &bodyElement)
	if err != nil {
		http.Error(w, "Failed to unmarhsal body", http.StatusBadRequest)
		log.Errorf("Failed to unmarshal body %s", err.Error())
		return
	}
	userid, ok := bodyElement["userid"]
	if !ok {
		http.Error(w, "Userid not found", http.StatusNoContent)
		log.Errorf("User-id not found in body")
		return
	}
	password, ok := bodyElement["password"]
	if !ok {
		http.Error(w, "Password not found", http.StatusNoContent)
		log.Errorf("Password not found in body")
		return
	}
	token, err := handler.Login(userid, password, PMgr.db)
	if err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		log.Errorf("Failed to get JWT token %s", err.Error())
		return
	}
	byteData, err := json.Marshal(token)
	if err != nil {
		if err != nil {
			http.Error(w, "Failed marshal data", http.StatusInternalServerError)
			log.Errorf("Failed to marshal data %s", err.Error())
			return
		}
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(byteData)
	if err != nil {
		http.Error(w, "Failed write response", http.StatusNotAcceptable)
		log.Errorf("Failed to write response to api %s", err.Error())
		return
	}
}

func (pMgr *PgManager) validation(w http.ResponseWriter, r *http.Request) {
	log.Infof("token validation called")
	token, ok := r.URL.Query()["token"]
	if !ok {
		http.Error(w, "token is not found", http.StatusNoContent)
		log.Errorf("toke is not found")
		return
	}
	message, err := handler.TokenValidation(token[0], pMgr.db)
	if err != nil {
		http.Error(w, "Token is not valid", http.StatusNotAcceptable)
		log.Errorf("Token failed error %s", err.Error())
		return
	}
	byteData, err := json.Marshal(message)
	if err != nil {
		if err != nil {
			http.Error(w, "Failed marshal data", http.StatusInternalServerError)
			log.Errorf("Failed to marshal data %s", err.Error())
			return
		}
	}
	w.Header().Set("content-type", "applicatioin/json")
	_, err = w.Write(byteData)
	if err != nil {
		http.Error(w, "Failed write response", http.StatusNotAcceptable)
		log.Errorf("Failed to write response to api %s", err.Error())
		return
	}
}

func main() {
	logger.SetLogLevel("*", "Debug")
	pMgr := &PgManager{
		db: pg.PostgresConfig(),
	}
	router := mux.NewRouter()
	router.HandleFunc("/signup", pMgr.signUp).Methods("POST")
	router.HandleFunc("/login", pMgr.login).Methods("POST")
	router.HandleFunc("/valid", pMgr.validation).Methods("GET")
	http.ListenAndServe(":8008", router)
	log.Infof("Serve is running on loclahost:8008")
}
