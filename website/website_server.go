package website

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var router *mux.Router

//StartServer is the entry point for this package
func StartServer(port string) {
	router := mux.NewRouter().StrictSlash(true)

	setupRoutes(router)

	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("website/templates/assets/"))),
	)

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	log.Printf("Serving on https://0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func setupRoutes(mRouter *mux.Router) {

	mRouter.Path("/").HandlerFunc(homeHandle)

	mRouter.Path("/login").HandlerFunc(login)
	mRouter.Path("/callback").HandlerFunc(callback)

	mRouter.Path("/queue").HandlerFunc(showqueue)
	mRouter.Path("/queue/add/{id}").HandlerFunc(queueadd)
	mRouter.Path("/queue/remove/{id}").HandlerFunc(queueRemove)
	mRouter.Path("/queue/up/{id}").HandlerFunc(queueUp)
	mRouter.Path("/queue/down/{id}").HandlerFunc(queueDown)
	mRouter.Path("/queue/skip").HandlerFunc(queueSkip)
	mRouter.Path("/queue/state").HandlerFunc(pauseQueue)

}
