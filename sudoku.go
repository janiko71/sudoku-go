package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"sort"
)

// Constantes
var (
	FULL  = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	DEBUG = false
)

// Variables globales
var (
	nbIter      int
	nbSolutions int
)

// Grille représente une grille de Sudoku
type Grille struct {
	Item []int
}

// NewGrille crée une nouvelle grille à partir d'une liste d'entiers
func NewGrille(itemList []int) *Grille {
	item := make([]int, len(itemList))
	copy(item, itemList)
	return &Grille{Item: item}
}

// String implémente l'interface Stringer pour afficher la grille
func (g *Grille) String() string {
	return g.JoliPrint()
}

// JoliPrintBrut affiche la grille sous forme brute 9x9
func (g *Grille) JoliPrintBrut() {
	for i := 0; i < 9; i++ {
		ligne := ""
		for j := 0; j < 9; j++ {
			ligne += fmt.Sprintf(" %d", g.Item[i*9+j])
		}
		fmt.Println(ligne)
	}
}

// JoliPrint affiche la grille avec des séparateurs améliorés
func (g *Grille) JoliPrint() string {
	disp := ""
	for i := 0; i < 9; i++ {
		if i%3 == 0 {
			disp += "------------+-----------+------------\n"
		}
		disp += "! "
		for j := 0; j < 9; j++ {
			disp += fmt.Sprintf(" %d ", g.Item[i*9+j])
			if j%3 == 2 {
				disp += " ! "
			}
		}
		disp += "\n"
	}
	disp += "------------+-----------+------------\n"
	return disp
}

// Ligne renvoie le contenu de la ligne i
func (g *Grille) Ligne(i int) []int {
	return g.Item[i*9 : i*9+9]
}

// Colonne renvoie le contenu de la colonne j
func (g *Grille) Colonne(j int) []int {
	col := make([]int, 0, 9)
	for i := 0; i < 9; i++ {
		col = append(col, g.Item[i*9+j])
	}
	return col
}

// BlocXY renvoie le contenu du bloc entourant la position (i, j)
func (g *Grille) BlocXY(i, j int) []int {
	bloc := []int{}
	startRow := (i / 3) * 3
	startCol := (j / 3) * 3
	for r := startRow; r < startRow+3; r++ {
		for c := startCol; c < startCol+3; c++ {
			bloc = append(bloc, g.Item[r*9+c])
		}
	}
	return bloc
}

// BlocIndex renvoie le contenu du bloc entourant l'index donné
func (g *Grille) BlocIndex(index int) []int {
	i, j := g.Coord(index)
	return g.BlocXY(i, j)
}

// Coord renvoie les coordonnées (x, y) d'une position donnée
func (g *Grille) Coord(pos int) (int, int) {
	return pos / 9, pos % 9
}

// ChercheOrdre détermine l'ordre de recherche des cases restantes
func (g *Grille) ChercheOrdre() []int {
	type cellulePossibilite struct {
		cellule int
		poss    int
	}

	possibiliteGrille := []cellulePossibilite{}

	for pos := 0; pos < 81; pos++ {
		if g.Item[pos] == 0 {
			nbPossibilites := 0
			for num := 1; num <= 9; num++ {
				if g.EstPossible(num, pos) {
					nbPossibilites++
				}
			}
			possibiliteGrille = append(possibiliteGrille, cellulePossibilite{cellule: pos, poss: nbPossibilites})
		}
	}

	// Trier par nombre de possibilités croissant
	sort.Slice(possibiliteGrille, func(i, j int) bool {
		return possibiliteGrille[i].poss < possibiliteGrille[j].poss
	})

	ordreCellules := []int{}
	for _, cp := range possibiliteGrille {
		ordreCellules = append(ordreCellules, cp.cellule)
	}

	return ordreCellules
}

// EstResolu vérifie si la grille est résolue
func (g *Grille) EstResolu() bool {
	// Vérification des lignes
	for i := 0; i < 9; i++ {
		lig := make([]int, len(g.Ligne(i)))
		copy(lig, g.Ligne(i))
		sort.Ints(lig)
		if !equalSlices(lig, FULL) {
			return false
		}
	}

	// Vérification des colonnes
	for j := 0; j < 9; j++ {
		col := make([]int, len(g.Colonne(j)))
		copy(col, g.Colonne(j))
		sort.Ints(col)
		if !equalSlices(col, FULL) {
			return false
		}
	}

	// Vérification des blocs
	for c := 0; c < 9; c++ {
		car := make([]int, len(g.BlocIndex(c)))
		copy(car, g.BlocIndex(c))
		sort.Ints(car)
		if !equalSlices(car, FULL) {
			return false
		}
	}

	// Toutes les conditions sont remplies
	nbSolutions++
	return true
}

// EstPossible vérifie si un numéro peut être placé à une position donnée
func (g *Grille) EstPossible(num, pos int) bool {
	lig, col := g.Coord(pos)
	if !contains(g.Ligne(lig), num) && !contains(g.Colonne(col), num) && !contains(g.BlocXY(lig, col), num) {
		if DEBUG {
			fmt.Printf("On peut mettre %d dans la grille à la position (%d, %d).\n", num, lig, col)
		}
		return true
	}
	return false
}

// Fonction auxiliaire pour vérifier l'égalité de deux slices d'entiers
func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, val := range a {
		if val != b[i] {
			return false
		}
	}
	return true
}

// Fonction auxiliaire pour vérifier si une slice contient un élément
func contains(slice []int, elem int) bool {
	for _, val := range slice {
		if val == elem {
			return true
		}
	}
	return false
}

// LectureFichier lit une grille de Sudoku à partir d'un fichier
func LectureFichier(nomFichier string) (*Grille, error) {
	file, err := os.Open(nomFichier)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de l'ouverture du fichier: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	grille := []int{}

	for scanner.Scan() {
		ligne := strings.TrimSpace(scanner.Text())
		ligne = strings.ReplaceAll(ligne, " ", "")
		if len(ligne) == 9 {
			for _, c := range ligne {
				num, err := strconv.Atoi(string(c))
				if err != nil {
					grille = append(grille, 0)
				} else {
					grille = append(grille, num)
				}
			}
		} else if len(ligne) == 0 {
			continue
		} else {
			return nil, fmt.Errorf("Mauvais format de ligne (longueur attendue : 9, lue : %d)", len(ligne))
		}
	}

	if len(grille) != 81 {
		return nil, fmt.Errorf("Mauvais nombre de cases (attendu : 81, lu : %d)", len(grille))
	}

	return NewGrille(grille), nil
}

// Cherche effectue la recherche de solution de la grille
func Cherche(grille *Grille) bool {
	nbIter++
	if nbIter%5000 == 0 {
		fmt.Printf("%d essais\n", nbIter)
	}
	if DEBUG {
		fmt.Printf("(%d)\n", nbIter)
		fmt.Println(grille.JoliPrint())
		fmt.Println()
	}

	if grille.EstResolu() {
		fmt.Println(strings.Repeat("=", 40))
		fmt.Printf("SOLUTION #%d TROUVEE (%d itérations) :\n", nbSolutions, nbIter)
		fmt.Println(strings.Repeat("=", 40))
		fmt.Println()
		fmt.Println(grille)
		return true
	} else {
		ordre := grille.ChercheOrdre()

		if DEBUG {
			fmt.Println(ordre)
		}

		for _, ind := range ordre {
			for num := 1; num <= 9; num++ {
				if grille.Item[ind] == 0 {
					if grille.EstPossible(num, ind) {
						newGrille := NewGrille(grille.Item)
						if DEBUG {
							fmt.Printf("grille    %p\n", grille)
							fmt.Printf("new_grille %p\n", newGrille)
						}
						newGrille.Item[ind] = num
						Cherche(newGrille)
					}
				} else {
					break
				}
			}
			// Utilité du break ?
			break
		}

		if DEBUG {
			fmt.Println("Fin des recherches (infructueuses)")
		}
	}

	return false
}

func main() {
	DEFAULT_FILE_NAME := "1.txt"

	// Gestion des arguments de ligne de commande
	nbArg := len(os.Args) - 1
	var fileName string

	if nbArg == 1 {
		fileName = os.Args[1]
	} else {
		fmt.Printf("Nom du fichier (grille), sans .txt [%s]: ", DEFAULT_FILE_NAME)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Erreur de lecture:", err)
			return
		}
		input = strings.TrimSpace(input)
		if len(input) == 0 {
			fileName = DEFAULT_FILE_NAME
		} else {
			if len(input) < 4 || !strings.HasSuffix(input, ".txt") {
				fileName = input + ".txt"
			} else {
				fileName = input
			}
		}
	}

	// Lecture de la grille
	grille, err := LectureFichier(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("\nGrille initiale")
	fmt.Println(grille)
	fmt.Println()

	if grille.EstResolu() {
		fmt.Println("La grille en entrée est déjà terminée ! Il n'y a rien à faire !")
		return
	}

	// Démarrage de la recherche
	fmt.Println("On démarre la recherche...\n")
	t0 := time.Now()
	Cherche(grille)
	t1 := time.Now()
	elapsed := t1.Sub(t0).Seconds()
	fmt.Printf("Problème résolu en %f secondes.\n\n", elapsed)
}
