package main

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
)

type JWTToken struct {
	Token string `json:"token"`
}

type CA struct {
	Cert string `json:"cert"`
}

func getJWTToken(access_token string, issuerUrl string) (string, error) {
	var jwt JWTToken

	resp, _, err := gorequest.New().Get(issuerUrl).
		Set("Authorization", "Bearer "+access_token).
		EndStruct(&jwt)

	if err != nil {
		log.Warn("Failed in fetching JWT Token ", err)
		return "", err[0]
	}

	if resp != nil && resp.StatusCode != 201 {
		log.Warn("Failed in fetching JWT Token, responsecode: ", resp.StatusCode)
		return "", errors.New("Failed in fetching JWT Token")
	}

	return jwt.Token, nil
}

func getCACert(issuerUrl string) ([]byte, error) {
	var ca CA

	resp, _, err := gorequest.New().Get(issuerUrl + "/ca").
		EndStruct(&ca)

	if err != nil {
		log.Warn("Failed in fetching CA certificate ", err)
		return nil, err[0]
	}

	if resp != nil && resp.StatusCode != 200 {
		log.Warn("Failed in fetching CA certificate, responsecode: ", resp.StatusCode)
		return nil, errors.New("Failed in fetching CA certificate")
	}
	return []byte(ca.Cert), nil
}
