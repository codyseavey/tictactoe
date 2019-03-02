package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

var (
	indexTemplate = template.Must(template.New("index.html").Funcs(template.FuncMap{
		"formatTile": func(input Tile) string {
			ret := "_"
			if input == XVar {
				ret = "X"
			} else if input == OVar {
				ret = "O"
			}
			return ret
		},
	}).Parse(`
		<style>
			table,td{text-align:center}table{width:40%;border-collapse:collapse}
			.cell{font:150px arial,sans-serif;text-decoration:none;color:#000}
			td{height:200px;width:33.333%;vertical-align:center;border:6px solid #222}
			td:first-of-type,td:nth-of-type(2),td:nth-of-type(3){border-top-color:transparent}
			td:first-of-type{border-left-color:transparent}td:nth-of-type(3){border-right-color:transparent}
			tr:nth-of-type(3) td{border-bottom-color:transparent}
		</style>
		<table align="center">
		    {{range $i, $b := $.Board }}
		    <tr>
		        {{range $j, $c := index $.Board $i}}
		           <td><a class="cell" href="/updateGame/{{ $.SessionID }}/{{ $i }}/{{ $j }}">{{ formatTile $c }}</a></td>
		        {{end}}
		    </tr>
		    {{end}}
		</table>
		<p>Currently it is {{ formatTile .Turn }}'s turn.</p>
		<p>Winner: {{ formatTile .Winner }}</p>
		<p>New Game? <a href=/newGame/1>1p</a> <a href=/newGame/2>2p</a></p>
	`))
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<p>New game? <a href=/newGame/1>1p</a> <a href=/newGame/2>2p</a>`))
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
	indexTemplate.Execute(w, g)
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

// redirect http to https
func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		http.StatusTemporaryRedirect)
}

func main() {
	log.SetOutput(os.Stderr)

	initDb()

	// go http.ListenAndServe(":8080", http.HandlerFunc(redirect))

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
