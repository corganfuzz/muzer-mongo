package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

// DB stores db session info
type DB struct {
	session    *mgo.Session
	collection *mgo.Collection
}

type Movie struct {
	ID        bson.ObjectId ` json:"id" bson:"_id,omitempty"`
	Name      string        ` bson:"name"`
	Year      string        ` bson:"year"`
	Directors []string      ` bson:"directors"`
	Writers   []string      ` bson:"writers"`
	BoxOffice ` bson:"boxOffice"`
}

type BoxOffice struct {
	Budget uint64 ` bson:"budget"`
	Gross  uint64 ` bson:"gross"`
}

// GetMovie handler function fetches a movie with given ID
func (db *DB) GetMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var movie Movie
	err := db.collection.Find(bson.M{"_id": bson.ObjectIdHex(vars["id"])}).One(&movie)

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(movie)
		w.Write(response)
	}
}

// PostMovie handler function adds a new movie to our DB collecion
func (db *DB) PostMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	postBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(postBody, &movie)

	//Create a hash ID to insert

	movie.ID = bson.NewObjectId()
	err := db.collection.Insert(movie)

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(movie)
		w.Write(response)
	}
}

// UpdateMovie is a put request
func (db *DB) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var movie Movie
	putBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(putBody, &movie)

	//Create a hash ID to insert

	err := db.collection.Update(bson.M{"_id": bson.IsObjectIdHex(vars["id"])}, bson.M{"$set": &movie})
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text")
		w.Write([]byte("Updated successfully!"))
	}
}

// DeleteMovie deletes a movie
func (db *DB) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Create a hash id to insert

	err := db.collection.Remove(bson.M{"_id": bson.ObjectIdHex(vars["id"])})
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text")
		w.Write([]byte("Deleted successfully!"))
	}
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	c := session.DB("appdb").C("movies")
	db := &DB{session: session, collection: c}

	if err != nil {
		panic(err)
	}

	defer session.Close()

	// Create a new Router
	r := mux.NewRouter()

	// Attach an elegant path with handler
	r.HandleFunc("/v1/movies/{id:[a-zA-Z0-9]*}", db.GetMovie).Methods("GET")
	r.HandleFunc("/v1/movies", db.PostMovie).Methods("POST")
	r.HandleFunc("/v1/movies/{id:[a-zA-Z0-9]*}", db.UpdateMovie).Methods("PUT")
	r.HandleFunc("/v1/movies/{id:[a-zA-Z0-9]*}", db.DeleteMovie).Methods("DELETE")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
