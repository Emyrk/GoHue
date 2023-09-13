package gohue

import (
	"context"
	"encoding/json"
	"github.com/Emyrk/gohue/hueclient"
	"net/http"
)

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
	req, err := http.NewRequest("GET", "https://"+b.Addr.String()+"/api/0/config", nil)
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
