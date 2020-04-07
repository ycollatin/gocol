// lemmatisation et analyse morpho de textes latins
package gocol

/*
	XXX - bogues

	TODO
	- écrire un lisdata() pour que les données du module soient toutes lues.
	  Seuls les lemmes sont lus. Il faudrait ajouter
	   . irréguliers
	   . vargraph
	   . modeles
	- moteur.go : il ne devrait pas y avoir de désinence nil dans la
	  boucle lemmatise()
	- paramètres supplémentaires possibles :
        -w sortie html
		-t traduction
		-l tri alpha des lemmes
		-m tri alpha des mots
		-e rejet des échecs en fin
		-c calcul des stats (fréquence des formes, des lemmes, des morphos, des pos)
		-o <fichier de sortie>
		-s serveur sur le port indiqué, ou un port par défaut
	- modeles.la : nombreuses désinences héritées redéfinies à l'identique
*/


import (
	"C"
)

var (
	dat bool
	module string
	modules []string
)

const version = "Alpha"

// lecture des données et affichage des effectifs
func Data(path string) {
	if dat {
		return
	}
	if path[len(path)-1] != '/' {
		path = path+"/"
	}
	lismorphos(path)
	//fmt.Println(len(morphos), "morphos")
	lismodeles(path)
	//fmt.Println(len(modeles), "modèles")
	var nc string
	if len(module) > 0 {
		nc = path + module + "/"
		lisLemmes(nc + "lemmes.la")
		lisTraductions(nc + "lemmes.fr")
		//fmt.Println("module", module, len(lemmes), "lemmes")
		lisExp(nc + "vargraph.la")
	}
	lisExp(path+"vargraph.la")
	lisLemmes(path+"lemmes.la")
	lisTraductions(path+"lemmes.fr")
	//fmt.Println(len(lemmes), "lemmes")
	ajRadicaux()
	//fmt.Println(len(radicaux), "radicaux")
	if len(module) > 0 {
		lisIrregs(path + module + "/irregs.la")
	}
	lisIrregs(path+"irregs.la")
	//fmt.Println(len(irregs), "irréguliers")
	//fmt.Println(len(lexp), "variantes graphiques\n")
	dat = true
}

// lemmatise tous les mots du texte s
func LemmatiseT(t string) string {
	lm := Mots(t)
	var sortie string
	for _, m := range lm {
		an, echec := Lemmatise(m)
		if echec {
			an, echec = Lemmatise(Majminmaj(m))
		}
		if echec {
				sortie = ("\n" + m + " échec")
		} else {
			sortie += Restostring(an) + "\n"
		}
	}
	return sortie
}

func main() {}
