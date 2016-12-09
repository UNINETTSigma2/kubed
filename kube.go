package main

import (
// 	"fmt"
// 	log "github.com/Sirupsen/logrus"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/pkg/api/v1"
// 	"k8s.io/client-go/tools/clientcmd"
)

func getAppInfo(secret string) string {
	return `{"name": "test-kube", "redirect_uri": ["http://localhost:9302"],
			"scopes_requested": ["userid", "profile"], "authproviders": ["feide|all"],
			"client_secret": "` + secret + "\"}"
}

// func connectKube(kubeconfig string) {
// 	// uses the current context in kubeconfig
// 	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
// 	if err != nil {
// 		log.Fatal("Error in getting kubeconfig ", err.Error())
// 	}
// 	// creates the clientset
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		log.Fatal("Error in creating client ", err.Error())
// 	}
// }
