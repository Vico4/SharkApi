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

func updateshark(w http.ResponseWriter, r *http.Request) {
	// on récup l'id passé en paramètes
	sharkID := mux.Vars(r)["id"]
	// on crée une variable updatedshark, destinée à recevoir un objet de type Shark
	var shark Shark
	// on récup le corps de la requête, en affichant une erreur si c'est mal formaté et on le 
	// Unmarshal afin de le stocker dans notre variable updatedShark 
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the shark infos in order to update")
	}
	json.Unmarshal(reqBody, &shark)

	stmt, err := db.Prepare("UPDATE allSharks SET name = ?, scienceName = ?, description = ?, imageUrl = ?, memeUrl = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(shark.Name, shark.ScienceName, shark.Description, shark.ImageURL, shark.MemeURL, sharkID)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Shark was updated")
}

func deleteshark(w http.ResponseWriter, r *http.Request) {
	
	sharkID := mux.Vars(r)["id"]

	_, err := db.Exec("DELETE FROM allSharks WHERE id = ?", sharkID)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "The shark with ID %v has been deleted successfully", sharkID)
}

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
	router.HandleFunc("/sharks/{id}", updateshark).Methods("PATCH")
	router.HandleFunc("/sharks/{id}", deleteshark).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
