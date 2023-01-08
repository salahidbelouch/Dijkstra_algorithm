package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"runtime"
)

//gérer les erreurs
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//type elementGraph, noius servira pour créer un slice qui represente chaque liason dans le graphe
type elementGraph struct {
	from   int
	to     int
	weight int
}

// sépare les données reçues en deux noeuds et un poids et les ajoute à un slice de elementGraph
func addToSlice(str string, slice []elementGraph) []elementGraph {

	line := strings.Split(str, " ") // retourne un tableau de strings line contenant les sous chaines séparées par un espace

	x, errr := strconv.Atoi(line[0]) // le noeud from
	y, errr := strconv.Atoi(line[1]) // le noeud to
	p, errr := strconv.Atoi(line[2]) // le poids

	check(errr)

	if p != 0 { //si le poids n'est pas nul

		slice = append(slice, elementGraph{x, y, p}) // on ajoute au slice from x to y et  weight p

	}

	return slice

}

// permet de representer un array en string ( plus simple pour faire entrer des chaines de string dans le channel)
func arrayToString(a []int, delim string) string {

	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

// determine si la liste L contient l'élément x
func contains(L []int, x int) bool {
	for i := 0; i < len(L); i++ {
		if L[i] == x {
			return true
		}
	}
	return false
}

// donne la distance (le poids) entre a et b dans un graph (slice)
func length(a int, b int, s []elementGraph) int {

	for i := 0; i < len(s); i++ { //modifier en FOR(while) pour moins de complexité

		if (s[i].from == a) && (s[i].to == b) { // si {a b p} existe dans le slice   return p

			return s[i].weight
		}
	}
	return 0
}

// génére une liste des noeuds [0 1 2 3 ... n ] pour n noeuds
//on aurait pu chercher le max dans slice.to ou/et slice.from et faire une boucle for jusqu'au max avec append dans l'array mais c'est plus propre comme ça
func list_noeuds(slice []elementGraph) []int {

	Node := []int{} //  contient tous les noeuds

	for i := 0; i < len(slice); i++ { //pour chaque arête

		if !contains(Node, slice[i].from) { // si le noeud n'est pas déja pris en compte l'ajouter dans Node
			Node = append(Node, slice[i].from)
		}
	}
	sort.Sort(sort.IntSlice(Node)) //trie

	return Node
}

// initialisation de deux arrays qu'on utilisera avec des grandes valeurs
func initialiser(n int) ([]int, []int) {

	prev := make([]int, n)
	distance := make([]int, n)
	for i := 0; i < n; i++ {

		distance[i] = 9998 // 9999 = oo : pour dire que la distance est max au debut (sert pour trouver le min)
		prev[i] = 9999     // 9998 = indéfinie : pour dire que les prédecesseurs ne sont pas encore définie

	}
	return distance, prev
}
func dijkstra(slice []elementGraph, liste []int, dist []int, pre []int, source int) string {

	var results string

	//instanciation
	Unvisited := make([]int, len(liste)) // noeuds non visité
	prev := make([]int, len(liste))      // prédecesseurs de chaque neouds pour arriver à la source
	distance := make([]int, len(liste))  // distance de la source à chaque noeud
	neighbors := make(map[int][]int)     
	ways := make(map[int][]int)          

	//pour eviter 2 boucles for à chaque appel de la fonction
	copy(Unvisited, liste)
	copy(distance, dist)
	copy(prev, pre)

	distance[source] = 0 // le poids (la distance) d'une source à elle même est nulle

	Q := []int{}

	for i := 0; i < len(slice); i++ {

		if slice[i].from != slice[i].to {

			if !contains(Q, slice[i].from) {

				neighbors[slice[i].from] = []int{slice[i].to}
				Q = append(Q, slice[i].from)

			} else {

				neighbors[slice[i].from] = append(neighbors[slice[i].from], slice[i].to)
			}
		}
	}

	
	for len(Unvisited) != 0 { // tant qu'il y a des noeuds non visités

		min := distance[Unvisited[0]] // Initialisation de min; min = min(distance[u]) avec u ∈ Unvisited
		u := Unvisited[0]
		j := 0

		for k := 0; k < len(Unvisited); k++ {

			if distance[Unvisited[k]] < min {

				min = distance[Unvisited[k]]
				u = Unvisited[k] // récupération du noeudu ∈ Unvisited, contenant la distance la plus courte
				j = k            // j = indice de u dans Unvisited

			}

		}
		Unvisited = append(Unvisited[:j], Unvisited[(j+1):]...) // retirer u de Unvisited

		
		for v := 0; v < len(neighbors[u]); v++ {

			var poids int = length(u, neighbors[u][v], slice)

			if contains(Unvisited, neighbors[u][v]) && (distance[u]+poids < distance[neighbors[u][v]]) {

				distance[neighbors[u][v]] = distance[u] + poids 
				prev[neighbors[u][v]] = u                       
			}
		}

		ways[u] = []int{} // le chemin le plus court de la source à un noeud u

		
		if u != source {
			v := prev[u]

			for v != source {
				ways[u] = append(ways[u], v)
				v = prev[v]
			}

		}

		// inverser une liste (ici ways[u]) car le chemin est generé à l'envers
		for i, j := 0, len(ways[u])-1; i < j; i, j = i+1, j-1 {
			ways[u][i], ways[u][j] = ways[u][j], ways[u][i]
		}
		ways[u] = append(ways[u], u)
	}

	for i := 0; i < len(liste); i++ {
		results = results + " from " + strconv.Itoa(source) + " to " + strconv.Itoa(i) + " [ " + arrayToString(ways[i], " ") + " ] " + strconv.Itoa(distance[i]) + "  " + "\n"
	}

	return results

}

// fonction qui va mettre dans le channel le resultat de chaque disjktra calculé pour chaque noeuds donné par jobs
func worker(wg sync.WaitGroup, slice []elementGraph, liste []int, dist []int, pre []int, jobs <-chan int, resu chan<- string) {
	wg.Add(1)
	for n := range jobs {
		resu <- dijkstra(slice, liste, dist, pre, n)
	}
	wg.Done()

	
}

// fonction qui permet de d'envoyer dans le channel jobs les neouds pour lesquels le worker doit le calculer et prend le resultat du channel resul pour l'ecrire dans le fichier
func calcule(conn net.Conn ,wg sync.WaitGroup ,graphe []elementGraph, liste []int, dist []int, pre []int, Nom_resultat string) {
	
	jobs := make(chan int, len(liste))
	resul := make(chan string, len(liste))
	file, err := os.OpenFile(Nom_resultat, os.O_CREATE|os.O_WRONLY, 0600)
	check(err)

	debut := time.Now()
	for i := 0; i < len(liste); i++ {
		jobs <- i //met tous les noeuds pour lesquels on doit calculer le dijsktra dans le channel jobs
	}
	for i := 0; i < runtime.NumCPU(); i++ {

		go worker(wg, graphe, liste, dist, pre, jobs, resul) //appel le worker dans chaque go routine pour calculer le dijkstra à partir des noeuds dans le channel jobs
	}
	go func() {
		wg.Add(1)
		for i := 0; i < len(liste); i++ {
			_, err = file.WriteString(<-resul) //écrit ce qui est mis dans le channel resul par le wroker dans le fichier
			check(err)
			if i == len(liste)-1 {
			   wg.Done()
			   fin := time.Now()
			   temps := fin.Sub(debut) // donne le temps d'éxécution de la fonction calcule
			   resultat_temps := temps.String()
			   fmt.Println("Time : " + resultat_temps)
               io.WriteString(conn, resultat_temps+"\n")
			   fmt.Println("DONE")
			}
		}


	}()

	
	close(jobs) // fermeture du channel

}

//fonction quin gère chaque nouvelle connexion
func handleConnection(conn net.Conn) {
	var wg sync.WaitGroup
	defer conn.Close()                //fermeture de la connexion à la fin du prog
	var Nom_resultat string           // on va la récuperer du client
	Graphe := []elementGraph{}        //intialise le graphe pour le remplire avec les arêtes
	scanner := bufio.NewScanner(conn) //scanner écoute dans la connexion établie
	i := 0                            //pour traiter le differente données recues

	for scanner.Scan() { // tant qu'il y a des données qui arrive

		message := scanner.Text()
		if i == 0 { // premier tour pour prendre le nom du fichier resultat
			i++ // s'excécute une seule fois
			Nom_resultat = message

		} else if message == "" { //message vide veut dire toutes les données du graphe ont été envoyées.

			liste_n := list_noeuds(Graphe)          // génére la liste des noeuds à partir de du graphe
			dist, prev := initialiser(len(liste_n)) // initialise les Arrays dont nous aurons besoins

			calcule(conn, wg ,Graphe, liste_n, dist, prev, Nom_resultat)
		

		} else {
			//ajout des des arêtes au slice
			Graphe = addToSlice(message, Graphe)

		}

	}

	check(scanner.Err())
}

func main() {

	ln, err := net.Listen("tcp", "localhost:6666")
	check(err)
	fmt.Println("Accept connection on port  : 6666")

	for {
		conn, err := ln.Accept() // accepte la connexion tcp
		check(err)
		fmt.Println("Connection accéptée")
		go handleConnection(conn)
	}
}
