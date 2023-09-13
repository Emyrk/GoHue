package gohue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Emyrk/gohue/hueclient"
	"io"
	"net/http"
)

type Client struct {
	Bridges  []Bridge
	Client   *http.Client
	username string
}

func NewClient(username string) (*Client, error) {
	return &Client{
		//Bridges:  bridges,
		Client:   hueclient.DefaultClient(),
		username: username,
	}, nil
}

type APIKeyResponse struct {
	Success struct {
		Username  string `json:"username"`
		ClientKey string `json:"clientkey"`
	} `json:"success"`
}

func (c *Client) GenerateAPIKey(ctx context.Context) ([]APIKeyResponse, error) {
	var resp []APIKeyResponse
	err := c.request(ctx, "POST", "/api", map[string]any{"devicetype": "gohue#myapp", "generateclientkey": true}, &resp)
	return resp, err
}

type Devices struct {
	Errors []interface{} `json:"errors"`
	Data   []Device      `json:"data"`
}
type ProductData struct {
	ModelID              string `json:"model_id"`
	ManufacturerName     string `json:"manufacturer_name"`
	ProductName          string `json:"product_name"`
	ProductArchetype     string `json:"product_archetype"`
	Certified            bool   `json:"certified"`
	SoftwareVersion      string `json:"software_version"`
	HardwarePlatformType string `json:"hardware_platform_type"`
}
type Metadata struct {
	Name      string `json:"name"`
	Archetype string `json:"archetype"`
}
type Identify struct {
}
type Services struct {
	Rid   string `json:"rid"`
	Rtype string `json:"rtype"`
}

type Usertest struct {
	Status   string `json:"status"`
	Usertest bool   `json:"usertest"`
}

type Device struct {
	ID          string      `json:"id"`
	ProductData ProductData `json:"product_data,omitempty"`
	Metadata    Metadata    `json:"metadata"`
	Identify    Identify    `json:"identify,omitempty"`
	Services    []Services  `json:"services"`
	Type        string      `json:"type"`
	IDV1        string      `json:"id_v1,omitempty"`
	Usertest    Usertest    `json:"usertest,omitempty"`
}

func (c *Client) Devices(ctx context.Context) (any, error) {
	err := c.request(ctx, "GET", "/clip/v2/resource/device", nil, nil)
	if err != nil {
		return nil, err
	}
	return nil, err
}

func (c *Client) request(ctx context.Context, method string, route string, body any, respStruct any) error {
	var reqBody io.Reader = nil
	if body != nil {
		var out bytes.Buffer
		err := json.NewEncoder(&out).Encode(body)
		if err != nil {
			return err
		}
		reqBody = &out
	}
	req, err := http.NewRequest(method, c.requestRoute(route), reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("hue-application-key", c.username)

	req = req.WithContext(ctx)
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	if err := debugHttpResponse(ctx, resp); err != nil {
		return fmt.Errorf("debug request: %w", err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(respStruct)
	return err
}

func (c *Client) requestRoute(route string) string {
	return "https://" + c.Bridges[0].Addr.String() + route
}
