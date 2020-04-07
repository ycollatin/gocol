package gocol

import (
	"fmt"
	"strings"
)

type Genrad struct {
	num   int
	oter  int
	ajout string
}

func (gr Genrad) doc() (d string) {
	return fmt.Sprintf("Genrad %d, oter %d, ajouter %s", gr.num, gr.oter, gr.ajout)
}

//  rad(k string)  permet de générer un radical à partir
//  de la forme canonique k
func (genr Genrad) rad(k string) (r string) {
	// éliminer le rune sans largeur
	k = strings.Replace(k, string(0x306), "", -1)
	rs := []rune(k)
	l := len(rs)
	rs = rs[0 : l-genr.oter]
	r = string(rs) + genr.ajout
	return r
}
