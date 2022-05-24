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
	
	requestBody, _ := ioutil.ReadAll(r.Body)
	var shark Shark 

	stmt, err := db.Prepare("INSERT INTO allSharks(name, scienceName, description, imageUrl, memeUrl) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(requestBody, &shark)
	_, err = stmt.Exec(shark.Name, shark.ScienceName, shark.Description, shark.ImageURL, shark.MemeURL)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}

func getOneshark(w http.ResponseWriter, r *http.Request) {
	// on récup l'ID passé en paramètres
	w.Header().Set("Content-Type", "application/json")
	sharkID := mux.Vars(r)["id"]
	var shark Shark 
	
	result, err := db.Query("SELECT id, name, scienceName, description, imageUrl, memeUrl from allSharks WHERE id = ?", sharkID)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
	err := result.Scan(&shark.ID, &shark.Name, &shark.ScienceName, &shark.Description, &shark.ImageURL, &shark.MemeURL)
	if err != nil {
		panic(err.Error())
	}
}

	json.NewEncoder(w).Encode(shark)
}

func getAllsharks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
  	var sharks []Shark
  	result, err := db.Query("SELECT id, name, scienceName, description, imageUrl, memeUrl from allSharks")
  if err != nil {
    panic(err.Error())
  }
  defer result.Close()

  for result.Next() {
    var shark Shark
    err := result.Scan(&shark.ID, &shark.Name, &shark.ScienceName, &shark.Description, &shark.ImageURL, &shark.MemeURL)
    if err != nil {
      panic(err.Error())
    }
    sharks = append(sharks, shark)
  }
  json.NewEncoder(w).Encode(sharks)
}

func getList(w http.ResponseWriter, r *http.Request) {
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

// func updateshark(w http.ResponseWriter, r *http.Request) {
// 	// on récup l'id passé en paramètes
// 	sharkID := mux.Vars(r)["id"]
// 	// on crée une variable updatedshark, destinée à recevoir un objet de type Shark
// 	var updatedshark Shark
	
// 	parseJson := parsingJson()

// 	// on récup le corps de la requête, en affichant une erreur si c'est mal formaté et on le 
// 	// Unmarshal afin de le stocker dans notre variable updatedShark 
// 	reqBody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Fprintf(w, "Kindly enter data with the shark infos in order to update")
// 	}
// 	json.Unmarshal(reqBody, &updatedshark)

// 	// une boucle for pour chercher le requin concerné et mettre à jour avec les infos d'updatedShark
// 	for i, singleshark := range parseJson.Allsharks {
// 		if singleshark.ID == sharkID {
// 			singleshark.Name = updatedshark.Name
// 			singleshark.Description = updatedshark.Description
// 			// on modifie le tableau de requin avec le requin mis à jour 
// 			parseJson.Allsharks = append(parseJson.Allsharks[:i], singleshark)
// 			// on envoi en réponse le requin modifié 
// 			json.NewEncoder(w).Encode(singleshark)
// 		}
// 	}

// 	// afin de pouvoir l'écrire dans le Json, on Marshal notre parseJson 
// 	modifJson, err := json.Marshal(parseJson)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// on écrit notre modifJson (parseJson "marshalisé") dans le fichier sharks.json grace à ioutil 
// 	err = ioutil.WriteFile("sharks.json", modifJson, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func deleteshark(w http.ResponseWriter, r *http.Request) {
// 	// on récup le Json parsé, comme d'hab 
// 	parseJson := parsingJson()
// 	// on récup l'id passé en paramètes
// 	sharkID := mux.Vars(r)["id"]

// 	// une boucle for pour chercher le requin concerné
// 	for i, singleshark := range parseJson.Allsharks {
// 		if singleshark.ID == sharkID {
// 			// on supprime le requin concerné en décalant les valeurs du tableau vers la gauche 
// 			// à partir de l'ID trouvé
// 			parseJson.Allsharks = append(parseJson.Allsharks[:i], parseJson.Allsharks[i+1:]...)
// 			fmt.Fprintf(w, "The shark with ID %v has been deleted successfully", sharkID)
// 		}
// 	}

// 	// afin de pouvoir l'écrire dans le Json, on Marshal notre parseJson 
// 	modifJson, err := json.Marshal(parseJson)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// on écrit notre modifJson (parseJson "marshalisé") dans le fichier sharks.json grace à ioutil 
// 	err = ioutil.WriteFile("sharks.json", modifJson, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

var db *sql.DB
var err error

func main() {
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
	// router.HandleFunc("/sharks/{id}", updateshark).Methods("PATCH")
	// router.HandleFunc("/sharks/{id}", deleteshark).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
