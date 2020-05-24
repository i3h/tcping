package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	API string = "https://ifconfig.is/json/"
)

type IPInfo struct {
	Continent string  `json:"Continent"`
	Country   string  `json:"Country"`
	City      string  `json:"City"`
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	TimeZone  string  `json:"TimeZone"`
	IsEU      bool    `json:"IsEU"`
	ASN       uint    `json:"ASN"`
	ORG       string  `json:"ORG"`
}

func queryInfo(address string) IPInfo {
	var info IPInfo

	res, err := http.Get(API + address)
	if err != nil {
		match, _ := regexp.MatchString("connection reset by peer", err.Error())
		if match {
			log.Fatal("Oops, your connection was reset by magic power. You may need to set env http_proxy.")
		} else {
			log.Fatal(err)
		}
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &info)
	if err != nil {
		fmt.Println(err)
	}
	//printInfo(info)
	return info
}
