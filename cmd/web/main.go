package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Lokiop/Bookings/internal/config"
	"github.com/Lokiop/Bookings/internal/handlers"
	"github.com/Lokiop/Bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portnumber = ":8000"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	//Changers to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Cannot Create template cache")
	}

	app.TemplateCache = tc
	app.Usecache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplate(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println("Server Starting at port", portnumber)
	// _ = http.ListenAndServe(portnumber, nil)

	srv := &http.Server{
		Addr:    portnumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
