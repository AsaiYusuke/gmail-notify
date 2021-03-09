package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type gmailAPI struct {
	paths   *paths
	conf    *config
	subDir  string
	service *gmail.Service
}

func (g *gmailAPI) getService() *gmail.Service {
	if g.service == nil {
		g.service = g.createService()
	}
	return g.service
}

func (g *gmailAPI) createService() *gmail.Service {
	jsonKey, err := ioutil.ReadFile(g.paths.getCredentialsFilePath(g.subDir))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(jsonKey, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tokenFile := g.paths.getTokenFilePath(g.subDir)
	token, err := g.tokenFromFile(tokenFile)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}
	ctx := context.Background()
	service, err := gmail.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	return service
}

func (g *gmailAPI) tokenFromFile(filename string) (*oauth2.Token, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)
	return token, err
}

func (g *gmailAPI) getEmail() string {
	profiles, err := g.getService().Users.GetProfile(gmailAPIUserID).Do()
	if err != nil {
		log.Fatalf("Unable to get gmail profile: %v", err)
	}
	return profiles.EmailAddress
}

func (g *gmailAPI) getUnreadMessageIDs() []string {
	retryCount := 3
	var err error
	for retryCount > 0 {
		var messages *gmail.ListMessagesResponse
		messages, err = g.getService().Users.Messages.List(gmailAPIUserID).LabelIds(g.conf.Labels...).Do()
		if err != nil {
			retryCount--
			if retryCount > 0 {
				log.Printf("Retry to get gmail messages: %v", err)
				time.Sleep(time.Second * g.conf.RetryInterval)
			}
			continue
		}
		var newMessageIDs []string
		for _, message := range messages.Messages {
			newMessageIDs = append(newMessageIDs, message.Id)
		}
		return newMessageIDs
	}
	log.Fatalf("Unable to get gmail messages: %v", err)
	return nil
}
