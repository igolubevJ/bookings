package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/igolubevJ/bookings/pkg/config"
	"github.com/igolubevJ/bookings/pkg/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates set the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all template
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)

	return td
}

// RenderTemplate renders template using html/template
func RenderTemplate(w http.ResponseWriter, r *http.Request,tmpl  string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cash from app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCash()
	}
	
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cash")
	} 

	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}
}

// Creates a template cash as a map
func CreateTemplateCash() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			fmt.Println("One place for generate error")
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			fmt.Println("Two place for generate error")
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err := ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				fmt.Println("Three place for generate error")
				return myCache, err
			}

			myCache[name] = ts
		}
	}

	return myCache, nil
}
