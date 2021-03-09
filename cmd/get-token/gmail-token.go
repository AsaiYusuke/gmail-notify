package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

type gmailToken struct {
}

func (g *gmailToken) createToken(subDir string) {
	jsonKey, err := ioutil.ReadFile(path.Join(subDir, pathCredentialsFilename))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	if _, err := os.Stat(path.Join(subDir, pathTokenFilename)); err == nil || os.IsExist(err) {
		log.Fatalf(`%s exist.`, pathTokenFilename)
	}

	config, err := google.ConfigFromJSON(jsonKey, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	token := g.getTokenFromWeb(config)
	g.saveToken(path.Join(subDir, pathTokenFilename), token)
}

func (g *gmailToken) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return token
}

func (g *gmailToken) saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
