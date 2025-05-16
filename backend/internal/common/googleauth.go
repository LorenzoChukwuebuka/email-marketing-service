package common

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleUserData represents the user information returned from Google
type GoogleUserData struct {
	ID            string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	VerifiedEmail bool   `json:"verified_email"`
}

// GoogleAuthConfig holds the configuration for Google OAuth
type GoogleAuthConfig struct {
	clientID     string
	clientSecret string
	redirectURL  string
	scopes       []string
	config       *oauth2.Config
}

// GoogleAuthOption defines the type for option functions
type GoogleAuthOption func(*GoogleAuthConfig)

// WithClientID sets the client ID
func WithClientID(clientID string) GoogleAuthOption {
	return func(c *GoogleAuthConfig) {
		c.clientID = clientID
	}
}

// WithClientSecret sets the client secret
func WithClientSecret(clientSecret string) GoogleAuthOption {
	return func(c *GoogleAuthConfig) {
		c.clientSecret = clientSecret
	}
}

// WithRedirectURL sets the redirect URL
func WithRedirectURL(redirectURL string) GoogleAuthOption {
	return func(c *GoogleAuthConfig) {
		c.redirectURL = redirectURL
	}
}

// WithScopes sets additional scopes
func WithScopes(scopes []string) GoogleAuthOption {
	return func(c *GoogleAuthConfig) {
		c.scopes = append(c.scopes, scopes...)
	}
}

// NewGoogleAuth creates a new GoogleAuthConfig with the provided options
func NewGoogleAuth(opts ...GoogleAuthOption) *GoogleAuthConfig {
	// Default configuration
	config := &GoogleAuthConfig{
		scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}

	// Apply all options
	for _, opt := range opts {
		opt(config)
	}

	// Initialize the OAuth2 config
	config.config = &oauth2.Config{
		ClientID:     config.clientID,
		ClientSecret: config.clientSecret,
		RedirectURL:  config.redirectURL,
		Scopes:       config.scopes,
		Endpoint:     google.Endpoint,
	}

	return config
}

// GetAuthURL generates the Google OAuth2 URL with the provided state
func (g *GoogleAuthConfig) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// GetUserData exchanges the authorization code for user data
func (g *GoogleAuthConfig) GetUserData(ctx context.Context, code string) (*GoogleUserData, error) {
	// Exchange the authorization code for a token
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Create an HTTP client with the token
	client := g.config.Client(ctx, token)

	// Get user info from Google
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	//   // Add this debug step to see raw response
	//   body, err := io.ReadAll(resp.Body)
	//   if err != nil {
	//       return nil, fmt.Errorf("failed to read response body: %w", err)
	//   }

	//   // Print raw response for debugging
	//   fmt.Printf("Raw response: %s\n", string(body))

	// Decode the user data
	var userData GoogleUserData
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return nil, fmt.Errorf("failed to decode user data: %w", err)
	}

	fmt.Printf("%+v", userData)
	return &userData, nil
}
