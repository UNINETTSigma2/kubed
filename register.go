package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
	"time"
)

func registerApp(token string, secret string) bool {
	if checkRegisteredApp(token, "test-kube") {
		log.Warn("Application is already registered with same name and will not register again")
		return false
	}

	resp, body, errs := gorequest.New().
		Timeout(20*time.Second).
		Post(regURL).
		Set("Authorization", "Bearer "+token).
		Send(getAppInfo(secret)).
		End()

	if errs != nil {
		log.Fatal("Failed in registering Application ", errs)
		return false
	}

	if resp.StatusCode != 201 {
		log.Fatal("Failed in registering Application, got code: ", resp.StatusCode, body)
		return false
	}

	if err := json.Unmarshal([]byte(body), appInfo); err != nil {
		log.Error("Could not map Application Info object", err)
	}
	return true
}

func checkRegisteredApp(token string, appName string) bool {
	appList := new([]AppInfo)

	resp, _, errs := gorequest.New().
		Timeout(20*time.Second).
		Get(regURL).
		Set("Authorization", "Bearer "+token).
		EndStruct(appList)

	if errs != nil || resp.StatusCode != 200 {
		log.Error("Failed in getting Applications list ", errs)
		return false
	}

	for _, app := range *appList {
		if appName == app.Name {
			appInfo = &app
			return true
		}
	}
	return false
}
