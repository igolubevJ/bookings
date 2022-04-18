package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/igolubevJ/bookings/internal/config"
	"github.com/igolubevJ/bookings/internal/handlers"
	"github.com/igolubevJ/bookings/internal/models"
	"github.com/igolubevJ/bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":5000"

var app config.AppConfig
var session *scs.SessionManager

// main is the start applcation point
func main() {
	// ! change this to true when in production
	app.InProduction = false

	 // What am I going to put in the session
	 gob.Register(models.Reservation{})

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCash()
	if err != nil {
		log.Fatal("can not create template cash")
	}

	app.TemplateCache = tc
	app.UseCache = false
	
	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Application Server is startup on http://localhost%s", portNumber))
	
	srv := &http.Server{
		Addr: portNumber,
		Handler: routes(app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
