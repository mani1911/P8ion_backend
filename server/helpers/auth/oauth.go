package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"p8ion/config"
	"time"
)

type OAuth2Token struct {
	AccessToken string
	IDToken     string
}

type OAuth2User struct {
	Email         string
	VerifiedEmail bool
	Name          string
}

func GetGoogleOAuthTokens(code string) (*OAuth2Token, error) {
	tokenEndpoint := "https://oauth2.googleapis.com/token"
	Auth := config.GetConfig().Auth

	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("client_id", Auth.OAuth2Key)
	values.Add("client_secret", Auth.OAuth2Secret)
	values.Add("redirect_uri", Auth.RedirectURL)
	query := values.Encode()

	req, err := http.NewRequest("POST", tokenEndpoint, bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Could Not Retrieve Token")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var OAuthTokenRes map[string]interface{}
	if err := json.Unmarshal(resBody, &OAuthTokenRes); err != nil {
		return nil, err
	}

	tokenBody := &OAuth2Token{
		AccessToken: OAuthTokenRes["access_token"].(string),
		IDToken:     OAuthTokenRes["id_token"].(string),
	}

	return tokenBody, nil
}

func GetOAuth2User(accessToken string, idToken string) (*OAuth2User, error) {
	userEndpoint := "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"

	req, err := http.NewRequest("GET", userEndpoint+"&access_token="+accessToken, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", idToken))

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var GoogleUserRes map[string]interface{}

	if err := json.Unmarshal(resBody, &GoogleUserRes); err != nil {
		return nil, err
	}

	userBody := &OAuth2User{
		Email:         GoogleUserRes["email"].(string),
		VerifiedEmail: GoogleUserRes["verified_email"].(bool),
		Name:          GoogleUserRes["name"].(string),
	}

	return userBody, nil
}
