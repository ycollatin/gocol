package gocol

import (
	"fmt"
	"strings"

	//"log"
)

type Lemme struct {
	Grq,
	Gr []string
	Cle,
	Indmorph,
	Pos,
	Genre,
	Traduction string
	nh,
	Freq int // fréquence en dernier champ
	Radicaux map[int][]*Rad
    renvoi string
	Modele   *Modele
	// TODO ajouter une map de traductions
}

func (l *Lemme) Doc() string {
	var lr []string
	for nr, rr := range l.Radicaux {
		for _, r := range rr {
			lr = append(lr, fmt.Sprintf("n° %d %s", nr, r.Grq))
		}
	}
	return fmt.Sprintf("clé %s %s, %s-%s : %s\n  rad. %s",
		l.Cle,
		strings.Join(l.Grq, ","), l.Indmorph, l.Pos, l.Traduction,
		strings.Join(lr, "\n       "))
}

func (l Lemme) habetRad(r string, n int) bool {
	for _, rad := range l.Radicaux[n] {
		if rad.gr == r {
			return true
		}
	}
	return false
}

var Lemmes = make(map[string]*Lemme)

func recupNh(k string) (string, string) {
	der := string(k[len(k)-1])
	if der == "2" || der == "3" || der == "4" {
		k = strings.TrimSuffix(k, der)
	}
	return k, der
}

// creeLemme(l string)
// créateur de lemme à partir de la ligne l de lemmes.la
func creeLemme(l string) *Lemme {
	var lem *Lemme = new(Lemme)
	lem.Radicaux = make(map[int][]*Rad)
	eclats := strings.Split(l, "|")
	for i, e := range eclats {
		switch i {
		case 0:
			eclg := strings.Split(e, "=")
			// supprimer et affecter le numéro d'homonymie
			var lgrq, lnh string
			lgrq, lnh = recupNh(eclg[0])
			lem.nh = Strtoint(lnh)
			lem.Grq = append(lem.Grq, lgrq)
			lem.Gr = append(lem.Gr, atone(eclg[0]))
			if len(eclg) > 1 {
				eclk := strings.Split(eclg[1], ",")
				for i, eclki := range eclk {
					lem.Grq = append(lem.Grq, eclki)
					lem.Gr = append(lem.Gr, atone(lem.Grq[i]))
				}
			}
			lem.Cle = atone(lem.Grq[0])
			if lem.nh > 1 {
				lem.Cle = lem.Cle + lnh
			}

		case 1:
			lem.Modele = modeles[e]
			lem.Pos = lem.Modele.pos
			if strings.ToLower(lem.Cle) != lem.Cle && lem.Pos == "n" {
				lem.Pos = "np"
			}
		case 2, 3:
			if e > "" {
				eclR := strings.Split(e, ",")
				for _, r := range eclR {
					rgr := deramAtone(r)
					rad := new(Rad)
					rad.Grq = r
					rad.gr = rgr
					rad.num = i - 1
					rad.lemme = lem
					lem.Radicaux[rad.num] = append(lem.Radicaux[rad.num], rad)
				}
			}
			// radicaux calculés
			for _, grad := range lem.Modele.lgenR {
				// si un radical de même num a été donné
				// on ne tient pas compte du calculé
				if len(lem.Radicaux[grad.num]) > 0 {
					continue
				}
				for _, grq := range lem.Grq {
					rgrq := grad.rad(grq)
					rgr := deramAtone(rgrq)
					if !lem.habetRad(rgr, grad.num) {
						rad := new(Rad)
						rad.Grq = rgrq
						rad.gr = rgr
						rad.num = grad.num
						rad.lemme = lem
						lem.Radicaux[rad.num] = append(lem.Radicaux[rad.num], rad)
					}
				}
			}
		case 4: // indMorph
			lem.Indmorph = e
            // renvois
            eclcf := strings.Split(e, "cf. ")
            if len(eclcf) > 1 {
                eclr := strings.Split(eclcf[1], " ")
                lem.renvoi = eclr[0]
            }

			if strings.Contains(lem.Indmorph, " f.") {
				lem.Genre = "féminin"
			} else if strings.Contains(lem.Indmorph, " m.") {
				lem.Genre = "masculin"
			} else if strings.Contains(lem.Indmorph, " n.") {
				lem.Genre = "neutre"
			}
			// pos des prépositions, négations et adverbes
			cacc := contient(lem.Indmorph, "+ acc.")
			cabl := contient(lem.Indmorph, "+ abl.")
			switch {
			case cacc && cabl:
				lem.Pos = "prepAA"
			case cabl:
				lem.Pos = "prepAbl"
			case cacc:
				lem.Pos = "prepAcc"
			case contient(lem.Indmorph, "neg."):
				lem.Pos = "neg"
			case contient(lem.Indmorph, "interj."):
				lem.Pos = "intj"
			case contient(lem.Indmorph, "adv."):
				lem.Pos = "adv"
			case contient(lem.Indmorph, "conj."):
				lem.Pos = "conj"
			}
		case 5: // fréquence
			lem.Freq = Strtoint(e)
		}
	}
	return lem
}

// lisLemmes()
// lecteur de la base de lemmes
func lisLemmes(nf string) {
	ll := Lignes(nf)
	for _, l := range ll {
		lem := creeLemme(l)
		// si le lemme existe déjà, passer)
		present, _ := Lemmes[lem.Cle]
		if present == nil {
			Lemmes[lem.Cle] = lem
		}
	}
}

// lisTraductions()
// lis les traductions à martir d'un fichier lemmes.??
func lisTraductions(nf string) {
	ll := Lignes(nf)
	for _, l := range ll {
		ecl := strings.Split(l, ":")
		cle := ecl[0]
		lem := Lemmes[cle]
		if lem != nil && lem.Traduction == "" {
			lem.Traduction = strings.Join(ecl[1:], ":")
		} /* else { log.Println(l) } */
	}
    // traduction des renvois
    for _, l := range Lemmes {
        if l.renvoi > "" {
            lrenv := Lemmes[l.renvoi]
            if lrenv != nil {
                l.Traduction = Lemmes[l.renvoi].Traduction
            } /*else {
                log.Println(l.Cle,"l.Traduction, NIL")
            }*/
        }
    }
}
