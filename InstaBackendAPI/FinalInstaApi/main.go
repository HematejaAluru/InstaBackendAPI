package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"quickstart/helper"
	"quickstart/models"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Connection mongoDB with helper class
var client = helper.ConnectDB()
var Userscollection = client.Database("Instagram_api").Collection("Users")
var Postscollection = client.Database("Instagram_api").Collection("Posts")
var Users []models.User
var Posts []models.Post

func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
func connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

func ping(client *mongo.Client, ctx context.Context) error {

	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occored, then
	// the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/users", createUser)
	http.HandleFunc("/users/", getUser)
	http.HandleFunc("/posts", createPost)
	http.HandleFunc("/posts/", getPost)
	http.HandleFunc("/posts/users/", getAllPosts)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tempUser models.User
	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&tempUser)
	sum := sha256.Sum256([]byte(tempUser.Password))
	pass := hex.EncodeToString(sum[:])
	tempUser.Password = pass

	// insert our User model.
	result, err := Userscollection.InsertOne(context.TODO(), tempUser)
	Userscollection.UpdateOne(context.TODO(), bson.M{"_id": result.InsertedID}, bson.D{{"$set", bson.D{{"id", result.InsertedID}}}})
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	// if only one expected
	param1 := path.Base(r.URL.Path)
	if param1 != "" {
		fmt.Print(param1)
	}
	var retreviedUser models.User
	ObjectIDParam, errr := primitive.ObjectIDFromHex(param1)
	if errr != nil {
		log.Fatal(errr)
		return
	}
	err := Userscollection.FindOne(r.Context(), bson.M{"id": ObjectIDParam}).Decode(&retreviedUser)
	if err != nil {
		log.Fatal(err)
		return
	}
	json.NewEncoder(w).Encode(retreviedUser)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//var tempUser models.User
	var recievedinfo bson.M
	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&recievedinfo)
	var tempID models.ID
	var tempPost models.Post
	if str, ok := recievedinfo["Id"].(string); ok {
		tempID.UserID = str
	}
	if str, ok := recievedinfo["Caption"].(string); ok {
		tempPost.Caption = str
	}
	if str, ok := recievedinfo["ImageURL"].(string); ok {
		tempPost.ImageURL = str
	}
	if str, ok := recievedinfo["PostedTS"].(string); ok {
		tempPost.PostedTS = str
	} else {
		tempPost.PostedTS = time.Now().Format("02-Jan-2006 15:04:05")
	}
	tempPost.Id = tempID
	result, err := Postscollection.InsertOne(context.TODO(), tempPost)
	Postscollection.UpdateOne(context.TODO(), bson.M{"_id": result.InsertedID}, bson.D{{"$set", bson.D{{"id.postid", result.InsertedID}}}})
	if err != nil {
		helper.GetError(err, w)
		return
	}
	Postscollection.FindOne(context.TODO(), bson.M{"_id": result.InsertedID}).Decode(&tempPost)
	json.NewEncoder(w).Encode(result)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	// if only one expected
	param1 := path.Base(r.URL.Path)
	if param1 != "" {
		fmt.Print(param1)
	}
	var retreviedPost models.Post
	ObjectIDParam, errr := primitive.ObjectIDFromHex(param1)
	if errr != nil {
		log.Fatal(errr)
		return
	}
	err := Postscollection.FindOne(r.Context(), bson.M{"id.postid": ObjectIDParam}).Decode(&retreviedPost)
	if err != nil {
		log.Fatal(err)
	}
	Post := retreviedPost
	json.NewEncoder(w).Encode(Post)
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
	Id := path.Base(r.URL.Path)
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	if limit != "" && offset != "" {
		if Id != "" {
			fmt.Print(Id)
		}
		var retrievedPosts []models.Post

		cursor, err := Postscollection.Find(context.TODO(), bson.M{"id.userid": Id})
		if err != nil {
			helper.GetError(err, w) // prints 'document is nil'
		}
		var results []models.Post
		if err = cursor.All(context.TODO(), &results); err != nil {
			log.Fatal(err)
		}
		var offsetInt int
		var limitInt int
		offsetInt, err = strconv.Atoi(offset)
		limitInt, err = strconv.Atoi(limit)
		var count int
		count = 0
		for i := offsetInt; i < len(results); i++ {
			if count == limitInt {
				break
			}
			if results[i].Id.UserID == Id {
				retrievedPosts = append(retrievedPosts, results[i])
				count += 1
			}
		}
		json.NewEncoder(w).Encode(retrievedPosts)
	} else {
		if Id != "" {
			fmt.Print(Id)
		}
		var retrievedPosts []models.Post
		cursor, err := Postscollection.Find(r.Context(), bson.M{"id.userid": Id})
		if err != nil {
			helper.GetError(err, w)
		}
		var results []models.Post
		if err = cursor.All(context.TODO(), &results); err != nil {
			log.Fatal(err)
		}
		for _, result := range results {
			if result.Id.UserID == Id {
				retrievedPosts = append(retrievedPosts, result)
			}
		}
		json.NewEncoder(w).Encode(retrievedPosts)
	}
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://dbUser:avnmht123@cluster0.4us7d.mongodb.net/Cluster0?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	ping(client, ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	handleRequests()
}
