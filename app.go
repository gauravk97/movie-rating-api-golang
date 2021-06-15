// app.go

package main

import (
	"database/sql"

	"fmt"
	"log"

	"encoding/json"
	"net/http"
	"strconv"

	// tom: go get required
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// tom: initial function is empty, it's filled afterwards
// func (a *App) Initialize(user, password, dbname string) { }

// tom: added "sslmode=disable" to connection string
func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	// tom: this line is added after initializeRoutes is created later on
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func (a *App) searchMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	movie, err := searchMovie(a.DB, name)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	}

	respondWithJSON(w, http.StatusOK, movie)
}

func (a *App) getUserActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	userActivity, err := getUserActivity(a.DB, id)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User activity not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	}

	respondWithJSON(w, http.StatusOK, userActivity)
}

func (a *App) addUserRatingOnMovie(w http.ResponseWriter, r *http.Request) {
	var p user_movie_rating
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.addUserRatingOnMovie(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) addUserCommentOnMovie(w http.ResponseWriter, r *http.Request) {
	var p user_movie_comments
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.addUserCommentOnMovie(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/movie/search/{name}", a.searchMovie).Methods("GET")
	a.Router.HandleFunc("/movie/get_user_activity/{id:[0-9]+}", a.getUserActivity).Methods("GET")
	a.Router.HandleFunc("/movie/add_user_rating/", a.addUserRatingOnMovie).Methods(("POST"))
	a.Router.HandleFunc("/movie/add_user_comment/", a.addUserCommentOnMovie).Methods("POST")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
