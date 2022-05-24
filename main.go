package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "os"
	"database/sql"
  	_"github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type Shark struct {
	ID          *int  	`json:"id"`
	Name        *string  `json:"name"`
	ScienceName	*string	`json:"scienceName"`
	Description *string  `json:"description"`
	ImageURL   	*string  `json:"imageUrl"`
	MemeURL	   	*string  `json:"memeUrl"`
}

type Allsharks struct {
	Allsharks []Shark `json:"sharks"`
}

type List struct {
	List []Shortshark `json:"sharks"`
}

type Shortshark struct {
	ID        *string `json:"id"`
	Name      *string `json:"name"`
	ImageURL  *string `json:"imageUrl"`
}


func createshark(w http.ResponseWriter, r *http.Request) {

	//on récup le corps de la requête et on l'Unmarshal dans la variable shark de type Shark afin de pouvoir l'exploiter
	requestBody, _ := ioutil.ReadAll(r.Body)
	var shark Shark 
	json.Unmarshal(requestBody, &shark)

	//on prépare notre requête SQL avec des placeholders "?" pour les valeurs
	stmt, err := db.Prepare("INSERT INTO allSharks(name, scienceName, description, imageUrl, memeUrl) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	// on exécute notre requête SQL, les valeurs passées dans Exec remplacerons les placeholders
	_, err = stmt.Exec(shark.Name, shark.ScienceName, shark.Description, shark.ImageURL, shark.MemeURL)
	if err != nil {
		panic(err.Error())
	}

	// on envoi la réponse 
	fmt.Fprintf(w, "New post was created")
}

func getOneshark(w http.ResponseWriter, r *http.Request) {
	// on récup l'ID passé en paramètres
	w.Header().Set("Content-Type", "application/json")
	sharkID := mux.Vars(r)["id"]

	// on crée une variable shark de type Shark
	var shark Shark 
	
	// on effectue la requête SQL et on stock le résultat dans result 
	result, err := db.Query("SELECT id, name, scienceName, description, imageUrl, memeUrl from allSharks WHERE id = ?", sharkID)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	// on Scan result pour inscrire ses valeurs dans notre variable shark 
	for result.Next() {
	err := result.Scan(&shark.ID, &shark.Name, &shark.ScienceName, &shark.Description, &shark.ImageURL, &shark.MemeURL)
	if err != nil {
		panic(err.Error())
	}
}
	// on encode et envoi la réponse 
	json.NewEncoder(w).Encode(shark)
}

func getAllsharks(w http.ResponseWriter, r *http.Request) {
	// on défini le header de la réponse
	w.Header().Set("Content-Type", "application/json")

	// on crée une variable sharks, qui est un tableau d'objets Shark 
  	var sharks []Shark

	// on lance une requête SQL dont on stocke le résultat dans result  
  	result, err := db.Query("SELECT id, name, scienceName, description, imageUrl, memeUrl from allSharks")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	// on boucle sur result et on scan chaque élément, c'est à dire qu'on l'envoi dans une variable shark, 
	// qu'on ajoute ensuite à notre tableau sharks 
	for result.Next() {
		var shark Shark
		err := result.Scan(&shark.ID, &shark.Name, &shark.ScienceName, &shark.Description, &shark.ImageURL, &shark.MemeURL)
		if err != nil {
		panic(err.Error())
		}
		sharks = append(sharks, shark)
	}

	// on encode et on renvoit le tableau sharks en réponse 
	json.NewEncoder(w).Encode(sharks)
}

func getList(w http.ResponseWriter, r *http.Request) {
	// tout pareil que getAllsharks sauf qu'on prend seulement les 3 infos nécessaires 
	// et on renvoie une liste d'objets Shortshark 
	w.Header().Set("Content-Type", "application/json")
  	var list []Shortshark
  	result, err := db.Query("SELECT id, name, imageUrl from allSharks")
  if err != nil {
    panic(err.Error())
  }
  defer result.Close()

  for result.Next() {
    var shark Shortshark
    err := result.Scan(&shark.ID, &shark.Name, &shark.ImageURL)
    if err != nil {
      panic(err.Error())
    }
    list = append(list, shark)
  }
  json.NewEncoder(w).Encode(list)
}

func updateshark(w http.ResponseWriter, r *http.Request) {
	// on récup l'id passé en paramètes
	sharkID := mux.Vars(r)["id"]
	// on crée une variable shark, destinée à recevoir un objet de type Shark
	var shark Shark
	// on récup le corps de la requête, en affichant une erreur si c'est mal formaté et on le 
	// Unmarshal afin de le stocker dans notre variable shark 
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the shark infos in order to update")
	}
	json.Unmarshal(reqBody, &shark)

	// on préparé la requête SQL avec des placeholders "?"
	stmt, err := db.Prepare("UPDATE allSharks SET name = ?, scienceName = ?, description = ?, imageUrl = ?, memeUrl = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	// on exécute la requête SQL, les valeurs passées remplacent les placeholders 
	_, err = stmt.Exec(shark.Name, shark.ScienceName, shark.Description, shark.ImageURL, shark.MemeURL, sharkID)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Shark was updated")
}

func deleteshark(w http.ResponseWriter, r *http.Request) {
	
	// on récup l'ID passé en paramètre 
	sharkID := mux.Vars(r)["id"]

	// on exécute une requête SQL pour supprimer l'élément 
	// (c'est court donc j'ai test d'exécuter la requête sans passer par l'étape db.Prepare ça marche)
	_, err := db.Exec("DELETE FROM allSharks WHERE id = ?", sharkID)
	if err != nil {
		panic(err.Error())
	}

	// on envoie la réponse 
	fmt.Fprintf(w, "The shark with ID %v has been deleted successfully", sharkID)
}

// on crée un objet db et un objet erreur 
var db *sql.DB
var err error

func main() {
	// on ouvre la connection à la bdd et on utilise defer pour lui demander de rester ouverte 
	// jusqu'à ce que qu'on ait fini.
	// attention "Root5003" est le mot de passe de mon serveur SQl, à remplacer par le votre
	// dans l'idéal ce serait une variable d'env 
	db, err = sql.Open("mysql", "root:Root5003@tcp(127.0.0.1:3306)/sharks")
  if err != nil {
    panic(err.Error())
  } else {
	  fmt.Println("Opened DB !")
	}
  defer db.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/shark", createshark).Methods("POST")
	router.HandleFunc("/sharks", getAllsharks).Methods("GET")
	router.HandleFunc("/sharks/{id}", getOneshark).Methods("GET")
	router.HandleFunc("/sharklist", getList).Methods("GET")
	router.HandleFunc("/sharks/{id}", updateshark).Methods("PATCH")
	router.HandleFunc("/sharks/{id}", deleteshark).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
