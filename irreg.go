package gocol

import (
	"strings"
)

type Irr struct {
	grq,
	gr string
	lem      *Lemme
	lmorph   []int
	exclusif bool
}

func creeIrr(l string) (irr *Irr) {
	ecl := strings.Split(l, ":")
	irr = new(Irr)
	irr.grq = ecl[0]
	if d := len(irr.grq) - 1; irr.grq[d] == '*' {
		irr.grq = irr.grq[:d]
		irr.exclusif = true
	}
	irr.gr = atone(irr.grq)
	irr.lem = lemmes[ecl[1]]
	irr.lmorph = listei(ecl[2])
	return irr
}

var irregs = make(map[string]*Irr)

func lisIrregs(nf string) {
	ll := Lignes(nf)
	for _, l := range ll {
		nirr := creeIrr(l)
		// L'irrégulier n'est ajouté que s'il
		// n'est pas encore dans la base.
		if irregs[nirr.gr] == nil {
			irregs[nirr.gr] = nirr
		}
	}
}

/*
   lem := lemmes[cle]
   if lem != nil && lem.Traduction == "" {
       lem.Traduction = strings.Join(ecl[1:], ":")
   }
*/
