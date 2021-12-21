package controllers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/brndedhero/finance/helpers"
	"github.com/brndedhero/finance/models"
	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"message": "Welcome to the Index Page"}`)
}

func AllAccountsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case "GET":
		data, err := models.GetAllAccounts()
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, data)
	default:
		message, err := helpers.PrepareString(405, nil)
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, message)
	}
}

func NewAccountHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		balance, _ := strconv.ParseFloat(req.FormValue("balance"), 32)
		data, err := models.CreateAccount(req.FormValue("name"), float32(balance))
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, data)
	default:
		message, err := helpers.PrepareString(405, nil)
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, message)
	}
}

func AccountHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(req)
	id, _ := strconv.ParseUint(params["id"], 10, 64)

	switch req.Method {
	case "GET":
		data, err := models.GetAccount(id)
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, data)
	case "POST":
		if err := req.ParseForm(); err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		balance, _ := strconv.ParseFloat(req.FormValue("balance"), 32)
		data, err := models.UpdateAccount(id, req.FormValue("name"), float32(balance))
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, data)
	case "DELETE":
		data, err := models.DeleteAccount(id)
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, data)
	default:
		message, err := helpers.PrepareString(405, nil)
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, message)
	}
}

func SearchAccountHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}
