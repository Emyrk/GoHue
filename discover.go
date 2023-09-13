package gohue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/mdns"
	"net"
	"net/http"
	"strings"
)

type BridgeDiscoverMethod int

const (
	DiscoveryManual BridgeDiscoverMethod = iota
	DiscoveryMDNS
	DiscoveryEndpoint
)

func DiscoverBridges(ctx context.Context) ([]Bridge, error) {
	bridges, err := MDNSDiscover()
	if err == nil && len(bridges) > 0 {
		return bridges, nil
	}

	bridges, err = CloudDiscovery(ctx)
	if err == nil && len(bridges) > 0 {
		return bridges, nil
	}

	return nil, errors.New("no bridges found")
}

type Bridge struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	ModelID string `json:"modelid"`

	AddrIPv4 net.IP `json:"addr_ipv4"`
	AddrIPv6 net.IP `json:"addr_ipv6"`
	// Addr is going to be either ipv4 or ipv6
	Addr net.IP `json:"addr"`
	Port int    `json:"port"`

	DiscoveredBy BridgeDiscoverMethod `json:"discovered_by"`
}

// MDNSDiscover will discover bridges using mDNS
func MDNSDiscover() ([]Bridge, error) {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	var dnsErr error
	go func() {
		dnsErr = mdns.Lookup("_hue._tcp", entriesCh)
		close(entriesCh)
	}()

	var bridges []Bridge
	for entry := range entriesCh {
		bridge := Bridge{
			ID:           "",
			Name:         entry.Name,
			ModelID:      "",
			AddrIPv4:     entry.AddrV4,
			AddrIPv6:     entry.AddrV6,
			Addr:         entry.Addr,
			Port:         entry.Port,
			DiscoveredBy: DiscoveryMDNS,
		}
		for _, info := range entry.InfoFields {
			parts := strings.Split(info, "=")
			if len(parts) != 2 {
				continue
			}
			switch parts[0] {
			case "bridgeid":
				bridge.ID = parts[1]
			case "modelid":
				bridge.ModelID = parts[1]
			}
		}
		bridges = append(bridges, bridge)
	}

	return bridges, dnsErr
}

// CloudDiscovery asks hue cloud for the list of bridges.
// This is rate limited.
func CloudDiscovery(ctx context.Context) ([]Bridge, error) {
	cli := http.DefaultClient
	req, err := http.NewRequest("GET", "https://discovery.meethue.com", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	if err := debugHttpResponse(ctx, resp); err != nil {
		return nil, fmt.Errorf("debug request: %w", err)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, errors.New("rate limited from discovery endpoint, try again later in 15minutes")
	}

	type discoverBridge struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		//  "id":"001788fffe100491",
		//  "internalipaddress":"192.168.2.23",
		//  "macaddress":"00:17:88:10:04:91",
		//  "name":"Philips Hue"
		InternapIPAddress string `json:"internalipaddress"`
		MACAddress        string `json:"macaddress"`
	}
	var bridges []discoverBridge
	err = json.NewDecoder(resp.Body).Decode(&bridges)
	if err != nil {
		return nil, err
	}

	var normalizedBridges []Bridge
	for _, bridge := range bridges {
		addr := net.ParseIP(bridge.InternapIPAddress)
		if addr == nil {
			return nil, fmt.Errorf("invalid ip address %q", bridge.InternapIPAddress)
		}

		nb := Bridge{
			ID:           bridge.ID,
			Name:         bridge.Name,
			Addr:         addr,
			Port:         0,
			DiscoveredBy: DiscoveryEndpoint,
		}
		if len(nb.Addr) == net.IPv4len {
			nb.AddrIPv4 = nb.Addr
		} else {
			nb.AddrIPv6 = nb.Addr
		}

		normalizedBridges = append(normalizedBridges, nb)
	}

	return normalizedBridges, nil
}
