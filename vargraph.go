package gocol

import "strings"

type re struct {
	g string
	d string
}

var lexp []re

func estDans(ll []string, l string) bool {
	for _, item := range ll {
		if item == l {
			return true
		}
	}
	return false
}

func finEn(s string, f string) bool {
	return strings.TrimSuffix(s, f) != s
}

func debEn(s string, d string) bool {
	return strings.TrimPrefix(s, d) != s
}

func chop(s string, n int) string {
	return s[:len(s)-n]
}

func vars(ante string, exp re) string {
	g := exp.g
	d := exp.d
	var recipr bool = true
	if finEn(g, "$") {
		g = chop(g, 1)
		if !finEn(ante, g) {
			return ante
		}
		recipr = false
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

func varsF(f string) (post []string) {
	post = append(post, f)
	var lf []string
	for _, r := range lexp {
		lf = varsL(post, r)
		for _, v := range lf {
			if !estDans(post, v) {
				post = append(post, v)
			}
		}
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
