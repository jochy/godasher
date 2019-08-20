package godasher

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(target string, insecure bool) (*Proxy, error) {
	targetUrl, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	pxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy := &Proxy{
		target: targetUrl,
		proxy:  pxy,
	}

	proxy.proxy.ModifyResponse= func(response *http.Response) error {
		// Allow IFrame even if target doesn't want to
		response.Header.Del("X-Frame-Options")
		return nil
	}

	if insecure {
		proxy.proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return proxy, nil
}

type Proxy struct {
	target  *url.URL
	proxy   *httputil.ReverseProxy
}

func (p *Proxy) handle(w http.ResponseWriter, r *http.Request) {
	r.Host = p.target.Host
	p.proxy.ServeHTTP(w, r)
}

func (p *Proxy) StartServer(port string) {
	// GO routine with a new proxy-server
	go func() {
		proxyserver := http.NewServeMux()
		proxyserver.HandleFunc("/", p.handle)
		log.Printf("%v", http.ListenAndServe(":"+port, proxyserver))
	}()
}