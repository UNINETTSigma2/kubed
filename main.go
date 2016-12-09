package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/m4rw3r/uuid"
	"github.com/pkg/browser"
	"os"
	"sync"
)

type AppInfo struct {
	Name         string   `json:"name"`
	ClientID     string   `json:"id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
	RedirectURI  []string `json:"redirect_uri"`
}

const authURL = "https://auth.dataporten.no/oauth/authorization"
const regURL = "https://clientadmin.dataporten-api.no/clients/"

var (
	kubeconfig  = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
	deployment  = flag.String("dep", "./dep.yaml", "absolute path to the Deployment file")
	ingress     = flag.String("ing", "./ing.yaml", "absolute path to the Ingress file")
	showVersion = flag.Bool("version", false, "Prints version information and exits")
	port        = flag.Int("port", 49999, "Port number where Oauth2 Provider will redirect Kubed")
	client_id   = flag.String("client_id", "3181c169-9fbe-4f31-8802-06e45cab9b00", "Client ID for Kubed app")
	version     = "none"
	appInfo     = new(AppInfo)
	token       string
	reqErr      error
	wg          sync.WaitGroup
)

func init() {
	// Log as JSON to stderr
	log.SetFormatter(&log.JSONFormatter{"2006-01-02T15:04:05.000Z07:00"})
	log.SetOutput(os.Stderr)
}

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Println("kubed version", version)
		os.Exit(0)
	}

	// Open brower to authenticate user and get access token
	browser.OpenURL(authURL + "?response_type=token&client_id=" + *client_id)
	if err := getToken(*port); err != nil {
		log.Fatal("Error in getting access token", err)
	}
	wg.Wait() // Wait until we get the token back
	if reqErr != nil {
		log.Fatal("Error in getting access token", reqErr)
		os.Exit(1)
	}

	// Generate client_secret, as Dataporten wants it to be sent by registering application
	secret, err := uuid.V4()
	if err != nil {
		log.Fatal("Error in generating UUID", err)
	}

	if registerApp(token, secret.String()) {
		fmt.Println("Registered Application: ", getAppInfo(secret.String()))
	}
	spew.Dump(appInfo)
}
