package main

import (
	"bufio"
	"fmt"
	"github.com/ycollatin/gocol"
	"os"
	"path"
	"strings"
)

func main() {
	dir, _ := os.Executable()
	ch := path.Dir(dir)
	chData := ch + "/data/"
	// lecture des données Collatinus
	gocol.Data(chData)
	var f string
	lecteur := bufio.NewReader(os.Stdin)
	for f != "x" {
		fmt.Print("> ")
		f, _ = lecteur.ReadString('\n')
		f = f[:len(f)-1]
		if f == "x" {
			break
		}
		ret, _ := gocol.Lemmatise(f)
		for _, r := range ret {
			fmt.Printf("   %s, %s, %s : %s\n",
			strings.Join(r.Lem.Grq, " "),
			r.Lem.Pos,
			r.Lem.Indmorph,
			r.Lem.Traduction)
			for _, m := range r.Morphos {
				fmt.Println("      ", m)
			}
		}
	}
}
