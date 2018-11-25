package main

import (
	"golang.org/x/oauth2"
	"net/http"
	"fmt"
	"log"
	"os"
	"encoding/json"
	"google.golang.org/api/calendar/v3"
	"golang.org/x/net/context"
	"errors"
	"io/ioutil"
	"golang.org/x/oauth2/google"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func getCalenderId(calService *calendar.Service, calendarName string) (string, error) {
	listRes, err := calService.CalendarList.List().Fields("items/id", "items/summary").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of calendars: %v", err)
	}
	for _, v := range listRes.Items {
		log.Printf("Calendar ID: %v\n with name %v", v.Id, v.Summary)
		if v.Summary == calendarName {
			return v.Id, nil
		}
	}
	return "", errors.New("calendar not found")
}

func getEventsFromCalender(calenderName string) (Items []*calendar.Event){
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	calId, err :=  getCalenderId(srv, calenderName)
	if err != nil {
		log.Fatalf("Unable to find calender: %v", err)
	}

	events := getEvents(srv, calId, "")
	eventItems := events.Items
	pageToken := events.NextPageToken

	 for pageToken != "" {
		 events = getEvents(srv, calId, pageToken)
		 eventItems  = append(eventItems, events.Items...)
		 pageToken = events.NextPageToken
	 }
	 return eventItems
}

func getEvents(srv *calendar.Service, calId string, pageToken string)(*calendar.Events) {
	events, err := srv.Events.List(calId).OrderBy("startTime").ShowDeleted(false).
		PageToken(pageToken).SingleEvents(true).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		panic(err)
	}

	return events
}