package main

import (
	"github.com/coreos/fleet/client"
	"net/url"
	"net"
	"time"
	"golang.org/x/net/proxy"
	"net/http"
	log "github.com/Sirupsen/logrus"
	"github.com/coreos/fleet/functional/util"
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

	log.Infof("Connecting to fleet API on %s", u)
	fleetHTTPAPIClient, err := client.NewHTTPClient(httpClient, *u)
	if err != nil {
		panic(err) // TODO handle this properly
	}
	return fleetHTTPAPIClient, err
}

func shutDownNeo(fleetClient client.API) {
	isDeployerActive, err := isServiceActive(fleetClient, "deployer.service")
	if isDeployerActive || err != nil {
		log.Warnf(`Problem: either the deployer is still active, or there was a problem checking its status.
We cannot complete the backup process in case neo4j is accidentally started up again during backup creation.`)
		panic(err) // TODO handle this properly.
	}
	log.Info("TODO IFWEHAVETO: Use the Go fleet API to shut down neo4j's dependencies.")
	neoServiceName := "neo4j-red@1.service"
	log.Info("Shutting down neo4j-red@1....")
	//fleetClient.SetUnitTargetState(neoServiceName, "inactive")
	util.RunFleetctl("stop neo4j-red@1")
	isNeoActive, err := isServiceActive(fleetClient, neoServiceName)
	if isNeoActive || err != nil {
		log.Warnf(`Problem: either neo4j is still active, or there was a problem checking its status.
We cannot complete the backup process when neo4j is still running.`)
		panic(err) // TODO handle this properly.
	}

}

func isServiceActive(fleetClient client.API, serviceName string) (bool, error) {
	unitStates, err := fleetClient.UnitStates()
	if err != nil {
		log.Panic("Could not retrieve list of units from fleet API, do you need to start a SOCKS proxy?")
		return true, err // TODO handle this properly
	}
	log.WithFields(log.Fields{"num": len(unitStates)}).Info("Retrieved services from fleet API.")
	for index, each := range unitStates {
		if each.Name == serviceName {
			log.WithFields(log.Fields{
				"index": index,
				"name": each.Name,
				"SystemdActiveState": each.SystemdActiveState,
				"SystemdLoadState": each.SystemdLoadState,
			}).Info("Processing service.")
			if each.SystemdActiveState == "active" {
				return true, err
			} else {
				return false, err
			}
		} else {
			log.WithFields(log.Fields{
				"index": index,
				"name": each.Name,
				"SystemdActiveState": each.SystemdActiveState,
				"SystemdLoadState": each.SystemdLoadState,
			}).Info("Processing service.")
		}
	}
	log.WithFields(log.Fields{"serviceName": serviceName}).Panic("Could not find service in list of services!")
	panic(err) // TODO handle this properly by returning a proper error from this function.
}

func startNeo(fleetClient client.API) {
	log.Info("Starting up neo4j...")
	fleetClient.SetUnitTargetState("neo", "active")
	isNeoActive, err := isServiceActive(fleetClient, "deployer.service")
	if !isNeoActive || err != nil {
		log.Warnf(`Problem starting up neo4j: either neo4j is not active, or there was a problem checking its status.
neo4j must be started up to conclude the backup process.`)
		panic(err) // TODO handle this properly.
	}
	log.Info("TODO: Start up neo4j's dependencies.")
}

