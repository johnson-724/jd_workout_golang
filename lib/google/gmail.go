package google

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	fileLib "jd_workout_golang/lib/file"
	"log"
	"net/http"
	"os"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "credential/google/token.json"
	tok, err := tokenFromFile(tokFile)

	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	return config.Client(context.Background(), tok)
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)

	path = getAbsPath(path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}

	defer f.Close()

	json.NewEncoder(f).Encode(token)
}

func getAbsPath(file string) string {
	path := fileLib.AccessFromCurrentDir(file)

	return path
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	path := getAbsPath(file)

	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return tok
}

func Send(from string, to string, title string, message string) {
	ctx := context.Background()
	path := getAbsPath("credential/google/credentials.json")

	c, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	config, err := google.ConfigFromJSON(c, gmail.GmailSendScope)

	if err != nil {
		panic(err)
	}

	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		panic(err)
	}

	msgStr := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, title, message)
	msg := []byte(msgStr)
	
	gMessage := &gmail.Message{Raw: base64.URLEncoding.EncodeToString(msg)}

	// Send the message
	res, err := srv.Users.Messages.Send("me", gMessage).Do()

	if err != nil {
		
		panic(err)
	}

	fmt.Printf("Message sent: %v\n", res)
}
