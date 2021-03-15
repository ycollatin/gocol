/*     collatinus - moteur.go   */
package gocol

import (
	"fmt"
	"strings"
)

type Sr struct {
	Lem     *Lemme
	Nmorph  []int
	Morphos []string
}

// Res Collection de résultats
type Res []Sr

func Restostring(r Res) string {
	var lr []string
	for _, rl := range r {
		if rl.Lem == nil {
			continue
		}
		l := fmt.Sprintf("   %s, %s : %s",
			strings.Join(rl.Lem.Grq, " "),
			rl.Lem.Indmorph,
			rl.Lem.Traduction)
		lr = append(lr, l)
		for _, m := range rl.Morphos {
			//l = l + ("\n      " + m)
			lr = append(lr, "      "+m)
		}
		//lr = append(lr, l)
	}
	return strings.Join(lr, "\n")
}

// addRes(r Res, l *Lemme, m string)
//  ajout d'une analyse lemme + morpho  à un Sr
func AddRes(r Res, l *Lemme, m string, nm int) Res {
	var contient bool
	for i := 0; i < len(r); i++ {
		if r[i].Lem == l {
			contient = true
			r[i].Morphos = append(r[i].Morphos, m)
			r[i].Nmorph = append(r[i].Nmorph, nm)
		}
	}
	if !contient {
		var rs Sr
		rs.Lem = l
		rs.Morphos = []string{m}
		rs.Nmorph = []int{nm}
		r = append(r, rs)
	}
	return r
}

/*
   Lemmatisation d'une seule forme f
*/

func LemmatiseF(f string) (result Res) {
	f = strings.TrimSuffix(f, "\n")
	r := f
	d := ""

	// irréguliers
	irr := irregs[f]
	if irr != nil {
		for _, nm := range irr.lmorph {
			result = AddRes(result, irr.lem, irr.grq+" "+Morphos[nm], nm)
			if irr.exclusif {
				continue
			}
		}
	}

	// romains
	if estRomain(f) {
		lin := fmt.Sprintf("%s|inv|||num.|1", f)
		romain := creeLemme(lin)
		romain.Traduction = aRomano(f)
		result = AddRes(result, romain, "inv\n", 0)
	}

	// radical-désinence
	for {
		lrad := radicaux[r]
		// contraction ii > i
		ii := strings.HasPrefix(d, "i") && !strings.HasPrefix(d, "ii") && !strings.HasSuffix(r, "i")
		if ii {
			lradi := radicaux[r+"i"]
			lrad = append(lrad, lradi...)
		}
		if len(lrad) > 0 {
			for _, rad := range lrad {
				for _, des := range rad.lemme.modele.desm[rad.num] {
					if des.gr == d {
						/*
						m := fmt.Sprintf("%s%s %s %s",
							rad.grq, des.grq, Morphos[des.morpho], rad.lemme.Genre)
						*/
						mg := Morphos[des.morpho]
						if rad.lemme.Genre != "" {
							mg = genreNom(mg, rad.lemme.Genre)
						}
						m := fmt.Sprintf("%s%s %s", rad.grq, des.grq, mg)
						result = AddRes(result, rad.lemme, m, des.morpho)
					}
				}
			}
		}
		l := len(r) - 1
		if r == "" {
			break
		}
		d = string(r[l]) + d
		r = r[:l]
	}
	return
}

// func Lemmatise(f string) (lsr Res, echec bool) {
// calcule toutes les variantes graphiques de f, lemmatise chacune et
// retourne avec leur analyse celles qui ont obtenu une lemmatisation
// renvoie une liste vide et echec à true si aucune lemmatisation
// n'a été trouvée
func Lemmatise(f string) (lsr Res, echec bool) {
	f = Deramise(f)
	liste := varsF(f)
	for _, el := range liste {
		if el == "" {
			continue
		}
		lel := LemmatiseF(el)
		if len(lel) > 0 {
			for _, l := range lel {
				lsr = append(lsr, l)
			}
		} else {
			echec = true
		}
	}
	echec = len(lsr) == 0
	return
}
