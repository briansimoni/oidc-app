/*
This is an example application to demonstrate parsing an ID Token.
*/
package main

import (
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tmpl = template.Must(template.ParseFiles("index.html"))

type UserInfo struct {
	Sub string `json:"sub"`
}

// These are required
var (
	tokenURL        = os.Getenv("TOKEN_URL")
	authURL         = os.Getenv("AUTH_URL")
	clientID        = os.Getenv("CLIENT_ID")
	clientSecret    = os.Getenv("CLIENT_SECRET")
	redirectURL     = os.Getenv("REDIRECT_URI")
	codeExchangeURL = os.Getenv("CODE_EXCHANGE_URL")
	port            = os.Getenv("PORT")
	userinfoURL     = os.Getenv("USER_INFO_URL")
)

func main() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/oidc-app/", home)
	http.HandleFunc("/oidc-app/redirect", callback)

	if port == "" {
		port = "8080"
	}
	log.Printf("listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //stupid
}
