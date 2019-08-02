package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	jwt "gopkg.in/square/go-jose.v2"
)

var state = "foobar"

// oidcResponse is the raw response from exchanging an authorization code
type oidcResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func home(w http.ResponseWriter, r *http.Request) {
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     oauth2.Endpoint{TokenURL: tokenURL, AuthURL: authURL},
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "profile"},
	}

	authCodeURL := config.AuthCodeURL(state)
	http.Redirect(w, r, authCodeURL, http.StatusFound)
}

func callback(w http.ResponseWriter, r *http.Request) {
	log.Println("performing OAuth code flow callback")
	for key, value := range r.Header {
		log.Println("HEADER:", key, value)
	}

	if r.URL.Query().Get("state") != state {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	tokens, err := exchange(code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.MarshalIndent(tokens, "", "    ")
	if err != nil {
		http.Error(w, "Failed to create JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	jsonString := string(j)

	idToken, err := jwt.ParseSigned(tokens.IDToken)
	if err != nil {
		http.Error(w, "Failed to create JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	prettyIDToken, err := prettyprint(idToken.UnsafePayloadWithoutVerification())
	if err != nil {
		http.Error(w, "Failed to pretty print ID Token "+err.Error(), http.StatusInternalServerError)
		return
	}

	var userInfo []byte
	if userinfoURL != "" {
		info, err := getUserInfo(userinfoURL, tokens.AccessToken)
		if err != nil {
			http.Error(w, "Failed to obtain userinfo "+err.Error(), http.StatusInternalServerError)
			return
		}
		prettyInfo, err := prettyprint(info)
		if err != nil {
			http.Error(w, "Failed to pretty print user info "+err.Error(), http.StatusInternalServerError)
			return
		}
		userInfo = prettyInfo
	}

	templateData := struct {
		Response string
		IDToken  string
		UserInfo string
	}{
		Response: jsonString,
		IDToken:  string(prettyIDToken),
		UserInfo: string(userInfo),
	}

	tmpl.Execute(w, templateData)
}

// exchange code for id_token and access_token
func exchange(code string) (*oidcResponse, error) {
	endpoint := codeExchangeURL
	data := url.Values{}
	data["grant_type"] = []string{"authorization_code"}
	data["redirect_uri"] = []string{redirectURL}
	data["code"] = []string{code}
	data["client_id"] = []string{clientID}
	data["client_secret"] = []string{clientSecret}

	res, err := http.PostForm(endpoint, data)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New("Error, status was " + res.Status)
	}
	defer res.Body.Close()
	var raw oidcResponse
	err = json.NewDecoder(res.Body).Decode(&raw)
	if err != nil {
		return nil, err
	}
	return &raw, nil
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func getUserInfo(url string, accessToken string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	log.Println(req.Header.Get("Authorization"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println("userinfo response", res.StatusCode, res.Status)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Println(string(body))
	}
	return body, nil
}
