package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session *mgo.Session
var db *mgo.Database

func main() {
	defer openDatabase()()
	router := mux.NewRouter()
	router.HandleFunc("/api/cats", getCats).Methods("GET")
	router.HandleFunc("/api/cats", postCats).Methods("POST")
	router.HandleFunc("/api/cats/{catId}", getOneCat).Methods("GET")
	router.HandleFunc("/api/cats/{catId}", deleteOneCat).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Application started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getCats(w http.ResponseWriter, r *http.Request) {
	cats := []Cat{}
	db.C("cats").Find(nil).All(&cats)
	w.Header().Set("Content-Type", "application/json")
	w.Write(toJSON(cats))
}

func postCats(w http.ResponseWriter, r *http.Request) {
	var cat = new(Cat)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(cat)
	db.C("cats").Insert(cat)
	w.WriteHeader(201)
}

func getOneCat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var cat = new(Cat)
	err := db.C("cats").Find(bson.M{"_id": bson.ObjectIdHex(vars["catId"])}).One(&cat)
	if err != nil {
		w.WriteHeader(404)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(toJSON(cat))
	}
}

func deleteOneCat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := db.C("cats").RemoveId(bson.ObjectIdHex(vars["catId"]))
	var statusCode int
	if err == nil {
		statusCode = 204
	} else {
		statusCode = 404
	}
	w.WriteHeader(statusCode)
}

var initContext sync.Once

func openDatabase() func() {
	initContext.Do(func() {
		var err error
		mongoDbAddr := os.Getenv("MONGO_PORT_27017_TCP_ADDR")
		if mongoDbAddr == "" {
			mongoDbAddr = "localhost"
		}
		session, err = mgo.Dial(mongoDbAddr)
		if err != nil {
			panic(err)
		}
		db = session.DB("cats_api")
	})
	return func() {
		initContext.Do(func() {
			session.Close()
		})
	}
}

func toJSON(input interface{}) []byte {
	b, _ := json.Marshal(input)
	return b
}

// Cat represents a cat
type Cat struct {
	ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name  string        `json:"name"`
	Color string        `json:"color"`
}
