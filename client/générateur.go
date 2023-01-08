package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func check1(e error) {
	if e != nil {
		panic(e)
	}
}

func random_network() {
	var n, a int
	var name string

	fmt.Println("Nom de votre graphe? exemple.txt")
	fmt.Scanf("%s", &name)

	fmt.Println("Nombre de noeuds?")
	fmt.Scanf("%d", &n)

	fmt.Println("poids max?")
	fmt.Scanf("%d", &a)

	fmt.Println(name, n, a)

	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0600)
	defer file.Close()
	check1(err)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {

				_, err = file.WriteString(strconv.Itoa(i) + " " + strconv.Itoa(j) + " " + strconv.Itoa(rand.Intn(a)) + "\n")

			}
		}
	}
	_, err = file.WriteString("\n")

}
func main() {
	d := time.Now()
	random_network()
	f := time.Now()
	temps := f.Sub(d)
	fmt.Println(temps)
}
