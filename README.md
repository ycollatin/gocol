# gocoll

_lemmatisation de textes latins_

## Installation
`$ go get github.com/ycollatin/gocoll

## Utilisation

~~~~
import (
	//...
	"github.com/ycollatin/gocoll"
	//...
)
~~~~

Il faut d'abord lire les données en leur fournissant
un chemin.

`data("le/chemin/des/données")`

## Principales fonction : 

- `godoc.Lemmatise(s string) (lsr Res, echec bool)`    
lemmatisation de la forme "uocem"
- Pour convertir ces résultats en chaîne lisible :    
`godoc.Restostring(r Res)`

La lemmatisation ѕe fait donc en deux temps. Par exemple :

~~~~
results, ech := godoc.Lemmasise("uocem")
if !ech {
	for r := range results {
	fmt.Println(godoc.Restostring(r)
	} else {
	fmt.Println("echec")
}
~~~~

- Pour une API complète, utiliser godoc.
