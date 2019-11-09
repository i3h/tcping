package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	API string = "http://ip-api.com/json/"
)

type IPInfo struct {
	IP      string `json:"query"`
	City    string `json:"city"`
	Country string `json:"country"`
	ISP     string `json:"isp"`
	AS      string `json:"as"`
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
