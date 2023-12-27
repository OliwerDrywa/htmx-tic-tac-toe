package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RandomUserResponse struct {
	Results []struct {
		Name struct {
			Title string `json:"title"`
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Picture struct {
			Large     string `json:"large"`
			Medium    string `json:"medium"`
			Thumbnail string `json:"thumbnail"`
		} `json:"picture"`
	} `json:"results"`
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
