package gocol

import (
	"fmt"
	"strings"
)

type Modele struct {
	id    string
	Desm  map[int][]*Des    // table des désinences classées par n° de morpho
	lgenR map[int]*Genrad   // table des radicaux classés par n° de radical
	pere  *Modele
	abs   []int
	suf   []string
	pos   string
	sufd  string
}

var modeles = make(map[string]*Modele)
var vardes = make(map[string][]string)

func (m Modele) doc() string {
	var mid string
	if m.pere != nil {
		mid = m.pere.id
	} else {
		mid = "nil"
	}
    var gg string
    for _, r := range m.lgenR {
        gg = gg + "\n  "+ r.doc()
    }
    var dd string
    for _, ld := range m.Desm {
        for _, d := range ld {
            dd = dd + "\n"+d.doc()
        }
    }
	return fmt.Sprintf("modele %s père %s, pos %s sufd %s, rads %s, des %s",
		m.id, mid, m.pos, m.sufd, gg, dd)
}

// m a-t-il une désinence de morpho d
func (m Modele) habetdes(d int) bool {
    for _, v := range m.Desm {
        for _,w := range v {
            if w.Morpho == d {
                return true
            }
        }
    }
    return false
}

// vrai si le modèle m a déjà le générateur de radical gnr
func (m *Modele) habetR(r int) bool {
    _, in := m.lgenR[r]
    return in
}

func (m *Modele) estabs(d int) bool {
    for _, i := range m.abs {
        if i == d {
            return true
        }
    }
    return false
}

// renvoie toutes les désinences de morpho morf
func (m *Modele) desmorph(morf int) []*Des {
    var dd []*Des
    for _, v := range m.Desm {
        for _, d := range v {
            if d.Morpho == morf {
                dd = append(dd, d)
            }
        }
    }
    return dd
}

// héritage du pos et des générateurs de radicaux
// les désinences sont héritées via heritedes()
func (m *Modele) herite() {
	if m.pere == nil {
		return
	}
	// héritage du pos
	if m.pos == "" {
		m.pos = m.pere.pos
	}
	/*  héritage des générateurs de radicaux */
	for k, genr := range m.pere.lgenR {
		if !m.habetR(k) {
			m.lgenR[k] = genr
		}
	}
}

// héritage des désinences, appelé séparément après
// les redéfinitions des désinences du modèle
func (m *Modele) heritedes() {
    if m.pere == nil {
        return
    }
    // héritage des desinences non absentes non redéfinies
    for n, ld := range m.pere.Desm {
        for _, d := range ld {
            if !m.estabs(d.Morpho) && !m.habetdes(d.Morpho) {
                nd := d.clone()
                nd.modele = m
                m.Desm[n] = append(m.Desm[n], nd)
            }
        }
    }
}

func (m *Modele) ajsuffd() {
	if m.sufd == "" {
		return
	}
	for _, ld := range m.Desm {
		for _, d := range ld {
			d.Grq = d.Grq + m.sufd
			d.gr = atone(d.Grq)
		}
	}
}

func (m *Modele) ajsuff() {
	for _, lsuf := range m.suf {
		ecl := strings.Split(lsuf, ":")
		li := listei(ecl[0])
		suff := ecl[1]
		for _, i := range li {
			for _, ld := range m.Desm {
				for _, d := range ld {
					{
						if d.Morpho == i {
							nd := d.clone()
							nd.Grq = d.Grq + suff
							nd.gr = atone(nd.Grq)
							m.Desm[d.Nr] = append(m.Desm[d.Nr], nd)
						}
					}
				}
			}
		}
	}
}

func lismodeles(nf string) {
	ll := Lignes(nf)
	var m *Modele
	for _, l := range ll {
		if l == "" {
			continue
		}
        // variables
		if strings.HasPrefix(l, "$") {
			l = strings.TrimPrefix(l, "$")
			ecl := strings.Split(l, "=")
			k := ecl[0]
			vardes[k] = strings.Split(ecl[1], ";")
			continue
		}
		ecl := strings.Split(l, ":")
		cle := ecl[0]
		val := strings.TrimPrefix(l, cle+":")
		switch cle {
		case "modele":
            // terminer le modèle précédent
			if m != nil {
                m.heritedes()
				m.ajsuffd()
				m.ajsuff()
				modeles[m.id] = m
			}
			m = new(Modele)
			m.pere = nil
			m.id = val
            m.lgenR = make(map[int]*Genrad)
			m.Desm = make(map[int][]*Des)
		case "R":
			num := Strtoint(ecl[1])
            switch ecl[2] {
            case "K":
                m.lgenR[num] = &Genrad{num, 0, ""}
            case "-":
                delete(m.lgenR, num)
            default:
                lp := strings.Split(ecl[2], ",")
                oter := Strtoint(lp[0])
                ajout := lp[1]
                if ajout == "0" {
                    ajout = ""
                }
                m.lgenR[num] = &Genrad{num, oter, ajout}
            }
		case "abs":
			m.abs = listei(ecl[1])
		case "des", "des+":
            // établir la liste des morphos décrites
			li := listei(ecl[1])
            // et le numéro de radical des des de la ligne
			nr := Strtoint(ecl[2])
            // partie désinence
			ld := ecl[3]
            // dd, liste brute des graphies des désinences
			var dd []string = strings.Split(ld, ";")
            // ddd, liste dont les variables ont été développées 
			var ddd []string
			for _, sdes := range dd {
                // cas particulier : variables, suffixes
				if strings.HasPrefix(sdes, "$") {
					ddv := strings.TrimPrefix(sdes, "$")
					ldv := vardes[ddv]
					for _, dv := range ldv {
						ddd = append(ddd, dv)
					}
				} else if strings.Contains(sdes, "$") {
					fix := strings.Split(sdes, "$")
					prefix := fix[0]
					suffix := fix[1]
					ldv := vardes[suffix]
					for _, dv := range ldv {
						ddd = append(ddd, prefix+dv)
					}
				} else {
                    // ajout des graphies sans variable 
					if sdes == "-" {
						sdes = ""
					}
					ddd = append(ddd, sdes)
				}
			}
            // mise en map les dés surnuméraires (après maxd)
            // sont des copies de la dernière
			maxd := len(ddd)
			for ides, ili := range li {
				if ides < maxd {
                    // sld: liste (sep ',') des variables d'une morpho
					sld := ddd[ides]
					ecld := strings.Split(sld, ",")
					for _, cld := range ecld {
                        // cld graphie
                        // m   modèle
                        // mr  n° morpho
                        // nr  numéro de radical)
                        nd := creeDes(cld, m, ili, nr)
                        nd.aj = cle == "des+"
						m.Desm[nr] = append(m.Desm[nr], nd)
					}
				} else {
					sld := ddd[maxd-1]
					ecld := strings.Split(sld, ",")
					for _, cld := range ecld {
                        nd := creeDes(cld, m, ili, nr)
                        nd.aj = cle == "des+"
						m.Desm[nr] = append(m.Desm[nr], nd)
                    }
                }
			}
		case "pos":
			m.pos = val
		case "suf":
			m.suf = append(m.suf, val)
		case "sufd":
			m.sufd = val
		case "pere":
            m.pere = modeles[val]
            m.herite()
        }
	}
	// il faut ajouter le dernier modèle lu
    m.heritedes()
	m.ajsuffd()
	m.ajsuff()
	modeles[m.id] = m
}
