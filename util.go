package gocol

/*
Le découpage en runes peut être utile :
https://www.dotnetperls.com/substring-go
*/

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var mapIVX = make(map[rune]int)

func aRomano(f string) string {
	if f == "" {
		return "-1"
	}

	if len(mapIVX) == 0 {
		mapIVX = map[rune]int{
			'I': 1,
			'V': 5,
			'X': 10,
			'L': 50,
			'C': 100,
			'D': 500,
			'M': 1000,
		}
	}
	f = strings.ToUpper(f)
	var conv_c, res int
	conv_s := mapIVX[rune(f[0])]
	//for i, _ := range f {
	for i := 0; i < len(f)-1; i++ {
		conv_c = conv_s
		conv_s = mapIVX[rune(f[i+1])]
		if conv_c < conv_s {
			res -= conv_c
		} else {
			res += conv_c
		}
	}
	res += conv_s
	return strconv.Itoa(res)
}

func atone(ch string) string {
	var trans = map[rune]string{
		'ă': "a",
		'ā': "a",
		'ĕ': "e",
		'ē': "e",
		'ĭ': "i",
		'ī': "i",
		'ŏ': "o",
		'ō': "o",
		'ŭ': "u",
		'ū': "u",
		'ụ': "u",
		'ў': "y",
		'ȳ': "y",
		'Ă': "A",
		'Ā': "A",
		'Ĕ': "E",
		'Ē': "E",
		'Ĭ': "I",
		'Ī': "I",
		'Ŏ': "O",
		'Ō': "O",
		'Ŭ': "U",
		'Ū': "U",
		'Ў': "Y",
		'Ȳ': "Y",
	}
	b := bytes.NewBufferString("")
	ch = strings.Replace(ch, string(0x306), "", -1)
	for _, c := range ch {
		if val, ok := trans[c]; ok {
			b.WriteString(val)
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func contient(tout string, part string) bool {
	return strings.Index(tout, part) > -1
}

func deramAtone(ch string) string {
	return Deramise(atone(ch))
}

func Deramise(ch string) string {
	nch := ch
	nch = strings.Replace(nch, "v", "u", -1)
	nch = strings.Replace(nch, "j", "i", -1)
	nch = strings.Replace(nch, "J", "I", -1)
	return nch
}

func estRomain(f string) bool {
	f = strings.ToUpper(f)
	lr := "IVXLCDM"
	for _, c := range f {
		if !contient(lr, string(c)) {
			return false
		}
	}
	if contient(f, "IL") || contient(f, "IVI") {
		return false
	}
	return true
}

// insère le genre du nom entre le cas et le nombre
func genreNom(m, g string) string {
	lm := strings.Split(m, " ")
	if len(lm) < 2 {
		return g
	}
	ls := []string{lm[0],g,lm[1]}
	return strings.Join(ls," ")
}

// donne sous forme de slice toutes
// les lignes du fichier nf. Les
// lignes commençant par '!' sont
// considérées comme commentaire
// et ignorées.
func Lignes(nf string) []string {
	f, err := os.Open(nf)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	lecteur := bufio.NewScanner(f)
	var ll []string
	for lecteur.Scan() {
		l := string(lecteur.Text())
		if (len(l) > 0) && (l[0] != '!') {
			ll = append(ll, l)
		}
	}
	return ll
}

func listei(s string) (li []int) {
	lvirg := strings.Split(s, ",")
	for _, virg := range lvirg {
		if strings.Contains(virg, "-") {
			tiret := strings.Split(virg, "-")
			deb := Strtoint(tiret[0])
			fin := Strtoint(tiret[1])
			for j := deb; j <= fin; j++ {
				li = append(li, j)
			}
		} else {
			li = append(li, Strtoint(virg))
		}
	}
	return
}

// inverse la casse de s.
func Majminmaj(s string) string {
	smin := strings.ToLower(s)
	if smin != s {
		return smin
	}
	return strings.Title(s)
}

func Mots(s string) []string {
	entremots := func(c rune) bool {
		return !unicode.IsLetter(c)
	}
	return strings.FieldsFunc(s, entremots)
}

func Strtoint(s string) (n int) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return
}
