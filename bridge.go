package gohue

import (
	"context"
	"encoding/json"
	"github.com/Emyrk/gohue/hueclient"
	"net/http"
)

type Bridge struct {
	ID                string `json:"id"`
	InternalIPAddress string `json:"internalipaddress"`
	MacAddress        string `json:"macaddress"`
	Name              string `json:"name"`
}

type BridgeConfig struct {
	Name       string `json:"name"`
	SWVersion  string `json:"swversion"`
	APIVersion string `json:"apiversion"`
	MAC        string `json:"mac"`
	BridgeID   string `json:"bridgeid"`
	FactoryNew bool   `json:"factorynew"`
	//ReplacesBridgeID string `json:"replacesbridgeid"`
	ModelID string `json:"modelid"`
}

func (b Bridge) Config(ctx context.Context) (BridgeConfig, error) {
	cli := hueclient.DefaultClient()
	req, err := http.NewRequest("GET", "https://"+b.InternalIPAddress+"/api/0/config", nil)
	if err != nil {
		return BridgeConfig{}, err
	}

	req = req.WithContext(ctx)
	resp, err := cli.Do(req)
	defer resp.Body.Close()

	var cfg BridgeConfig
	err = json.NewDecoder(resp.Body).Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Implement this: https://developers.meethue.com/develop/application-design-guidance/hue-bridge-discovery/
func Discover(ctx context.Context) ([]Bridge, error) {
	cli := hueclient.DefaultClient()
	req, err := http.NewRequest("GET", "https://discovery.meethue.com", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var bridges []Bridge
	err = json.NewDecoder(resp.Body).Decode(&bridges)
	if err != nil {
		return nil, err
	}

	return bridges, nil
}
