[![Build Status](https://travis-ci.org/UNINETT/kubed.png)](https://travis-ci.org/UNINETT/kubed)

# Kubed (Kubernetes + Dataporten)
Get a JWT token from Dataporten enabled Token Issuer to be used for communication with Kubernetes cluster. This utility configures `kubectl` configuration with information about Kubernertes API server and authentication details. If the file already present, it will merge the configuration for the obtained cluster with already present ones. Example run is

```
kubed -name test-cluster -api-server https://kubernetes.apiserver.com -client-id client-id-from-your-cluster -issuer https://token.issuer.com
```

After successful authentication, kubed will store the credentials in `$HOME/.kube/config` file, by default. You can specify `kubectl config` file with parameter `-kube-config`. Now you can run your favourite `kubectl` commands against `https://kubernetes.apiserver.com`.

Kubed will also store this cluster configuration, so for JWT token renewal, you can simply run the command
```
kubed -renew test-cluster
```
