package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GetGoogleClient is a function for obtaining oAuth 2.0 Client from Google
func GetGoogleClient(credentialPath string) *http.Client {
	rawClientSecret, err := ioutil.ReadFile(credentialPath)
	if err != nil {
		log.Fatal("Error reading client secret file", err)
	}

	oAuthConfig, err := google.ConfigFromJSON(rawClientSecret, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatal("Error configuring Google API client", err)
	}

	tokenFileName := "token.json"
	token, err := getTokenFromFile(tokenFileName)
	if err != nil {
		// Obtain token from web and save
		token = getTokenFromWeb(oAuthConfig)
		log.Printf("Saving credential file to %s\n", tokenFileName)
		tokenFile, err := os.OpenFile(tokenFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatalf("Unable to cache token: %v", err)
		}
		defer tokenFile.Close()
		json.NewEncoder(tokenFile).Encode(token)
	}
	return oAuthConfig.Client(context.Background(), token)
}

func getTokenFromFile(tokenFileName string) (*oauth2.Token, error) {
	tokenFile, err := os.Open(tokenFileName)
	if err != nil {
		return nil, err
	}

	defer tokenFile.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(tokenFile).Decode(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	// TODO: Pingpong Flow에서는 이 Scan을 어떻게 처리했는지 확인해보기
	if _, err := fmt.Scan(&code); err != nil {
		panic("Can't read authorization code from user")
	}

	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		panic("Can't get token from user")
	}
	return token
}
