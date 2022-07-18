package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/RasmusLindroth/tut-web/web"
	"github.com/gorilla/mux"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
)

type App struct {
	Templates map[string]*template.Template
}

type PageData struct {
	Pages       []Page
	CurrentPage string
	CSS         string
}

type Page struct {
	Name string
	URL  string
}

func main() {
	port := ""
	portEnv, portEnvSet := os.LookupEnv("port")

	if portEnvSet {
		port = portEnv
	}

	portFlag := flag.String("port", "", "web server port")
	flag.Parse()

	if *portFlag != "" {
		port = *portFlag
	}

	if port == "" {
		fmt.Println("You need to set a port. Use env variable or flag.")
		os.Exit(1)
	}

	baseData := PageData{
		Pages: []Page{
			{"About", "/"},
			{"Install", "/install"},
			{"Config", "/config"},
			{"Keys", "/keys"},
			{"Commands", "/commands"},
			{"Flags", "/flags"},
			{"Templates", "/templates"},
			{"Contact", "/contact"},
		},
	}
	app := &App{
		Templates: make(map[string]*template.Template),
	}
	sites := []string{"about", "install", "config", "keys", "commands", "flags", "templates", "contact"}
	for _, s := range sites {
		tmpl, err := template.ParseFS(web.Templates, "templates/base.tmpl", "templates/"+s+".tmpl")
		if err != nil {
			log.Fatalln(err)
		}
		app.Templates[s] = tmpl
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	style, err := web.Static.Open("static/style.css")
	if err != nil {
		log.Fatalln(err)
	}
	b, err := ioutil.ReadAll(style)
	if err != nil {
		log.Fatalln(err)
	}
	styleB, err := m.Bytes("text/css", b)
	if err != nil {
		log.Fatalln(err)
	}
	baseData.CSS = string(styleB)

	preview, err := web.Static.Open("static/preview.png")
	if err != nil {
		log.Fatalln(err)
	}
	c, err := ioutil.ReadAll(preview)
	if err != nil {
		log.Fatalln(err)
	}

	theme, err := web.Static.Open("static/theme.html")
	if err != nil {
		log.Fatalln(err)
	}
	t, err := ioutil.ReadAll(theme)
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		d := baseData
		d.CurrentPage = "About"
		app.Templates["about"].ExecuteTemplate(w, "base", d)
	})
	r.HandleFunc("/install", func(w http.ResponseWriter, r *http.Request) {
		d := baseData
		d.CurrentPage = "Install"
		app.Templates["install"].ExecuteTemplate(w, "base", d)
	})
	r.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		d := baseData
		d.CurrentPage = "Config"
		app.Templates["config"].ExecuteTemplate(w, "base", d)
	})
	r.HandleFunc("/keys", func(w http.ResponseWriter, r *http.Request) {
		d := baseData
		d.CurrentPage = "Keys"
		app.Templates["keys"].ExecuteTemplate(w, "base", d)
	})
	r.HandleFunc("/commands", func(w http.ResponseWriter, r *http.Request) {
		d := baseData
		d.CurrentPage = "Commands"
		app.Templates["commands"].ExecuteTemplate(w, "base", d)
	})
	r.HandleFunc("/flags", func(w http.ResponseWriter, r *http.Request) {
		d := baseData
		d.CurrentPage = "Flags"
		app.Templates["flags"].ExecuteTemplate(w, "base", d)
	})
	r.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		d := baseData
		d.CurrentPage = "Templates"
		app.Templates["templates"].ExecuteTemplate(w, "base", d)
	})
	r.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		d := baseData
		d.CurrentPage = "Contact"
		app.Templates["contact"].ExecuteTemplate(w, "base", d)
	})
	r.HandleFunc("/theme", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/html")
		w.Write(t)
	})
	r.HandleFunc("/static/preview.png", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "image/png")
		w.Write(c)
	})
	http.Handle("/", r)

	http.ListenAndServe(":"+port, r)
}
