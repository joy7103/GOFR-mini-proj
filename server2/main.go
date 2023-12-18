package main

import (
	_"gofr.dev/pkg/gofr"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Task represents the data structure for a task
type Task struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Task      string             `json:"task,omitempty" bson:"task,omitempty"`
	Completed bool               `json:"completed,omitempty" bson:"completed,omitempty"`
}
type TaskResponse struct {
	IDHex     string `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

var client *mongo.Client

func init() {
	// Initialize MongoDB connection
	mongoURI := "mongodb+srv://joy7103:rexxEq-3nujvu-vokdac@gofr.yiqorce.mongodb.net/?retryWrites=true&w=majority"
	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		os.Exit(1)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Error pinging MongoDB:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to MongoDB.")
}

func main() {
	router := mux.NewRouter()

	// Set up routes
	router.HandleFunc("/api/tasks", GetTasks).Methods("GET")
	router.HandleFunc("/api/tasks", CreateTask).Methods("POST")
	router.HandleFunc("/api/tasks/{id}", UpdateTask).Methods("PUT")
	router.HandleFunc("/api/tasks/{id}", DeleteTask).Methods("DELETE")

	// Create CORS handler
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
		Debug:          true,
	})

	// Use CORS handler as the router's handler
	http.Handle("/", corsHandler.Handler(router))

	port := ":8080"
	fmt.Printf("Listening on port %s...\n", port)
	http.ListenAndServe(port, nil)
}


func GetTasks(w http.ResponseWriter, r *http.Request) {
    log.Printf("Request URL: %s %s", r.Method, r.URL.Path)

    collection := client.Database("gofr").Collection("tasks")

    cursor, err := collection.Find(context.Background(), bson.D{})
    if err != nil {
        log.Printf("Error fetching tasks from MongoDB: %v", err)
        http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    var tasks []Task
    if err := cursor.All(context.Background(), &tasks); err != nil {
        log.Printf("Error decoding tasks: %v", err)
        http.Error(w, "Error decoding tasks", http.StatusInternalServerError)
        return
    }

    // Convert ObjectIDs to hex representation
    var tasksResponse []TaskResponse
    for _, t := range tasks {
        tasksResponse = append(tasksResponse, TaskResponse{
            IDHex:     t.ID.Hex(),
            Task:      t.Task,
            Completed: t.Completed,
        })
    }

    if len(tasksResponse) == 0 {
        log.Println("No tasks found")
    } else {
        log.Println("Retrieved tasks:", tasksResponse)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasksResponse)
}




// CreateTask handles POST requests to create a new task
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	collection := client.Database("gofr").Collection("tasks")
	result, err := collection.InsertOne(context.Background(), newTask)
	if err != nil {
		http.Error(w, "Error inserting task", http.StatusInternalServerError)
		return
	}

	// Use the InsertedID type assertion to get the ObjectID
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		http.Error(w, "Error getting inserted ID", http.StatusInternalServerError)
		return
	}

	newTask.ID = insertedID // Assign the ObjectID directly

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}
// UpdateTask handles PUT requests to update a task
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskID := params["id"]

	// Convert the taskID from string to ObjectID
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		http.Error(w, "Invalid task ID format", http.StatusBadRequest)
		return
	}

	var updatedTask Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	collection := client.Database("gofr").Collection("tasks")
	_, err = collection.UpdateByID(context.Background(), objID, bson.M{"$set": updatedTask})
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}


// DeleteTask handles DELETE requests to delete a task
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskID, ok := params["id"]
	if !ok {
		log.Println("Task ID not provided")
		http.Error(w, "Task ID not provided", http.StatusBadRequest)
		return
	}

	log.Printf("Deleting task with ID: %v", taskID)

	// Convert the taskID from string to ObjectID
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		log.Printf("Invalid task ID format: %v", err)
		http.Error(w, "Invalid task ID format", http.StatusBadRequest)
		return
	}

	collection := client.Database("gofr").Collection("tasks")
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		log.Printf("Error deleting task: %v", err)
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		log.Println("Task not found")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	log.Println("Task deleted successfully")

	w.WriteHeader(http.StatusNoContent)
}

