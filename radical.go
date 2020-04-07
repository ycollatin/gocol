package gocol

import (
	"fmt"
)

type Rad struct {
	gr,
	grq string
	num   int
	lemme *Lemme
}

func (r Rad) doc() string {
	return fmt.Sprintf("grq %s gr %s  num %d modele %s, lemme %s",
		r.grq, r.gr, r.num, r.lemme.modele.id, r.lemme.Grq)
}

var radicaux = make(map[string][]*Rad)

func ajRadicaux() {
	for _, lem := range lemmes {
		for _, lrad := range lem.radicaux {
			for _, rad := range lrad {
				radicaux[rad.gr] = append(radicaux[rad.gr], rad)
			}
		}
	}
}
