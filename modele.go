package gocol

import (
	"fmt"
	"strings"
)

type Modele struct {
	id    string
	desm  map[int][]*Des
	lgenR []*Genrad
	pere  *Modele
	abs   []int
	suf   []string
	pos   string
	sufd  string
}

func (m Modele) doc() string {
	var mid string
	if m.pere != nil {
		mid = m.pere.id
	} else {
		mid = "nil"
	}
	return fmt.Sprintf("modele %s père %s, pos %s sufd %s, %d Genr",
		m.id, mid, m.pos, m.sufd, len(m.lgenR))
}

// vrai si le modeèle m a des désinences de numéro n
func (m Modele) habetD(nr, nd int) bool {
	for key, d := range m.desm {
		if key == nr && d[0].morpho == nd {
			return true
		}
	}
	return false
}

// habetR(gnr Genrad) bool
// vrai si le modèle m a déjà le générateur de radical gnr
func (m Modele) habetR(gnr *Genrad) bool {
	for _, genrad := range m.lgenR {
		if genrad.num == gnr.num {
			return true
		}
	}
	return false
}

// estabs(des *Des) bool
// vrai si le modèle m refuse d'hériter de la déninence
// de morpho n° m, parce ce n° figure dans ses abs
func (m Modele) estabs(des *Des) bool {
	for _, i := range m.abs {
		if i == des.morpho {
			return true
		}
	}
	return false
}

func (m *Modele) herite() {
	if m.pere == nil {
		return
	}
	// héritage du pos
	if m.pos == "" {
		m.pos = m.pere.pos
	}
	/*  héritage des générateurs de radicaux */
	for _, genr := range m.pere.lgenR {
		if !m.habetR(genr) {
			m.lgenR = append(m.lgenR, genr)
		}
	}
	// héritage des désinences
	for key, value := range m.pere.desm {
		for _, d := range value {
			if !m.estabs(d) && !m.habetD(d.nr, d.morpho) {
				nd := d.clone()
				nd.modele = m
				m.desm[key] = append(m.desm[key], nd)
			}
		}
	}
}

func (m *Modele) ajsuffd() {
	if m.sufd == "" {
		return
	}
	for _, ld := range m.desm {
		for _, d := range ld {
			d.grq = d.grq + m.sufd
			d.gr = atone(d.grq)
		}
	}
}

func (m *Modele) ajsuff() {
	for _, lsuf := range m.suf {
		ecl := strings.Split(lsuf, ":")
		li := listei(ecl[0])
		suff := ecl[1]
		for _, i := range li {
			for _, ld := range m.desm {
				for _, d := range ld {
					{
						if d.morpho == i {
							nd := d.clone()
							nd.grq = d.grq + suff
							nd.gr = atone(nd.grq)
							m.desm[d.nr] = append(m.desm[d.nr], nd)
						}
					}
				}
			}
		}
	}
}

func (m *Modele) ldesmorph(morph int) (ld []*Des) {
	for _, dd := range m.desm {
		for _, d := range dd {
			if d.morpho == morph {
				ld = append(ld, d)
			}
		}
	}
	return ld
}

var modeles = make(map[string]*Modele)
var vardes = make(map[string][]string)

func lismodeles(nf string) {
	//ll := Lignes(path + "modeles.la")
	ll := Lignes(nf)
	var m *Modele
	for _, l := range ll {
		if l == "" {
			continue
		}
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
			if m != nil {
				m.herite()
				m.ajsuffd()
				m.ajsuff()
				modeles[m.id] = m
			}
			m = new(Modele)
			m.pere = nil
			m.id = val
			m.desm = make(map[int][]*Des)
		case "R":
			num := Strtoint(ecl[1])
			if ecl[2] == "K" {
				m.lgenR = append(m.lgenR, &Genrad{num, 0, ""})
			} else {
				lp := strings.Split(ecl[2], ",")
				oter := Strtoint(lp[0])
				ajout := lp[1]
				if ajout == "0" {
					ajout = ""
				}
				m.lgenR = append(m.lgenR, &Genrad{num, oter, ajout})
			}
		case "abs":
			m.abs = listei(ecl[1])
		case "des", "des+":
			li := listei(ecl[1])
			// cas des variables
			nr := Strtoint(ecl[2])
			ld := ecl[3]
			var dd []string = strings.Split(ld, ";")
			var ddd []string
			for _, sdes := range dd {
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
					if sdes == "-" {
						sdes = ""
					}
					ddd = append(ddd, sdes)
				}
			}
			maxd := len(ddd)
			var nd *Des
			for ides, ili := range li {
				if ides < maxd {
					sld := ddd[ides]
					ecld := strings.Split(sld, ",")
					for _, cld := range ecld {
						nd = creeDes(cld, m, ili, nr)
						m.desm[nd.nr] = append(m.desm[nd.nr], nd)
					}
				} else {
					sld := ddd[maxd-1]
					ecld := strings.Split(sld, ",")
					for _, cld := range ecld {
						nd = creeDes(cld, m, ili, nr)
						m.desm[nd.nr] = append(m.desm[nd.nr], nd)
					}
				}
			}
			// si les désinences sont des+, le modèle doit
			// hériter des désinences de son père de même
			// morpho
			// TODO les des+ peuvent utiliser les $listes
			if cle == "des+" && m.pere != nil {
				li := listei(ecl[1])
				for _, nmorph := range li {
					// XXX non, les des de morph nmorph
					ldesp := m.pere.ldesmorph(nmorph)
					for _, desp := range ldesp {
						nd = desp.clone()
						nd.modele = m
						m.desm[nd.nr] = append(m.desm[nd.nr], nd)
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
		}
	}
	// il faut ajouter le dernier modèle lu
	m.herite()
	m.ajsuffd()
	m.ajsuff()
	modeles[m.id] = m
}
