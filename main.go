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
	kubeconfig  = flag.String("kubeconfig", "./config", "Path to the kubeconfig config to manage settings")
	apiserver   = flag.String("apiserver", "https://localhost", "Address of Kubernetes API server")
	issuerUrl   = flag.String("issuer-url", "https://token.example.no", "Address of JWT Token Issuer")
	issuerScope = flag.String("issuer-scope", "gk_jwt", "Scope name of JWT Token Issuer")
	clusterName = flag.String("name", "test", "Name of this Kubernetes cluster, used for context as well")
	showVersion = flag.Bool("version", false, "Prints version information and exits")
	port        = flag.Int("port", 49999, "Port number where Oauth2 Provider will redirect Kubed")
	client_id   = flag.String("client_id", "daa8f3c8-422f-40b5-a045-06e86b987557", "Client ID for Kubed app")
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
	browser.OpenURL(authURL + "?response_type=token&scope=userid " + *issuerScope + "&client_id=" + *client_id)
	if err := getToken(*port); err != nil {
		log.Fatal("Error in getting access token", err)
	}
	wg.Wait() // Wait until we get the token back
	if reqErr != nil {
		log.Fatal("Error in getting access token ", reqErr)
		os.Exit(1)
	}

	jwt, err := getJWTToken(token, *issuerUrl)
	if err != nil {
		log.Fatal("Failed in getting JWT token ", err)
		os.Exit(1)
	}
	cert, err := getCACert(*issuerUrl)
	if err != nil {
		log.Warn("No custom CA certificate provided, assuming running with standard certificate")
	}
	spew.Dump(jwt)
	spew.Dump(cert)

}
