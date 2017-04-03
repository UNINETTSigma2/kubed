package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	yaml "gopkg.in/yaml.v2"
)

type Cluster struct {
	Name        string `yaml:"name"`
	APIServer   string `yaml:"apiserver"`
	IssuerUrl   string `yaml:"issuer"`
	ClientID    string `yaml:"clientid"`
	KubeConfig  string `yaml:kubeconfig"`
	KeepContext bool   `yaml:"keepcontext"`
	Port        int    `yaml:"port"`
}

func readConfig(name string) (*Cluster, error) {
	home := os.Getenv("HOME")
	path := filepath.Join(home, kubedConf)
	confBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Warn("Failed in reading kubed config file ", err)
		return nil, err
	}

	var clusters []Cluster
	err = yaml.Unmarshal(confBytes, &clusters)
	if err != nil {
		log.Error("Failed in parsing config file ", err)
	}

	for _, c := range clusters {
		if c.Name == name {
			return &c, nil
		}
	}

	spew.Dump(clusters)
	return nil, errors.New("Provided cluster not configured, run with full config parameters to configure it")
}

func setConfig(
	name string,
	apiserver string,
	issuerUrl string,
	client_id string,
	kubeconfig string,
	keepContext bool,
	port int) *Cluster {

	return &Cluster{
		Name:        name,
		APIServer:   apiserver,
		IssuerUrl:   issuerUrl,
		ClientID:    client_id,
		KubeConfig:  kubeconfig,
		KeepContext: keepContext,
		Port:        port,
	}
}

func saveConfig(cluster *Cluster) error {
	home := os.Getenv("HOME")
	path := filepath.Join(home, kubedConf)

	var clusters []Cluster
	clusters = append(clusters, *cluster)
	confBytes, err := yaml.Marshal(clusters)
	if err != nil {
		log.Warn("Failed in marshaling kubedconfig ", err)
		return err
	}

	err = ioutil.WriteFile(path, confBytes, 0644)
	if err != nil {
		log.Warn("Failed in saving kubedconfig ", err)
		return err
	}

	return nil
}
