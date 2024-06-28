package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

const tokenFilename = "token.json"

var oauthConfig = oauth2.Config{
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize",
		TokenURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/token",
	},
	RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
	ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
	Scopes: []string{
		"https://graph.microsoft.com/Files.ReadWrite",
		"offline_access",
	},
}

func getStoredOrNewOAuthToken(ctx context.Context) (*oauth2.Token, error) {
	tokenJson, err := os.ReadFile(tokenFilename)
	if err == nil {
		var token oauth2.Token
		err = json.Unmarshal(tokenJson, &token)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal token file content: %w", err)
		}
		return &token, err
	}

	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("read token file: %w", err)
	}

	token, err := getNewOAuthToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get oauth token: %w", err)
	}

	err = storeToken(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func storeToken(token *oauth2.Token) error {
	tokenJson, err := json.Marshal(*token)
	if err != nil {
		return fmt.Errorf("json marshal token: %w", err)
	}

	err = os.WriteFile(tokenFilename, tokenJson, 0644)
	if err != nil {
		return fmt.Errorf("write token file: %w", err)
	}

	return nil
}

func getNewOAuthToken(ctx context.Context) (*oauth2.Token, error) {
	tokens, errors, httpServer := startTemporaryOAuthCallbackServer()

	defer func() {
		_ = httpServer.Shutdown(ctx)
		close(tokens)
		close(errors)
	}()

	select {
	case err := <-errors:
		return nil, err
	case token := <-tokens:
		return token, nil
	}
}

func startTemporaryOAuthCallbackServer() (chan *oauth2.Token, chan error, *http.Server) {
	tokenChannel := make(chan *oauth2.Token)
	errorChannel := make(chan error)

	state := randomString(32)

	loginUrl := oauthConfig.AuthCodeURL(state)
	fmt.Printf("Please log in:\n%v\n", loginUrl)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /auth/callback", func(w http.ResponseWriter, r *http.Request) {
		handleOAuthCallback(w, r, state, tokenChannel, errorChannel)
	})

	srv := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errorChannel <- fmt.Errorf("http listen and serve: %w", err)
		}
	}()

	return tokenChannel, errorChannel, srv
}

func handleOAuthCallback(w http.ResponseWriter, r *http.Request, state string, tokens chan *oauth2.Token, errors chan error) {
	errorCode := r.FormValue("error")
	if errorCode != "" {
		errorDescription := r.FormValue("error_description")
		errorText := fmt.Sprintf("Error: %v, Description: %v", errorCode, errorDescription)

		http.Error(w, errorText, http.StatusBadRequest)
		errors <- fmt.Errorf(errorText)
		return
	}

	if r.FormValue("state") != state {
		http.Error(w, "invalid state", http.StatusBadRequest)
		errors <- fmt.Errorf("state from callback request did not match")
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		errors <- fmt.Errorf("code not found in callback request")
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code, oauth2.AccessTypeOffline)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		errors <- fmt.Errorf("code exchange: %w", err)
		return
	}

	fmt.Println("OAuth2 token successfully received")
	_, _ = fmt.Fprintf(w, "Authentication successfull. You can close this page.")
	tokens <- token
}

func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
