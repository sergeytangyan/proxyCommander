package sources

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type SslProxiesOrgSource struct {
	lastProxy string
}

func (p *SslProxiesOrgSource) GetProxyList() ([]string, error) {
	log.Println("Getting proxies")

	res, err := http.Get("https://sslproxies.org/")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	proxies := []string{}

	doc.Find(".fpl-list tbody tr").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		ip := s.Find("td:nth-child(1)").Text()
		port := s.Find("td:nth-child(2)").Text()
		proxy := fmt.Sprintf("http://%s:%s", ip, port)

		if proxy == p.lastProxy {
			log.Println("last proxy match: ", p.lastProxy)
			return false
		}

		proxies = append(proxies, proxy)

		return true
	})

	p.lastProxy = proxies[0]

	log.Printf("Received %d proxies\n", len(proxies))

	return proxies, nil
}
