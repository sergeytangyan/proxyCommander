package main

import (
	"log"
	"time"

	"sergeytangyan/proxyCommander/sources"
)

func main() {
	proxiesSource := &sources.SslProxiesOrgSource{}
	myIpProxyRequest := newMyIpProxyRequest()

	pm := newProxyManager(proxiesSource, myIpProxyRequest.timeout)

	for {
		err := pm.Do(myIpProxyRequest)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(1500 * time.Millisecond)
	}
}
