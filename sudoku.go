package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------
// Constantes et Variables Globales
// -----------------------------------------------------------------------------

// FULL représente la séquence complète des chiffres de 1 à 9.
// Utilisée pour vérifier si une ligne, colonne ou bloc est correctement rempli.
var FULL = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

// DEBUG permet d'activer ou désactiver les messages de débogage.
const DEBUG = false

// Variables globales pour le suivi des itérations et des solutions trouvées.
var (
	nbIter      int
	nbSolutions int
)

// -----------------------------------------------------------------------------
// Structure Grille et Méthodes Associées
// -----------------------------------------------------------------------------

// Grille représente une grille de Sudoku avec 81 cases.
type Grille struct {
	Item []int // Slice contenant les valeurs des cases de la grille (0 pour vide)
}

// NewGrille crée une nouvelle instance de Grille en copiant la liste d'éléments fournie.
func NewGrille(itemList []int) *Grille {
	item := make([]int, len(itemList))
	copy(item, itemList)
	return &Grille{Item: item}
}

// String implémente l'interface Stringer pour afficher la grille de manière lisible.
func (g *Grille) String() string {
	return g.JoliPrint()
}

// JoliPrintBrut affiche la grille sous forme brute 9x9 sans séparateurs.
func (g *Grille) JoliPrintBrut() {
	for i := 0; i < 9; i++ {
		ligne := ""
		for j := 0; j < 9; j++ {
			ligne += fmt.Sprintf(" %d", g.Item[i*9+j])
		}
		fmt.Println(ligne)
	}
}

// JoliPrint affiche la grille avec des séparateurs pour une meilleure lisibilité.
func (g *Grille) JoliPrint() string {
	disp := ""
	for i := 0; i < 9; i++ {
		// Ajout des lignes séparatrices tous les 3 blocs
		if i%3 == 0 {
			disp += "------------+-----------+------------\n"
		}
		disp += "! "
		for j := 0; j < 9; j++ {
			disp += fmt.Sprintf(" %d ", g.Item[i*9+j])
			// Ajout des séparateurs verticaux tous les 3 blocs
			if j%3 == 2 {
				disp += " ! "
			}
		}
		disp += "\n"
	}
	disp += "------------+-----------+------------\n"
	return disp
}

// Ligne renvoie le contenu de la ligne i sous forme de slice d'entiers.
func (g *Grille) Ligne(i int) []int {
	return g.Item[i*9 : i*9+9]
}

// Colonne renvoie le contenu de la colonne j sous forme de slice d'entiers.
func (g *Grille) Colonne(j int) []int {
	col := make([]int, 0, 9)
	for i := 0; i < 9; i++ {
		col = append(col, g.Item[i*9+j])
	}
	return col
}

// BlocXY renvoie le contenu du bloc entourant la position (i, j).
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

// BlocIndex renvoie le contenu du bloc entourant l'index donné.
func (g *Grille) BlocIndex(index int) []int {
	i, j := g.Coord(index)
	return g.BlocXY(i, j)
}

// Coord renvoie les coordonnées (ligne, colonne) d'une position donnée dans la grille.
func (g *Grille) Coord(pos int) (int, int) {
	return pos / 9, pos % 9
}

// ChercheOrdre détermine l'ordre de recherche des cases restantes.
// Les cases sont classées par nombre croissant de possibilités pour optimiser la recherche.
func (g *Grille) ChercheOrdre() []int {
	type cellulePossibilite struct {
		cellule int // Index de la cellule
		poss    int // Nombre de possibilités pour cette cellule
	}

	// Slice pour stocker les cellules et leur nombre de possibilités
	possibiliteGrille := []cellulePossibilite{}

	// Parcours de toutes les cellules de la grille
	for pos := 0; pos < 81; pos++ {
		if g.Item[pos] == 0 { // Si la cellule est vide
			nbPossibilites := 0
			// Comptage des possibilités pour chaque chiffre de 1 à 9
			for num := 1; num <= 9; num++ {
				if g.EstPossible(num, pos) {
					nbPossibilites++
				}
			}
			// Ajout de la cellule et de son nombre de possibilités
			possibiliteGrille = append(possibiliteGrille, cellulePossibilite{cellule: pos, poss: nbPossibilites})
		}
	}

	// Tri des cellules par nombre de possibilités croissant
	sort.Slice(possibiliteGrille, func(i, j int) bool {
		return possibiliteGrille[i].poss < possibiliteGrille[j].poss
	})

	// Extraction de l'ordre des cellules triées
	ordreCellules := []int{}
	for _, cp := range possibiliteGrille {
		ordreCellules = append(ordreCellules, cp.cellule)
	}

	return ordreCellules
}

// EstResolu vérifie si la grille est complètement résolue et valide.
// Retourne true si toutes les lignes, colonnes et blocs contiennent les chiffres de 1 à 9 sans répétition.
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

	// Vérification des blocs 3x3
	for c := 0; c < 9; c++ {
		car := make([]int, len(g.BlocIndex(c)))
		copy(car, g.BlocIndex(c))
		sort.Ints(car)
		if !equalSlices(car, FULL) {
			return false
		}
	}

	// Si toutes les vérifications passent, la grille est résolue
	nbSolutions++ // Incrément du nombre de solutions trouvées
	return true
}

// EstPossible vérifie si un numéro peut être placé à une position donnée.
// Il doit être unique dans sa ligne, colonne et bloc.
func (g *Grille) EstPossible(num, pos int) bool {
	lig, col := g.Coord(pos)
	// Vérifie l'absence du numéro dans la ligne, la colonne et le bloc
	if !contains(g.Ligne(lig), num) &&
		!contains(g.Colonne(col), num) &&
		!contains(g.BlocXY(lig, col), num) {
		if DEBUG {
			fmt.Printf("On peut mettre %d dans la grille à la position (%d, %d).\n", num, lig, col)
		}
		return true
	}
	return false
}

// -----------------------------------------------------------------------------
// Fonctions Auxiliaires
// -----------------------------------------------------------------------------

// equalSlices vérifie l'égalité de deux slices d'entiers.
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

// contains vérifie si une slice contient un élément spécifique.
func contains(slice []int, elem int) bool {
	for _, val := range slice {
		if val == elem {
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------------------
// Lecture de la Grille depuis un Fichier
// -----------------------------------------------------------------------------

// LectureFichier lit une grille de Sudoku à partir d'un fichier texte.
// Le fichier doit contenir 9 lignes de 9 caractères chacune.
// Les chiffres représentent les cases remplies et tout autre caractère (sauf espace) représente une case vide.
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
		ligne = strings.ReplaceAll(ligne, " ", "") // Suppression des espaces
		if len(ligne) == 9 {                       // Chaque ligne doit contenir exactement 9 caractères
			for _, c := range ligne {
				num, err := strconv.Atoi(string(c))
				if err != nil {
					// Si le caractère n'est pas un chiffre, considérer la case comme vide (0)
					grille = append(grille, 0)
				} else {
					grille = append(grille, num)
				}
			}
		} else if len(ligne) == 0 {
			// Ignorer les lignes vides
			continue
		} else {
			return nil, fmt.Errorf("Mauvais format de ligne (longueur attendue : 9, lue : %d)", len(ligne))
		}
	}

	// Vérification du nombre total de cases
	if len(grille) != 81 {
		return nil, fmt.Errorf("Mauvais nombre de cases (attendu : 81, lu : %d)", len(grille))
	}

	return NewGrille(grille), nil
}

// -----------------------------------------------------------------------------
// Fonction de Résolution de la Grille (Algorithme de Backtracking)
// -----------------------------------------------------------------------------

// Cherche effectue la recherche de solution de la grille de Sudoku de manière récursive.
// Utilise l'algorithme de backtracking pour explorer les possibilités.
func Cherche(grille *Grille) bool {
	nbIter++ // Incrément du compteur d'itérations

	// Affichage périodique du nombre d'essais pour suivre la progression
	if nbIter%5000 == 0 {
		fmt.Printf("%d essais\n", nbIter)
	}

	// Affichage de débogage si activé
	if DEBUG {
		fmt.Printf("(%d)\n", nbIter)
		fmt.Println(grille.JoliPrint())
		fmt.Println()
	}

	// Vérification si la grille est résolue
	if grille.EstResolu() {
		// Affichage de la solution trouvée
		fmt.Println(strings.Repeat("=", 40))
		fmt.Printf("SOLUTION #%d TROUVEE (%d itérations) :\n", nbSolutions, nbIter)
		fmt.Println(strings.Repeat("=", 40))
		fmt.Println()
		fmt.Println(grille)
		return true // Retourne true pour indiquer qu'une solution a été trouvée
	} else {
		// Détermination de l'ordre de recherche des cases restantes
		ordre := grille.ChercheOrdre()

		// Affichage de l'ordre de recherche si en mode débogage
		if DEBUG {
			fmt.Println(ordre)
		}

		// Parcours des cellules dans l'ordre déterminé
		for _, ind := range ordre {
			// Essai de placer chaque chiffre de 1 à 9 dans la cellule
			for num := 1; num <= 9; num++ {
				if grille.Item[ind] == 0 { // Si la cellule est vide
					if grille.EstPossible(num, ind) { // Vérifie si le chiffre peut être placé
						// Création d'une nouvelle grille avec le chiffre placé
						newGrille := NewGrille(grille.Item)
						if DEBUG {
							fmt.Printf("grille    %p\n", grille)
							fmt.Printf("new_grille %p\n", newGrille)
						}
						newGrille.Item[ind] = num

						// Appel récursif de la fonction Cherche avec la nouvelle grille
						Cherche(newGrille)
					}
				} else {
					// Si la cellule n'est pas vide, passer à la suivante
					break
				}
			}
			// Après avoir tenté de remplir une cellule, ne pas continuer avec les autres
			// Cela permet de revenir en arrière si nécessaire (backtracking)
			break
		}

		// Affichage de fin de recherches infructueuses si en mode débogage
		if DEBUG {
			fmt.Println("Fin des recherches (infructueuses)")
		}
	}

	return false // Retourne false si aucune solution n'a été trouvée à ce niveau
}

// -----------------------------------------------------------------------------
// Fonction Principale
// -----------------------------------------------------------------------------

func main() {
	// Nom de fichier par défaut
	DEFAULT_FILE_NAME := "1.txt"

	// Gestion des arguments de ligne de commande
	nbArg := len(os.Args) - 1
	var fileName string

	if nbArg == 1 {
		// Si un argument est fourni, l'utiliser comme nom de fichier
		fileName = os.Args[1]
	} else {
		// Sinon, demander à l'utilisateur de saisir le nom du fichier
		fmt.Printf("Nom du fichier (grille), sans .txt [%s]: ", DEFAULT_FILE_NAME)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Erreur de lecture:", err)
			return
		}
		input = strings.TrimSpace(input) // Suppression des espaces et des retours chariot
		if len(input) == 0 {
			// Utiliser le nom par défaut si aucune saisie
			fileName = DEFAULT_FILE_NAME
		} else {
			// Ajouter l'extension .txt si nécessaire
			if len(input) < 4 || !strings.HasSuffix(input, ".txt") {
				fileName = input + ".txt"
			} else {
				fileName = input
			}
		}
	}

	// Lecture de la grille depuis le fichier
	grille, err := LectureFichier(fileName)
	if err != nil {
		// Affichage de l'erreur et terminaison si le fichier est mal formaté
		fmt.Println(err)
		return
	}
	fmt.Println("\nGrille initiale")
	fmt.Println(grille)
	fmt.Println()

	// Vérification si la grille est déjà résolue
	if grille.EstResolu() {
		fmt.Println("La grille en entrée est déjà terminée ! Il n'y a rien à faire !")
		return
	}

	// Démarrage de la recherche de solution
	fmt.Println("On démarre la recherche...\n")
	t0 := time.Now() // Début du chronomètre
	Cherche(grille)
	t1 := time.Now() // Fin du chronomètre
	elapsed := t1.Sub(t0).Seconds()
	fmt.Printf("Problème résolu en %f secondes.\n\n", elapsed)
}
