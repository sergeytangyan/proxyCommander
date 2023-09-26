package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

type ProxyProvider interface {
	GetNextBatch() ([]string, error)
}

type ProxyRequest interface {
	Request(pm *ProxyManager) error
}

type ProxyManager struct {
	proxyList []string
	provider  ProxyProvider
	client    *http.Client
}

func newProxyManager(provider ProxyProvider, timeout time.Duration) *ProxyManager {
	pm := &ProxyManager{
		provider: provider,
		client: &http.Client{
			Timeout: timeout,
		},
	}

	return pm
}

func (pm *ProxyManager) getNextProxy() (string, error) {
	length := len(pm.proxyList)

	if length == 0 {
		res, err := pm.provider.GetNextBatch()
		if err != nil {
			return "", err
		}

		pm.proxyList = res
		length = len(pm.proxyList)
	}

	lastElement := pm.proxyList[length-1]
	pm.proxyList = pm.proxyList[:length-1]

	return lastElement, nil
}

func (pm *ProxyManager) RotateProxy() error {
	proxy, err := pm.getNextProxy()
	if err != nil {
		return err
	}

	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return err
	}
	log.Printf("Using proxy '%s'\n", proxy)

	pm.client.CloseIdleConnections()
	pm.client.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	return nil
}

func (pm *ProxyManager) GetClient() *http.Client {
	return pm.client
}

func (pm *ProxyManager) Do(pr ProxyRequest) error {
	return pr.Request(pm)
}
