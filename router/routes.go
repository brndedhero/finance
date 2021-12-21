package router

import (
	"github.com/brndedhero/finance/controllers"
	"github.com/brndedhero/finance/middleware"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.PrometheusMiddleware)
	r.Path("/metrics").Handler(promhttp.Handler())

	r.HandleFunc("/", controllers.HomeHandler)
	r.HandleFunc("/accounts", controllers.AllAccountsHandler)
	r.HandleFunc("/accounts/new", controllers.NewAccountHandler)
	r.HandleFunc("/accounts/{id:[\\d]+}", controllers.AccountHandler)
	r.HandleFunc("/accounts/search", controllers.SearchAccountHandler)

	return r
}
