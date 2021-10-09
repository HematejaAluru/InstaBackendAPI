package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CreateUser(name string, email string, password string) {
	values := map[string]string{"Name": name, "Email": email, "Password": password}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://localhost:10000/users", "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Println(res["json"])
}
func CreatePost(id string, caption string, imageurl string, postts string) {
	values := map[string]string{"Id": id, "Caption": caption, "ImageURL": imageurl, "PostTS": postts}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://localhost:10000/posts", "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)
}
func GetPosts(id string) {
	resp, err := http.Get("http://localhost:10000/posts/users/" + id)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)
}
func GetPostsPagenation(id string, limit string, offset string) {
	resp, err := http.Get("http://localhost:10000/posts/users/" + id + "?limit=" + limit + "&offset=" + offset)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)
}

func main() {
	CreateUser("Teja2", "avnmht123@gmail.com", "avnmht123")

	for i := 0; i < 10; i++ {
		CreatePost("6161a8a3e04d0f0312f854cc", "Can you do it?"+strconv.Itoa(i), "www.google.com"+strconv.Itoa(i), time.Now().Format("02-Jan-2006 15:04:05"))
	}
	GetPosts("6161a8a3e04d0f0312f854cc")

	GetPostsPagenation("6161a0d1b250c93c56da96a4", "5", "1")

	CreatePost("6161a8a3e04d0f0312f854cc", "Can you do it?Final3", "www.google.com", time.Now().Format("02-Jan-2006 15:04:05"))
}
