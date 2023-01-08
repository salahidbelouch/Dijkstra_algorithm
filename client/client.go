package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// demande nom de fichier résultatl et l'emplacement du graphe à résoudre
func Arg() (string, string) {
	var path, res string

	fmt.Println("Comment souhaitez-vous appeler votre fichier résultat?  : exemple ")
	fmt.Scanf("%s", &res)
	
	fmt.Println("emplacement de votre graphe? exemple :/emplacement/du/fichier.txt")
	fmt.Scanf("%s", &path)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		fmt.Println("la saisie de l'emplacement de votre graphe est incorrect... réessayez:" + "\n")
		return Arg()
	}

	return res, path

}

// pour gerer les erreur
func check(e error) {
	if e != nil {
		panic(e)
	}
}
func main() {

	// connexion au serveur

	con, err := net.Dial("tcp", "localhost:6666")
	check(err)
	defer con.Close() //fermeture de la connexion à la fin d'éxécution
	fmt.Println("connecté")

	// demande nom de fichier résultatl et l'emplacement du graphe à résoudre
	res, path := Arg()

	if _, err = con.Write([]byte(res + "\n")); err != nil {
		log.Printf("envoie au serveur échoué: %v\n", err)
	}

	fmt.Println("envoie des donnée au serveur")
	//crée un reader sur de la connexion pour écouter le serveur
	serverReader := bufio.NewReader(con)

	//ouverture du fichier graphe
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	//itilialisation du scanner pour lire sur le fichier
	scanner := bufio.NewScanner(file)
	check(scanner.Err())

	//boucle infinie
	done := true
	for done {

		//envoie du fichier
		for scanner.Scan() { // tant qu'il y a des token à lire dans le fichier (retourne un false si plus de token ou erreur, retourne True sinon)

			clientRequest := scanner.Text() //lecture de la ligne

			switch err {
			case nil:
				clientRequest_trimed := strings.TrimSpace(clientRequest)
				if _, err = con.Write([]byte(clientRequest_trimed + "\n")); err != nil {
					log.Printf("failed to send the client request: %v\n", err)
				}
			case io.EOF:
				log.Println("client closed the connection")
				return
			default:
				log.Printf("client error:")
				check(err)
				return
			}
			if scanner.Scan() == false {
				fmt.Println("envoie terminé, Merci d'attendre le résultat")
			}
		}
		// attends en boucle le reulstat
		serverResponse, _ := serverReader.ReadString('\n')
		serverResponse = strings.TrimSpace(serverResponse) //suprimme les caractères inutiles (espaces..)
		if len(serverResponse) > 0 {
			fmt.Println("temps de l'algorithme: " + serverResponse)
			done = false
		}

	}

}
