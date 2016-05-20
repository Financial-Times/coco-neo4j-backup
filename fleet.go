package main

import (
	"github.com/coreos/fleet/client"
	"net/url"
	"net"
	"time"
	"golang.org/x/net/proxy"
	"net/http"
	"log"
)

func newFleetClient() (client.API, error) {
	u, err := url.Parse(*fleetEndpoint)
	if err != nil {
		panic(err) // TODO handle this properly
	}
	httpClient := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 100}}

	if *socksProxy != "" {
		log.Printf("using proxy %s\n", *socksProxy)
		netDialler := &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		dialer, err := proxy.SOCKS5("tcp", *socksProxy, nil, netDialler)
		if err != nil {
			log.Fatalf("error with proxy %s: %v\n", *socksProxy, err)
		}
		httpClient.Transport = &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			Dial:                dialer.Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			MaxIdleConnsPerHost: 100,
		}
	}

	fleetHTTPAPIClient, err := client.NewHTTPClient(httpClient, *u)
	if err != nil {
		panic(err) // TODO handle this properly
	}
	return fleetHTTPAPIClient, err
}

func shutDownNeo() {
	checkDeployerIsStopped()
	// TODO implement this function
	info.Printf("TODO IFWEHAVETO: Use the Go fleet API to shut down neo4j's dependencies.")
	info.Printf("TODO DEFINITELY: Shut down neo4j.")
}

func checkDeployerIsStopped() {
	// TODO implement this function
	info.Printf("TODO DEFINITELY: check using the fleet API that the deployer isn't running.")
	fleetClient, err := newFleetClient()
	if err {
		panic(err) // TODO handle this properly
	}
	unitStates, err := fleetClient.UnitStates()
	info.Printf(unitStates)
}