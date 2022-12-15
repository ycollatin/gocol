package gocol

import "strings"

type re struct {
	g string
	d string
}

var lexp []re

// vrai si l est dans la liste ll
func estDans(ll []string, l string) bool {
	for _, item := range ll {
		if item == l {
			return true
		}
	}
	return false
}

// vrai si s se termine en f
func finEn(s string, f string) bool {
	return strings.TrimSuffix(s, f) != s
}

// vrai si s commence en d
func debEn(s string, d string) bool {
	return strings.TrimPrefix(s, d) != s
}

// supprime les n derniers caractères de s 
func chop(s string, n int) string {
	return s[:len(s)-n]
}

// retourne ante modifié par re
func vars(ante string, exp re) string {
	g := exp.g
	d := exp.d
	var recipr bool = true
	if finEn(g, "$") {
		g = chop(g, 1)
		if !finEn(ante, g) {
			return ante
		}
        return chop(ante, len(g))
	} else if finEn(g, ">") {
		g = chop(g, 1)
		recipr = false
	}
	if debEn(g, "[^v]") {
		g = g[4:]
		if debEn(ante, "v"+g) {
			return ante
		}
	}
	if finEn(g, "[^n]") {
		g = chop(g, 4)
		if finEn(ante, "n"+g) {
			return ante
		}
	}
	if contient(ante, g) {
		ecl := strings.Split(ante, g)
		return ecl[0] + d + ecl[1]
	}
	if recipr && contient(ante, d) {
		ecl := strings.Split(ante, d)
		return ecl[0] + g + ecl[1]
	}
	return ante
}

func varsL(ante []string, exp re) (post []string) {
	for _, f := range ante {
		v := vars(f, exp)
		if !estDans(post, v) {
			post = append(post, v)
		}
	}
	return
}

// retourne toutes les variantes possibles de f
func VarsF(f string) (post []string) {
	post = append(post, f)
	for _, r := range lexp {
		lf := varsL(post, r)
		for _, v := range lf {
			if !estDans(post, v) {
				post = append(post, v)
			}
		}
	}
	ftl := Deramise(strings.ToLower(f))
	if ftl != f {
		post = append(post, VarsF(ftl)...)
	}
	return
}

func lisExp(nf string) {
	lreg := Lignes(nf)
	for _, l := range lreg {
		ecl := strings.Split(l, ":")
		lexp = append(lexp, re{g: ecl[0], d: ecl[1]})
	}
}
