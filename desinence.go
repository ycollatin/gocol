package gocol

import "fmt"

type Des struct {
	gr,
	Grq string
	Morpho int
	modele *Modele
	Nr     int      // numéro de radical
}

func (d Des) doc() string {
	return fmt.Sprintf("des %s -> %s num %d modele %s %s",
		d.Grq, d.gr, d.Nr, d.modele.id, Morphos[d.Morpho])
}

func (d Des) clone() (dc *Des) {
	dc = new(Des)
	dc.gr = d.gr
	dc.Grq = d.Grq
	dc.Morpho = d.Morpho
	dc.modele = d.modele
	dc.Nr = d.Nr
	return dc
}

var desinences = make(map[string][]*Des)

// creeDes( g graphie, md modèle, mr morpho, n numéro de radical)
func creeDes(g string, md *Modele, mr, n int) *Des {
	var d *Des = new(Des)
	if g == "-" {
		d.Grq = ""
		d.gr = ""
	} else {
		d.Grq = g
		d.gr = deramAtone(g)
	}
	d.Morpho = mr
	d.modele = md
	d.Nr = n
	return d
}
