package gocol

import (
	"fmt"
	"strconv"
	"strings"
)

var lenm int
var Morphos = make(map[int]string)

func lismorphos(path string) {
	ll := Lignes(path+"morphos.fr")
	for i := 0; i < len(ll); i++ {
		ecl := strings.Split(ll[i], ":")
		if len(ecl) > 1 {
			n, _ := strconv.Atoi(ecl[0])
			Morphos[n] = ecl[1]
		}
	}
	lenm = len(Morphos)
}

func morpho(i int) string {
	if i > lenm {
		return fmt.Sprintf("err. il n'y a que %d morphos\n", lenm)
	}
	return Morphos[i]
}
