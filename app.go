package main

import (
	"github.com/julienschmidt/httprouter"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"fmt"
	"strconv"
	"encoding/json"

	"github.com/juheyne/ShoppingLedger/model"
)

type App struct {
	Router *httprouter.Router
	DB *sql.DB
}

func (a *App) Initialize(path string) {
	a.DB = model.Open(path)

	a.Router = httprouter.New()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.GET("/", a.index)
	a.Router.GET("/expenses", a.getExpenses)
	a.Router.POST("/expense", a.createExpense)
	a.Router.GET("/expense/:id", a.getExpense)
	a.Router.PUT("/expense/:id", a.updateExpense)
	a.Router.DELETE("/expense/:id", a.deleteExpense)
}

func (a *App) index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func (a *App) getExpenses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 100 || count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	expenses, err := model.GetExpenses(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, expenses)
}

func (a *App) getExpense(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	e := model.Expense{ID: id}
	if err := e.Get(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Expense not found")
		default:
			fmt.Println(err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, e)
}

func (a *App) createExpense(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var e model.Expense
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		fmt.Println("CreateExpense Decode Error: ", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := e.Create(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, e)
}

func (a *App) updateExpense(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	var e model.Expense
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	e.ID = id

	if err := e.Update(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, e)
}

func (a *App) deleteExpense(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	e := model.Expense{ID: id}
	if err := e.Delete(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
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


