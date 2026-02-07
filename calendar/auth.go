package calendar

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// ReAuthenticate forces a new token retrieval from the web.
func (c *GoogleCalendarIntegration) ReAuthenticate(ctx context.Context) error {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return err
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return err
	}

	tok := getTokenFromWeb(config)
	saveToken("token.json", tok)

	client := config.Client(ctx, tok)
	svr, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	c.calendarService = svr
	return nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Generate a random state for security
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("Unable to generate random state: %v", err)
	}
	state := base64.URLEncoding.EncodeToString(b)

	// Use a channel to receive the code from the HTTP handler
	codeCh := make(chan string)
	errCh := make(chan error)

	// Define a simple handler to catch the code from the redirect
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify state
		if r.URL.Query().Get("state") != state {
			errCh <- fmt.Errorf("invalid state token")
			fmt.Fprintf(w, "Error: Invalid state token.")
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no code in redirect")
			fmt.Fprintf(w, "Error: No code found in the redirect.")
			return
		}
		fmt.Fprintf(w, "Authorization successful! You can close this window.")
		codeCh <- code
	})

	// Start a local server on the port expected by the redirect URI
	// Since credentials.json has http://localhost, it usually implies port 80 or a random port if configured.
	// However, Google allows http://localhost without a port for some app types,
	// but it's better to be explicit if we can.
	// Looking at the redirect_uris: ["http://localhost"]
	port := "8888"
	config.RedirectURL = "http://localhost:" + port
	server := &http.Server{Addr: "localhost:" + port, Handler: handler}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser: \n%v\n", authURL)

	var authCode string
	select {
	case authCode = <-codeCh:
		// Received the code!
	case err := <-errCh:
		log.Fatalf("Error during authentication: %v", err)
	case <-time.After(2 * time.Minute):
		log.Fatalf("Timeout waiting for authentication code")
	}

	// Shut down the server
	err := server.Close()
	if err != nil {
		log.Fatalf("Error shutting down the server: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		errClose := f.Close()
		if errClose != nil {
			log.Fatalf("Error closing token file: %v", errClose)
		}
	}()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer func() {
		errClose := f.Close()
		if errClose != nil {
			log.Fatalf("Error closing token file: %v", errClose)
		}
	}()
	if err := json.NewEncoder(f).Encode(token); err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
}
