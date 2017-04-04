package main

import (
	"flag"
	"fmt"
	"strings"

	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/bclicn/color"
	"github.com/pkg/browser"
)

const authURL = "https://auth.dataporten.no/oauth/authorization"
const kubedConf = ".kubedconf"

var (
	kubeconfig  = flag.String("kubeconfig", "~/.kube/config", "Absolute path to the kubeconfig config to manage settings")
	apiserver   = flag.String("apiserver", "https://localhost", "Address of Kubernetes API server")
	issuerUrl   = flag.String("issuer", "https://token.example.no", "Address of JWT Token Issuer")
	clusterName = flag.String("name", "test", "Name of this Kubernetes cluster, used for context as well")
	showVersion = flag.Bool("version", false, "Prints version information and exits")
	keepContext = flag.Bool("keep-context", false, "Keep the current context or switch to newly created one")
	port        = flag.Int("port", 49999, "Port number where Oauth2 Provider will redirect Kubed")
	client_id   = flag.String("client-id", "daa8f3c8-422f-40b5-a045-06e86b987557", "Client ID for Kubed app")
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

	var cluster *Cluster
	var err error
	if len(os.Args) == 3 && os.Args[1] == "-name" {
		cluster, err = readConfig(*clusterName)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		cluster = setConfig(
			*clusterName,
			*apiserver,
			*issuerUrl,
			*client_id,
			*kubeconfig,
			*keepContext,
			*port)

		// Save the current cluster config, so we can reuse it during token renewal
		err = saveConfig(cluster)
		if err != nil {
			log.Fatal("Failed in saving kubedconfig ", err)
		}
	}

	// Fix Home Path for Kubeconfig
	if strings.HasPrefix(cluster.KubeConfig, "~") {
		home := os.Getenv("HOME")
		cluster.KubeConfig = strings.Replace(cluster.KubeConfig, "~", home, 1)
	}

	// Open brower to authenticate user and get access token
	browser.OpenURL(authURL + "?response_type=token&client_id=" + cluster.ClientID)
	if err = getToken(cluster.Port); err != nil {
		log.Fatal("Error in getting access token", err)
	}
	wg.Wait() // Wait until we get the token back
	if reqErr != nil {
		log.Fatal("Error in getting access token ", reqErr)
		os.Exit(1)
	}

	log.Info("Requesting JWT Token from ", cluster.IssuerUrl)

	cfg := new(KubeConfigSetup)
	cfg.Token, err = getJWTToken(token, cluster.IssuerUrl)
	if err != nil {
		log.Fatal("Failed in getting JWT token ", err)
		os.Exit(1)
	}
	cfg.CertificateAuthorityData, err = getCACert(cluster.IssuerUrl)
	if err != nil {
		log.Warn("No custom CA certificate provided, assuming running with standard certificate")
	}

	cfg.ClusterName = cluster.Name
	cfg.ClusterServerAddress = cluster.APIServer
	cfg.kubeConfigFile = cluster.KubeConfig
	cfg.KeepContext = cluster.KeepContext

	err = SetupKubeConfig(cfg)
	if err != nil {
		log.Fatal("Failed in setting the kubeconfig ", err)
	}

	log.Info("Kubernetes configuration has been saved in ", cluster.KubeConfig, " with context ", cluster.Name)
	fmt.Println(color.Green("To renew JWT token for this cluster run: " + color.BGreen("kubed -name "+cluster.Name)))
}
