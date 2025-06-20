package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andrii-mykhailiuk/bookings/pkg/config"
	"github.com/andrii-mykhailiuk/bookings/pkg/handlers"
	"github.com/andrii-mykhailiuk/bookings/pkg/render"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager

func main() {
	App()
}

// App mai application function
func App() {

	// change this to true when in producton
	appConfig.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	appConfig.TemplateCache = tc
	appConfig.UseCache = false

	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)

	render.NewTemplates(&appConfig)

	// _ = http.ListenAndServe(portNumber, nil)
	fmt.Println(fmt.Sprintf("Starting application on port: %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
