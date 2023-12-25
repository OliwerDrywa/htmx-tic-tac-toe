package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RandomUserResponse struct {
	Results []struct {
		// Gender string `json:"gender"`
		Name struct {
			Title string `json:"title"`
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		// Location struct {
		// 	Street struct {
		// 		Number int    `json:"number"`
		// 		Name   string `json:"name"`
		// 	} `json:"street"`
		// 	City        string `json:"city"`
		// 	State       string `json:"state"`
		// 	Country     string `json:"country"`
		// 	Postcode    int    `json:"postcode"`
		// 	Coordinates struct {
		// 		Latitude  string `json:"latitude"`
		// 		Longitude string `json:"longitude"`
		// 	} `json:"coordinates"`
		// 	Timezone struct {
		// 		Offset      string `json:"offset"`
		// 		Description string `json:"description"`
		// 	} `json:"timezone"`
		// } `json:"location"`
		// Email string `json:"email"`
		// Login struct {
		// 	UUID     string `json:"uuid"`
		// 	Username string `json:"username"`
		// 	Password string `json:"password"`
		// 	Salt     string `json:"salt"`
		// 	MD5      string `json:"md5"`
		// 	SHA1     string `json:"sha1"`
		// 	SHA256   string `json:"sha256"`
		// } `json:"login"`
		// DOB struct {
		// 	Date string `json:"date"`
		// 	Age  int    `json:"age"`
		// } `json:"dob"`
		// Registered struct {
		// 	Date string `json:"date"`
		// 	Age  int    `json:"age"`
		// } `json:"registered"`
		// Phone string `json:"phone"`
		// Cell  string `json:"cell"`
		// ID    struct {
		// 	Name  string `json:"name"`
		// 	Value string `json:"value"`
		// } `json:"id"`
		Picture struct {
			Large     string `json:"large"`
			Medium    string `json:"medium"`
			Thumbnail string `json:"thumbnail"`
		} `json:"picture"`
		// Nat string `json:"nat"`
	} `json:"results"`
	// Info struct {
	// 	Seed    string `json:"seed"`
	// 	Results int    `json:"results"`
	// 	Page    int    `json:"page"`
	// 	Version string `json:"version"`
	// } `json:"info"`
}

type User struct {
	Icon string
	Name string
}

func RandomUser() (user *User, err error) {
	// Make an HTTP GET request to the randomuser.me API
	response, err := http.Get("https://randomuser.me/api/")
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Parse the JSON response
	var randomUserResponse RandomUserResponse
	err = json.Unmarshal(body, &randomUserResponse)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return nil, err
	}

	randomUser := User{
		randomUserResponse.Results[0].Picture.Thumbnail,
		randomUserResponse.Results[0].Name.First,
	}
	return &randomUser, nil
}
