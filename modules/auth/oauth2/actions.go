package oauth2

import (
	"context"
	"io"
	"net/http"
)

func Login() {
	url := GoogleConfig().GoogleLoginConfig.AuthCodeURL("randomstate") // setting a randomstate to prevent CSRF attack
	http.RedirectHandler(url, http.StatusTemporaryRedirect)
}

func GoogleCallback(state, code string) string {
	// The code will be in the *http.Request.FormValue("code")
	// Try fetching it from request not as a funtion arg
	token, err := GoogleConfig().GoogleLoginConfig.Exchange(context.Background(), code)
	if err != nil {
		panic("code-token exchange failed.")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		panic("User Data Fetch Failed")
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("JSON Parsing Failed")
	}

	return string(userData)
}
