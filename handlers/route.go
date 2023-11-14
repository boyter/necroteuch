package handlers

import (
	"embed"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"io/fs"
	"necroteuch/common"
	"necroteuch/service"
	"net/http"
)

type Application struct {
	Service     *service.Service
	StaticFiles embed.FS
}

func NewApplication(service *service.Service, staticFiles embed.FS) (Application, error) {
	application := Application{
		Service:     service,
		StaticFiles: staticFiles,
	}
	err := application.ParseTemplates()
	return application, err
}

func (app *Application) Routes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// JSON routes
	router.HandleFunc("/health-check/", app.HealthCheck).Methods("GET")

	// ActivityPub routes MUST NOT HAVE / ON END!!!!
	router.HandleFunc("/.well-known/webfinger", app.WebFinger).Methods("GET")
	//router.Handle("/u/{username:.*?}/inbox", IpRestrictorHandler(app.ActivityPubInbox)).Methods("POST") // route where messages and stuff get posted NEVER GET
	//router.Handle("/u/{username:.*?}/outbox", IpRestrictorHandler(app.ActivityPubOutbox)).Methods("GET")
	//router.Handle("/u/{username:.*?}/followers", IpRestrictorHandler(app.Empty)).Methods("GET")
	//router.Handle("/u/{username:.*?}/following", IpRestrictorHandler(app.Empty)).Methods("GET")
	//router.Handle("/u/{username:.*?}/image", IpRestrictorHandler(app.Image)).Methods("GET")
	//router.Handle("/u/{username:.*?}", IpRestrictorHandler(app.ActivityPubPerson)).Methods("GET", "POST")

	// Regular routes
	router.HandleFunc("/", app.Index).Methods("GET")

	// Setup to serve files from the supplied directory
	staticFS := fs.FS(app.StaticFiles)
	// strip off the location such that we route correctly
	staticContent, err := fs.Sub(staticFS, "assets/static")
	if err != nil {
		log.Fatal().Str(common.UniqueCode, "f8819e2a").Err(err).Msg("error with static content")
	}
	f := http.FileServer(http.FS(staticContent))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", f))

	return router
}
