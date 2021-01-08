package data

import (
	"aurashort/server/database"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"strconv"
)

var (
	db        = database.Connect("aurashort")
	shortInfo = db.Collection("shortInfo")
)

// CreateLink is used to create a new short link with a random id
func CreateLink(w http.ResponseWriter, r *http.Request) {
	if CheckRequestLimit(r.RequestURI, w) {
		w.Header().Set("connection-type", "application/json")

		redirect := r.Header.Get("Redirect")

		if redirect == "" {
			http.Error(w, "missing redirect header", http.StatusBadRequest)
			return
		}

		length, err := strconv.Atoi(os.Getenv("LENGTH_OF_LINKS"))
		if err != nil {
			http.Error(w, "failed to create short link", http.StatusBadRequest)
			return
		}

		randomId := randNum(length)

		if checkLinkAvailability(randomId) {

			link := ShortInfo{
				Id:       randomId,
				Redirect: redirect,
				Uses:     0,
			}

			_, err = shortInfo.InsertOne(context.TODO(), &link)
			if err != nil {
				http.Error(w, "failed to create short link", http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(link)
		}
	} else {
		CreateLink(w, r)
	}
}

// CreateCustomLink is used to create a custom short link
func CreateCustomLink(w http.ResponseWriter, r *http.Request) {
	if CheckRequestLimit(r.RequestURI, w) {
		w.Header().Set("connection-type", "application/json")

		redirect := r.Header.Get("Redirect")
		custom := r.Header.Get("Custom")

		if redirect == "" {
			http.Error(w, "missing redirect header field", http.StatusBadRequest)
			return
		}

		if custom == "" {
			http.Error(w, "missing custom header field", http.StatusBadRequest)
			return
		}

		if checkLinkAvailability(custom) {

			link := ShortInfo{
				Id:       custom,
				Redirect: redirect,
				Uses:     0,
			}

			_, err := shortInfo.InsertOne(context.TODO(), &link)
			if err != nil {
				http.Error(w, "failed to create short link", http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(link)
		} else {
			http.Error(w, "custom link already in use", http.StatusBadRequest)
			return
		}
	}
}

// Redirect is used to redirect a user to a certain domain
func Redirect(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)
	var id = params["id"]

	var info ShortInfo

	filter := bson.M{"id": id}
	err := shortInfo.FindOne(context.TODO(), filter).Decode(&info)
	if err != nil {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	update := bson.D{
		{"$set", bson.D{
			{"id", info.Id},
			{"redirect", info.Redirect},
			{"uses", info.Uses + 1},
		}},
	}

	result := shortInfo.FindOneAndUpdate(context.TODO(), filter, update)
	if result.Err() == mongo.ErrNoDocuments {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, info.Redirect, http.StatusMovedPermanently)
}

func checkLinkAvailability(id string) bool {
	filter := bson.M{"id": id}
	err := shortInfo.FindOne(context.TODO(), filter)
	if err != nil {
		return false
	}
	return true
}

func randNum(length int) string {
	num := make([]byte, length)
	_, err := rand.Read(num)
	if err != nil {
		fmt.Println("Failed to generate random number")
		fmt.Println(err)
		return randNum(length)
	}
	return fmt.Sprintf("%x", num)
}
