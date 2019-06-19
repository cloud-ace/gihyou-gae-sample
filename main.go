package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age 	  int32  `json:"age"`
}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/users", putUserHandler).Methods("PUT")
	r.HandleFunc("/users", getUserHandler).Methods("GET")
	http.Handle("/", r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_, err := fmt.Fprint(w, "Hello, World!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func putUserHandler(w http.ResponseWriter, r *http.Request) {
	projectID := os.Getenv("PROJECT_ID")
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	defer client.Close()

	var u *User
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Fatalf("Failed to parse request body: %v", err)
	}
	user := User {
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Age:	   u.Age,
	}
	_, _, err = client.Collection("users").Add(ctx, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println(w)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	projectID := os.Getenv("PROJECT_ID")
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user User
	userList := []User{}
	docs, err := client.Collection("users").Documents(ctx).GetAll()
	if err != nil {
		// handle error
	}
	for _, doc := range docs {
		doc.DataTo(&user)
		userList = append(userList, user)
	}
	w.Header().Set("Content-Type", "application/json")
	hoge, _ := json.Marshal(userList)
	fmt.Fprintln(w, string(hoge))
}
