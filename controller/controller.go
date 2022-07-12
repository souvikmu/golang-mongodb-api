package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/souvikmukherjee/mongoapi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// for connection string using mongodb atlas
// fetch the password
// var password string = os.Getenv("passwordstring")
// var username string = os.Getenv("usernamestring")

// for local system storage
// connection string
var connectionString string =  "mongodb://127.0.0.1:27017"

const dbName = "netflix"
const colName = "watchlist"

// define collection
var collection *mongo.Collection

// connect to mongodb database
func init() {
	// client options
	fmt.Println(connectionString)
	clientOption := options.Client().ApplyURI(connectionString)

	// connect to db
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("mongodb connection success...")

	// collection
	collection = client.Database(dbName).Collection(colName)

	// collection instance
	fmt.Println("collection reference is ready")

}

// database helper functions - seperate file
func insertOneRecord(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted one movie into db with id: ", inserted.InsertedID)
}

// update one record
func updateOneRecord(movieId string) {
	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified entry: ", result.ModifiedCount)

}

// delete one record
func deleteOneRecord(movieId string) {
	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("deleted entry count: ", deleteCount)
}

// delete all records
func deleteAllRecords() {
	delCount, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("entry delete count is: ", delCount.DeletedCount)
}

// get all records
func getAllRecords() []primitive.M {
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var movies []primitive.M

	for cursor.Next(context.Background()) {
		var movie bson.M
		err := cursor.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}

	return movies
}

// Router controllers - meant for this file originally

// get all movies
func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	allMovies := getAllRecords()

	if len(allMovies) == 0 {
		json.NewEncoder(w).Encode("database empty")
		return
	}
	json.NewEncoder(w).Encode(allMovies)
}

// create a movie
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	body := r.Body
	defer body.Close()
	if body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "empty body sent"})
		return
	}

	var movie model.Netflix
	err := json.NewDecoder(body).Decode(&movie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "could not decode json"})
		log.Fatal(err)
	}
	insertOneRecord(movie)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("data entry successful, movie id: %s", movie.ID),
	})
}

// mark movies as watched
func WatchMarker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	updateOneRecord(params["id"])
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("movie marked as watched, movie id: %s", params["id"]),
	})
}

// delete one movie
func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneRecord(params["id"])
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("deleted movie id: %s", params["id"]),
	})
}

// delete all movies
func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	deleteAllRecords()
	json.NewEncoder(w).Encode(map[string]string{
		"message": "all records perished",
	})
}
