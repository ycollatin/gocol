package gocol

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Lsr struct {
	Forme string
	Llem  Res
}

func lemm(w http.ResponseWriter, r *http.Request) {
	var lres []Lsr
	r.ParseForm()
	p := r.Form["texte"][0]
	lm := Mots(p)
	for _, m := range lm {
		lemm, _ := Lemmatise(m)
		nl := Lsr{m, lemm}
		lres = append(lres, nl)
	}
	t, err := template.New("tlem").Funcs(template.FuncMap{
		"Deref": func(lem *Lemme) Lemme {
			return *lem
		},
	}).Parse(
		`<!DOCTYPE HTML>
       <html>
       <head>
         <meta charset="utf-8">
       </head>
       <body>
         <h2>résultats</h2>
         <div>
         <a href= "/">nouvelle recherche</a><br/>
         {{range .}}
            <p><em>{{.Forme}}</em>
            <ul>
            {{range $i, $l := .Llem}}
              {{$d := (Deref $l.Lem)}}
              <li>
              <strong>{{range $d.Grq}}
                      &nbsp;{{.}}
                      {{end}}
              </strong> {{$d.Indmorph}} : {{$d.Traduction}}
              <ul>
              {{range .Morphos}}
                <li>{{.}}</li> 
              {{end}}
              </ul>
              </li>
            {{end}}
            </ul>
         {{end}}
         </div>
       </body>
       </head>
       </html>`)
	if err != nil {
		fmt.Println(err)
	} else {
		e := t.Execute(w, lres)
		if e != nil {
			log.Fatal(e)
		}
	}
}

func req(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("index.gtpl"))
		t.Execute(w, nil)
	}
}

func Serveur(port string) {
	http.HandleFunc("/lem", lemm)
	http.HandleFunc("/", req)
	fmt.Println("serveur port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

/*
func serveur(port string) {
    m := http.NewServeMux()
    s := http.Server{Addr: ":"+port, Handler: m}
    http.HandleFunc("/lem", lemm)
    http.HandleFunc("/", req)
    m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
        s.Shutdown(context.Background())
    })
    if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }
    log.Printf("Finished")
}
*/
