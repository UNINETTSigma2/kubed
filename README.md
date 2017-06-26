[![Build Status](https://travis-ci.org/UNINETT/kubed.png)](https://travis-ci.org/UNINETT/kubed)

# Kubed (Kubernetes + Dataporten)
This utility manages `kubectl` configuration with information about Kubernertes API server and authentication details. Kubed gets a JWT token from Dataporten enabled token issuer to be used for communication with Kubernetes cluster. If the file already present, it will merge the configuration for the obtained cluster with already present ones. Example run is

```
kubed -name test-cluster -api-server https://kubernetes.apiserver.com -client-id client-id-from-your-cluster -issuer https://token.issuer.com
```

After successful authentication, kubed will store the credentials in `$HOME/.kube/config` file, by default. You can specify `kubectl config` file with parameter `-kube-config`. Now you can run your favourite `kubectl` commands against `https://kubernetes.apiserver.com`.

Kubed will also store this cluster configuration, so for JWT token renewal, you can simply run the command
```
kubed -renew test-cluster
```

## Installation
To instal, run the following commands based on your operating system

**For MAC OSX (amd64)**
```
curl -LO https://github.com/UNINETT/kubed/releases/download/0.1.7/kubed-darwin-amd64 && sudo mv kubed-darwin-amd64 /usr/local/bin/kubed && sudo chmod +x /usr/local/bin/kubed
```

**For Linux (amd64)**
```
curl -LO https://github.com/UNINETT/kubed/releases/download/0.1.7/kubed-linux-amd64 && sudo mv kubed-linux-amd64 /usr/local/bin/kubed && sudo chmod +x /usr/local/bin/kubed
```

**For Windows (amd64)**

Download <a href="https://github.com/UNINETT/kubed/releases/download/0.1.7/kubed-windows-amd64.exe" target="_blank">Kubed</a> and then open <b>cmd</b>. On the command prompt run the following command
```
copy %HOMEPATH%\Downloads\kubed-windows-amd64.exe C:\Windows\System32\kubed.exe
```
