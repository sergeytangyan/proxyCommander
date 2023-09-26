package main

import (
	"log"
	"time"

	"sergeytangyan/proxyCommander/providers"
)

func main() {
	myIpProxyRequest := newMyIpProxyRequest()

	provider := &providers.SslProxiesOrgProvider{}
	pm := newProxyManager(provider, myIpProxyRequest.timeout)

	for {
		err := pm.Do(myIpProxyRequest)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(1500 * time.Millisecond)
	}
}
