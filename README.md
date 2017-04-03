# kubed
Get User JWT token from Dataporten enabled Token Issuer to be used for communication with Kubernetes cluster.

It also stores the obtsined token along with configuration under $HOME/.kube/config file. If the file already present it will merge the configuration for the obtained cluster with already present ones.