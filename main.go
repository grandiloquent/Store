package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"store/common"
	"store/common/datastore"
	"store/middleware"

	"store/handlers"
)

const (
	WEBSERVERPORT = ":5050"
)

func loadSettings() map[string]interface{} {
	buf, err := ioutil.ReadFile("./settings/settings.json")
	if err != nil {
		log.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(buf, &m)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func setupEnv() *common.Env {
	settings := loadSettings()
	debug := false
	if len(os.Args) > 1 {
		for _, v := range os.Args[1:] {
			if v == "debug" {
				debug = true
			}
		}
	}
	database := settings["DSN"].(string)

	db, err := datastore.NewDataStore(database)

	if err != nil {
		log.Fatal(err)
	}
	return &common.Env{Debug: debug, AccessToken: settings["AccessToken"].(string), DB: *db}
}

func main() {
	// -----------------------------------
	env := setupEnv()
	defer env.DB.Close(context.Background())
	// -----------------------------------

	r := mux.NewRouter()

	// -----------------------------------

	r.Handle("/", handlers.HomeHandler(env)).Methods("GET")
	r.Handle("/store", handlers.HomeHandler(env)).Methods("GET")
	r.Handle("/store/results", handlers.ResultsHandler(env)).Methods("GET")
	r.Handle("/store/categories", handlers.CategoriesHandler(env)).Methods("GET")
	r.Handle("/store/details/{uid}", handlers.DetailsHandler(env)).Methods("GET")
	r.Handle("/store/service", handlers.ServiceHandler(env)).Methods("GET")
	r.Handle("/store/cart", handlers.CartHandler(env)).Methods("GET")
	r.Handle("/store/calculator", handlers.CalculatorHandler(env)).Methods("GET")

	// -----------------------------------

	r.Handle("/store/api/search", handlers.ApiSearchHandler(env)).Methods("GET", "POST")
	r.Handle("/store/api/slide", handlers.ApiSlideHandler(env)).Methods("POST")
	r.Handle("/store/api/category", handlers.ApiCategoryHandler(env)).Methods("POST")
	r.Handle("/store/api/store", handlers.ApiStoreHandler(env)).Methods("GET", "POST")
	r.Handle("/store/api/sell", handlers.ApiSellHandler(env)).Methods("GET", "POST")

	// -----------------------------------

	//loggedRouter := ghandlers.LoggingHandler(os.Stdout, r)
	// http.Handle("/", stdChain.Then(loggedRouter))

	stdChain := alice.New(middleware.PanicRecoveryHandler, middleware.RemoveTrailingSlashHandler)
	http.Handle("/", stdChain.Then(r))

	r.PathPrefix("/store/static/").Handler(http.StripPrefix("/store/static/", http.FileServer(http.Dir("./static"))))

	err := http.ListenAndServe(WEBSERVERPORT, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
