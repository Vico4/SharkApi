package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Shark struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	RealName	string	`json:"realName"`
	Description string  `json:"description"`
	ImageName   string  `json:"imageName"`
	MemeName    string  `json:"memeName"`
}

type Allsharks struct {
	Allsharks []Shark `json:"sharks"`
}

type List struct {
	List []Shortshark `json:"sharks"`
}

type Shortshark struct {
	ID        string `json:"id"`
	ImageName string `json:"imageName"`
	Name      string `json:"name"`
}


func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createshark(w http.ResponseWriter, r *http.Request) {
	var newshark Shark
	parseJson := parsingJson()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the shark informations in order to update")
	}

	json.Unmarshal(reqBody, &newshark)
	parseJson.Allsharks = append(parseJson.Allsharks, newshark)

	fmt.Println(parseJson)

	modifJson, err := json.Marshal(parseJson)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("sharks.json", modifJson, 0644)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newshark)
}

func getOneshark(w http.ResponseWriter, r *http.Request) {

	parseJson := parsingJson()
	sharkID := mux.Vars(r)["id"]

	for _, singleshark := range parseJson.Allsharks {
		if singleshark.ID == sharkID {
			json.NewEncoder(w).Encode(singleshark)
		}
	}
}

func getAllsharks(w http.ResponseWriter, r *http.Request) {
	parseJson := parsingJson()
	json.NewEncoder(w).Encode(parseJson)
}

func getList(w http.ResponseWriter, r *http.Request) {
	parseJson := parsingJson()
	leng := len(parseJson.Allsharks)
	vlist := make([]Shortshark, leng)
	for i, singleshark := range parseJson.Allsharks {
		vlist[i].Name = singleshark.Name
		vlist[i].ID = singleshark.ID
		vlist[i].ImageName = singleshark.ImageName
	}
	json.NewEncoder(w).Encode(vlist)
}

func updateshark(w http.ResponseWriter, r *http.Request) {
	sharkID := mux.Vars(r)["id"]
	var updatedshark Shark
	parseJson := parsingJson()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the shark infos in order to update")
	}
	json.Unmarshal(reqBody, &updatedshark)

	for i, singleshark := range parseJson.Allsharks {
		if singleshark.ID == sharkID {
			singleshark.Name = updatedshark.Name
			singleshark.Description = updatedshark.Description
			parseJson.Allsharks = append(parseJson.Allsharks[:i], singleshark)
			json.NewEncoder(w).Encode(singleshark)
		}
	}

	modifJson, err := json.Marshal(parseJson)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("sharks.json", modifJson, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func deleteshark(w http.ResponseWriter, r *http.Request) {
	parseJson := parsingJson()
	sharkID := mux.Vars(r)["id"]

	for i, singleshark := range parseJson.Allsharks {
		if singleshark.ID == sharkID {
			parseJson.Allsharks = append(parseJson.Allsharks[:i], parseJson.Allsharks[i+1:]...)
			fmt.Fprintf(w, "The shark with ID %v has been deleted successfully", sharkID)
		}
	}

	fmt.Println(parseJson)

	modifJson, err := json.Marshal(parseJson)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("sharks.json", modifJson, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func parsingJson() Allsharks {

	// Open our jsonFile
	jsonFile, err := os.Open("sharks.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened sharks.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var sharks Allsharks
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &sharks)
	fmt.Println(sharks)

	return sharks
}

func main() {
	parseSharks := parsingJson()
	fmt.Println(parseSharks.Allsharks)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/shark", createshark).Methods("POST")
	router.HandleFunc("/sharks", getAllsharks).Methods("GET")
	router.HandleFunc("/sharks/{id}", getOneshark).Methods("GET")
	router.HandleFunc("/sharklist", getList).Methods("GET")
	router.HandleFunc("/sharks/{id}", updateshark).Methods("PATCH")
	router.HandleFunc("/sharks/{id}", deleteshark).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
