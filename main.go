package main

import (
	"flag"
	"fmt"
	"strings"

	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bclicn/color"
	"github.com/pkg/browser"
)

const authURL = "https://auth.dataporten.no/oauth/authorization"
const kubedConf = ".kubedconf"

var (
	kubeconfig  = flag.String("kube-config", "~/.kube/config", "Absolute path to the kubeconfig config to manage settings")
	apiserver   = flag.String("api-server", "", "Address of Kubernetes API server (Required)")
	issuerUrl   = flag.String("issuer", "", "Address of JWT Token Issuer (Required)")
	clusterName = flag.String("name", "", "Name of this Kubernetes cluster, used for context as well (Required)")
	showVersion = flag.Bool("version", false, "Prints version information and exits")
	keepContext = flag.Bool("keep-context", false, "Keep the current context or switch to newly created one")
	port        = flag.Int("port", 49999, "Port number where Oauth2 Provider will redirect Kubed")
	renew       = flag.String("renew", "", "Name of the cluster to renew JWT token for")
	client_id   = flag.String("client-id", "", "Client ID for Kubed app (Required)")
	version     = "none"
	reqErr      error
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

	if len(os.Args) < 3 {
		fmt.Println(color.Red("Please provide parameters to run Kubed, refer " + color.BRed("kubed -h")))
		os.Exit(1)
	}

	var cluster *Cluster
	var err error
	if *renew != "" {
		cluster, err = readConfig(*renew)
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

		// Check if we have all the required parameters
		if cluster.Name == "" || cluster.IssuerUrl == "" || cluster.APIServer == "" || cluster.ClientID == "" {
			fmt.Println(color.Red("Please provide all the required parameter, refer " + color.BRed("kubed -h")))
			os.Exit(1)
		}

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

	// Open browser to authenticate user and get access token
	browser.OpenURL(authURL + "?response_type=token&client_id=" + cluster.ClientID)
	token, err := getToken(cluster.Port)
	if err != nil {
		log.Fatal("Error in getting access token", err)
	}
	//wg.Wait() // Wait until we get the token back
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
	fmt.Println(color.Green("To renew JWT token for this cluster run: " + color.BGreen("kubed -renew "+cluster.Name)))
}
