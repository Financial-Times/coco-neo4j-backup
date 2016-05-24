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

func newFleetClient(fleetEndpoint string, socksProxy string) (client.API, error) {
	u, err := url.Parse(fleetEndpoint)
	if err != nil {
		panic(err) // TODO handle this properly
	}
	httpClient := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 100}}

	if socksProxy != "" {
		log.Printf("using SOCKS proxy %s\n", socksProxy)
		netDialler := &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		dialer, err := proxy.SOCKS5("tcp", socksProxy, nil, netDialler)
		if err != nil {
			log.Fatalf("error with proxy %s: %v\n", socksProxy, err)
		}
		httpClient.Transport = &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			Dial:                dialer.Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			MaxIdleConnsPerHost: 100,
		}
	}

	info.Printf("Connecting to fleet API on %s", u)
	fleetHTTPAPIClient, err := client.NewHTTPClient(httpClient, *u)
	if err != nil {
		panic(err) // TODO handle this properly
	}
	return fleetHTTPAPIClient, err
}

func shutDownNeo(fleetClient client.API) {
	isDeployerActive, err := isServiceActive(fleetClient, "deployer.service")
	if isDeployerActive || err != nil {
		warn.Printf(`Problem: either the deployer is still active, or there was a problem checking its status.
We cannot complete the backup process in case neo4j is accidentally started up again during backup creation.`)
		panic(err) // TODO handle this properly.
	}
	// TODO implement this function
	info.Printf("TODO IFWEHAVETO: Use the Go fleet API to shut down neo4j's dependencies.")
	info.Printf("TODO DEFINITELY: Shut down neo4j.")
}

func isServiceActive(fleetClient client.API, serviceName string) (bool, error) {
	// TODO implement this function
	info.Printf("TODO DEFINITELY: check using the fleet API that the deployer isn't running.")
	unitStates, err := fleetClient.UnitStates()
	if err != nil {
		warn.Printf("Could not retrieve list of units from fleet API")
		//panic(err) // TODO handle this properly
	}
	info.Printf("%d units retrieved", len(unitStates))
	for index, each := range unitStates {
		if each.Name == serviceName {
			info.Printf("index=%d Name=%s SystemdActiveState=%s SystemdLoadState=%s", index, each.Name, each.SystemdActiveState, each.SystemdLoadState)
			if each.SystemdActiveState == "active" {
				return true, err
			} else {
				return false, err
			}
		}
	}
	warn.Printf("Could not find deployer in list of services!")
	panic(err) // TODO handle this properly by returning a proper error from this function.
}
