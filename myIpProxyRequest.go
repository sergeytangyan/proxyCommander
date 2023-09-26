package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type MyIp struct {
	Ip string `json:"ip"`
}

type myIpProxyRequest struct {
	url     *url.URL
	timeout time.Duration
}

func newMyIpProxyRequest() *myIpProxyRequest {
	url, err := url.Parse("https://api64.ipify.org?format=json")
	if err != nil {
		log.Fatal(err)
	}

	return &myIpProxyRequest{
		url:     url,
		timeout: 20 * time.Second,
	}
}

func (pr *myIpProxyRequest) Request(pm *ProxyManager) error {
	// -- setup proxy ---------------------------------------------
	err := pm.RotateProxy()
	if err != nil {
		return err
	}

	// -- actual request ---------------------------------------------
	res, err := pm.GetClient().Do(&http.Request{
		Method: http.MethodPost,
		URL:    pr.url,
	})
	if err != nil {
		return err
	}

	// -- validate ---------------------------------------------
	if res.StatusCode >= 400 {
		return fmt.Errorf("non 200 status code returned: '%d'", res.StatusCode)
	}

	// -- parse response ---------------------------------------------
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	myip := &MyIp{}
	err = json.Unmarshal(body, myip)
	if err != nil {
		return err
	}

	// -- the end ---------------------------------------------
	fmt.Println(myip)
	return nil
}
