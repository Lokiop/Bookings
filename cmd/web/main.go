package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Lokiop/Bookings/internal/config"
	"github.com/Lokiop/Bookings/internal/drivers"
	"github.com/Lokiop/Bookings/internal/handlers"
	"github.com/Lokiop/Bookings/internal/helpers"
	"github.com/Lokiop/Bookings/internal/models"
	"github.com/Lokiop/Bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portnumber = ":8000"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

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

func run() (*drivers.DB, error) {
	//what I am going to store in the session
	gob.Register(models.User{})
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.ErrorPage{})

	//Changers to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to database
	log.Println("Connecting to Database....")
	db, err := drivers.ConnectSQL("host=localhost port=5432 dbname=Bookings user=postgres password=D@rshanheda24")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying....")
	}
	log.Println("Connected to Database....")

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Cannot Create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.Usecache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
