package main

// Reference Implementation from Minikube

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/clientcmd/api/latest"
)

type KubeConfigSetup struct {
	// The name of the cluster for this context
	ClusterName string

	// ClusterServerAddress is the address of of the kubernetes cluster
	ClusterServerAddress string

	// CertificateAuthority is the path to a cert file for the certificate authority.
	CertificateAuthorityData []byte

	// ClientToken is the path to a client key file for TLS.
	Token string

	// Should the current context be kept when setting up this one
	KeepContext bool

	// kubeConfigFile is the path where the kube config is stored
	kubeConfigFile string

	// NameSpace is the default namespace used with kubectl. May be blank.
	NameSpace string
}

// SetupKubeconfig reads config from disk, adds the minikube settings, and writes it back.
// activeContext is true when minikube is the CurrentContext
// If no CurrentContext is set, the given name will be used.
func SetupKubeConfig(cfg *KubeConfigSetup) error {
	// read existing config or create new if does not exist
	config, err := ReadConfigOrNew(cfg.kubeConfigFile)
	if err != nil {
		return err
	}

	clusterName := cfg.ClusterName
	cluster := api.NewCluster()
	cluster.Server = cfg.ClusterServerAddress
	cluster.CertificateAuthorityData = cfg.CertificateAuthorityData
	config.Clusters[clusterName] = cluster

	// user
	userName := cfg.ClusterName
	user := api.NewAuthInfo()
	user.Token = cfg.Token
	config.AuthInfos[userName] = user

	// context
	contextName := cfg.ClusterName
	context := api.NewContext()
	context.Cluster = cfg.ClusterName
	context.AuthInfo = userName
	if (cfg.NameSpace != "") {
		context.Namespace = cfg.NameSpace
	}
	config.Contexts[contextName] = context

	// Only set current context to minikube if the user has not used the keepContext flag
	if !cfg.KeepContext {
		config.CurrentContext = contextName
	}

	// write back to disk
	if err := WriteConfig(config, cfg.kubeConfigFile); err != nil {
		return err
	}
	return nil
}

// ReadConfigOrNew retrieves Kubernetes client configuration from a file.
// If no files exists, an empty configuration is returned.
func ReadConfigOrNew(filename string) (*api.Config, error) {
	data, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		return api.NewConfig(), nil
	} else if err != nil {
		return nil, errors.Wrapf(err, "Error reading file %q", filename)
	}

	// decode config, empty if no bytes
	config, err := decode(data)
	if err != nil {
		return nil, errors.Errorf("could not read config: %v", err)
	}

	// initialize nil maps
	if config.AuthInfos == nil {
		config.AuthInfos = map[string]*api.AuthInfo{}
	}
	if config.Clusters == nil {
		config.Clusters = map[string]*api.Cluster{}
	}
	if config.Contexts == nil {
		config.Contexts = map[string]*api.Context{}
	}

	return config, nil
}

// WriteConfig encodes the configuration and writes it to the given file.
// If the file exists, it's contents will be overwritten.
func WriteConfig(config *api.Config, filename string) error {
	if config == nil {
		log.Error("could not write to '%s': config can't be nil", filename)
	}

	// encode config to YAML
	data, err := runtime.Encode(latest.Codec, config)
	if err != nil {
		return errors.Errorf("could not write to '%s': failed to encode config: %v", filename, err)
	}

	// create parent dir if doesn't exist
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return errors.Wrapf(err, "Error creating directory: %s", dir)
		}
	}

	// write with restricted permissions
	if err := ioutil.WriteFile(filename, data, 0600); err != nil {
		return errors.Wrapf(err, "Error writing file %s", filename)
	}
	return nil
}

// decode reads a Config object from bytes.
// Returns empty config if no bytes.
func decode(data []byte) (*api.Config, error) {
	// if no data, return empty config
	if len(data) == 0 {
		return api.NewConfig(), nil
	}

	config, _, err := latest.Codec.Decode(data, nil, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Error decoding config from data: %s", string(data))
	}

	return config.(*api.Config), nil
}
