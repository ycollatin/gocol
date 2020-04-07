package gocol

import "fmt"

type Des struct {
	gr,
	grq string
	morpho int
	modele *Modele
	nr     int
}

func (d Des) doc() string {
	return fmt.Sprintf("des %s -> %s num %d modele %s %s",
		d.grq, d.gr, d.nr, d.modele.id, Morphos[d.morpho])
}

func (d Des) clone() (dc *Des) {
	dc = new(Des)
	dc.gr = d.gr
	dc.grq = d.grq
	dc.morpho = d.morpho
	dc.modele = d.modele
	dc.nr = d.nr
	return dc
}

var desinences = make(map[string][]*Des)

// creeDes( g graphie, md modèle, mr morpho, n numéro de radical)
func creeDes(g string, md *Modele, mr, n int) *Des {
	var d *Des = new(Des)
	if g == "-" {
		d.grq = ""
		d.gr = ""
	} else {
		d.grq = g
		d.gr = deramAtone(g)
	}
	d.morpho = mr
	d.modele = md
	d.nr = n
	return d
}
