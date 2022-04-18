package gocol

import (
	"fmt"
)

type Rad struct {
	gr string
	Grq string
	num   int
	lemme *Lemme
}

func (r Rad) doc() string {
	return fmt.Sprintf("grq %s gr %s  num %d modele %s, lemme %s",
		r.Grq, r.gr, r.num, r.lemme.Modele.id, r.lemme.Grq)
}

var radicaux = make(map[string][]*Rad)

func ajRadicaux() {
	for _, lem := range Lemmes {
		for _, lrad := range lem.Radicaux {
			for _, rad := range lrad {
				radicaux[rad.gr] = append(radicaux[rad.gr], rad)
			}
		}
	}
}
