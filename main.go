package main

import (
	"flag"
	"fmt"

	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/browser"
)

const authURL = "https://auth.dataporten.no/oauth/authorization"

var (
	kubeconfig  = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
	apiserver   = flag.String("apiserver", "https://localhost", "Address of Kubernetes API server")
	showVersion = flag.Bool("version", false, "Prints version information and exits")
	port        = flag.Int("port", 49999, "Port number where Oauth2 Provider will redirect Kubed")
	client_id   = flag.String("client_id", "3181c169-9fbe-4f31-8802-06e45cab9b00", "Client ID for Kubed app")
	version     = "none"
	token       string
	reqErr      error
	wg          sync.WaitGroup
)

func init() {
	// Log as JSON to stderr
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)

	flag.Parse()
	if *showVersion {
		fmt.Println("kubed version", version)
		os.Exit(0)
	}
}

func main() {

	// Open brower to authenticate user and get access token
	browser.OpenURL(authURL + "?response_type=token&scope=gk_jwt+userid&client_id=" + *client_id)
	if err := getToken(*port); err != nil {
		log.Fatal("Error in getting access token", err)
	}
	wg.Wait() // Wait until we get the token back
	if reqErr != nil {
		log.Fatal("Error in getting access token", reqErr)
		os.Exit(1)
	}

	spew.Dump(token)
}
