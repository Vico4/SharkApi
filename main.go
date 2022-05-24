package main

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
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


func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

// func createshark(w http.ResponseWriter, r *http.Request) {
// 	// on déclare une variable de type Shark qui recevra les infos de notre nouveau requin 
// 	var newshark Shark
// 	// on récup de Json parsé 
// 	parseJson := parsingJson()
// 	// on récup le corps de la requête, en affichant une erreur si c'est mal formaté et on le 
// 	// Unmarshal afin de le stocker dans notre variable newshark 
// 	reqBody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Fprintf(w, "Kindly enter data with the shark informations in order to update")
// 	}

// 	json.Unmarshal(reqBody, &newshark)

// 	// on ajoute newShark au tableau de requins dans notre objet parseJson 
// 	// mais attention à ce stade là, rien n'est encore écrit dans notre fichier Json
// 	parseJson.Allsharks = append(parseJson.Allsharks, newshark)

// 	fmt.Println(parseJson)

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

// 	// on envoi le nouveau requin créé en réponse à la requête 
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(newshark)
// }

// func getOneshark(w http.ResponseWriter, r *http.Request) {

// 	// on récup le json parsé dans parseJson
// 	parseJson := parsingJson()
// 	// on récup l'ID passé en paramètres
// 	sharkID := mux.Vars(r)["id"]

// 	// avec une boucle for, on cherche le requin dont l'id correspond dans notre tableau de requins 
// 	for _, singleshark := range parseJson.Allsharks {
// 		if singleshark.ID == sharkID {
// 			// on envoi le requin trouvé en réponse à la requête 
// 			json.NewEncoder(w).Encode(singleshark)
// 		}
// 	}
// }

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

// func parsingJson() Allsharks {

// 	// Open our jsonFile
// 	jsonFile, err := os.Open("sharks.json")
// 	// if we os.Open returns an error then handle it
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println("Successfully Opened sharks.json")
// 	// defer the closing of our jsonFile so that we can parse it later on
// 	defer jsonFile.Close()

// 	byteValue, _ := ioutil.ReadAll(jsonFile)

// 	// we initialize our Users array
// 	var sharks Allsharks
// 	// we unmarshal our byteArray which contains our
// 	// jsonFile's content into 'users' which we defined above
// 	json.Unmarshal(byteValue, &sharks)
// 	fmt.Println(sharks)

// 	return sharks
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

//   insert, err := db.Query("INSERT INTO allSharks(name) VALUES('Great White Shark')")
//   if err != nil {
//     panic(err.Error())
//   }
//   defer insert.Close()


//	parseSharks := parsingJson()
	router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/", homeLink)
	// router.HandleFunc("/shark", createshark).Methods("POST")
	router.HandleFunc("/sharks", getAllsharks).Methods("GET")
	// router.HandleFunc("/sharks/{id}", getOneshark).Methods("GET")
	router.HandleFunc("/sharklist", getList).Methods("GET")
	// router.HandleFunc("/sharks/{id}", updateshark).Methods("PATCH")
	// router.HandleFunc("/sharks/{id}", deleteshark).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
