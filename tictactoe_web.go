package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	tttTemplate = template.Must(template.New("ttt.html").Funcs(template.FuncMap{
		"formatTile": func(input Tile) string {
			ret := "_"
			if input == XVar {
				ret = "X"
			} else if input == OVar {
				ret = "O"
			}
			return ret
		},
	}).ParseFiles("assets/templates/ttt.html"))
	indexTemplate = template.Must(template.New("index.html").ParseFiles("assets/templates/index.html"))
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, nil)
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	players, err := strconv.Atoi(vars["players"])
	if err != nil {
		log.Println(err)
	}
	id, _ := NewGame(players)
	http.Redirect(w, r, "/ttt/"+strconv.Itoa(id), http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
	}
	g := GetGame(id)
	tttTemplate.Execute(w, g)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
	}
	rowInt, err := strconv.Atoi(vars["row"])
	if err != nil {
		log.Println(err)
	}
	columnInt, err := strconv.Atoi(vars["column"])
	if err != nil {
		log.Println(err)
	}
	row := Tile(rowInt)
	column := Tile(columnInt)
	g := GetGame(id)
	g, err = g.Update(row, column) // Update with the users choice
	if err == nil && g.Players == 1 {
		computerRow, computerColumn, err := g.getComputerChoice()
		if err == nil {
			g.Update(computerRow, computerColumn) // Update with our choice (O)
		}
	}

	http.Redirect(w, r, "/ttt/"+strconv.Itoa(id), http.StatusFound)
}

func main() {
	log.SetOutput(os.Stdout)

	initDb()

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/newGame/{players:[1-2]}", newGameHandler)
	r.HandleFunc("/ttt/{id:[0-9]+}", viewHandler)
	r.HandleFunc(
		"/updateGame/{id:[0-9]+}/{row:[0-9]+}/{column:[0-9]+}",
		updateHandler,
	)
	log.Fatal(http.ListenAndServeTLS(":8443",
		"/certs/cert.pem",
		"/certs/privkey.pem", r),
	)
}
