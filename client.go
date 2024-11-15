// Package larditrans provides a client for interacting with the Lardi-Trans API.
package lardiAPI

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// API endpoints
const (
	defaultBaseURL = "https://api.lardi-trans.com"
	defaultTimeout = 30 * time.Second
)

// Endpoint paths
const (
	pathCargo        = "/v2/proposals/my/add/cargo"
	pathCurrencies   = "/v2/references/currencies"
	pathUnits        = "/v2/references/payment/units"
	pathMoments      = "/v2/references/payment/moments"
	pathTypes        = "/v2/references/body/types"
	pathPackage      = "/v2/references/cargo/package"
	pathTypesPayment = "/v2/references/payment/types"
	pathLoadTypes    = "/v2/references/load/types"
	pathAreas        = "/v2/references/areas"
	pathContacts     = "/v2/users/user/contacts"
)

// Config contains the configuration for the API client
type Config struct {
	BaseURL  string
	APIKey   string
	Timeout  time.Duration
	Language string
}

// Client represents a client for the Lardi-Trans API
type Client struct {
	config Config
	http   HTTPClient
}

// HTTPClient interface allows for easy mocking in tests
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// NewClient creates a new Lardi-Trans API client with the given configuration
func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}
	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}
	if config.Language == "" {
		config.Language = "uk"
	}

	return &Client{
		config: config,
		http: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

type ResponseContacts struct {
	ContactID   int    `json:"contactId"`
	ContactName string `json:"face"`
}

// Response represents a generic API response
type Response struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CargoRequest represents the request body for creating a cargo proposal
type CargoRequest struct {
	ContactID          int           `json:"contactId,omitempty"`
	DateFrom           string        `json:"dateFrom,omitempty"`
	DateTo             string        `json:"dateTo,omitempty"`
	PaymentValue       int           `json:"paymentPrice,omitempty"`
	PaymentCurrencyID  int           `json:"paymentCurrencyId,omitempty"`
	PaymentUnitID      int           `json:"paymentUnitId,omitempty"`
	PaymentMomentID    int           `json:"paymentMomentId,omitempty"`
	CargoBodyTypeIDs   []int         `json:"cargoBodyTypeIds,omitempty"`
	CargoPackaging     []CargoPack   `json:"cargoPackaging,omitempty"`
	PaymentForms       []PaymentForm `json:"paymentForms,omitempty"`
	LorryAmount        int           `json:"lorryAmount,omitempty"`
	LoadTypes          []int         `json:"loadTypes,omitempty"`
	Groupage           bool          `json:"groupage,omitempty"`
	ContentName        string        `json:"contentName,omitempty"`
	SizeMass           float64       `json:"sizeMass,omitempty"`
	SizeVolume         float64       `json:"sizeVolume,omitempty"`
	WaypointListSource []LoadParams  `json:"waypointListSource"`
	WaypointListTarget []LoadParams  `json:"waypointListTarget"`
}

type PaymentForm struct {
	ID  int  `json:"id,omitempty"`
	Vat bool `json:"vat,omitempty"`
}

type CargoPack struct {
	ID    int `json:"id,omitempty"`
	Count int `json:"count,omitempty"`
}

// Validate checks if the required fields are set
func (r *CargoRequest) Validate() error {
	if len(r.WaypointListSource) == 0 {
		return fmt.Errorf("waypointListSource is required")
	}
	if len(r.WaypointListTarget) == 0 {
		return fmt.Errorf("waypointListTarget is required")
	}
	return nil
}

// LoadParams represents loading/unloading point parameters
type LoadParams struct {
	TownName    string   `json:"townName"`
	AreaID      int      `json:"areaId"`
	CountrySign string   `json:"countrySign"`
	RegionID    int      `json:"regionId"`
	PostCodes   []string `json:"postCode"`
}

type Request struct {
	Name string
}

// CargoResponse represents the response from creating a cargo proposal
type CargoResponse struct {
	ID int `json:"id"`
}

// APIError represents an error response from the API
type APIError struct {
	Status  int    `json:"status"`
	Err     string `json:"error"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: status=%d, error=%s, message=%s", e.Status, e.Err, e.Message)
}

// CreateCargo creates a new cargo proposal
func (c *Client) CreateCargo(ctx context.Context, req *CargoRequest) (*CargoResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	var resp CargoResponse
	err := c.post(ctx, pathCargo, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("create cargo request failed: %w", err)
	}

	return &resp, nil
}

// GetContacts retrieves available contacts
func (c *Client) GetContacts(ctx context.Context) ([]ResponseContacts, error) {
	var resp []ResponseContacts
	err := c.get(ctx, pathContacts, &resp)
	if err != nil {
		return nil, fmt.Errorf("get contacts failed: %w", err)
	}
	return resp, nil
}

// GetAreas retrieves available areas
func (c *Client) GetAreas(ctx context.Context, area Request) (*Response, error) {
	var resp []Response
	err := c.get(ctx, pathAreas, &resp)
	if err != nil {
		return nil, fmt.Errorf("get areas failed: %w", err)
	}
	for _, v := range resp {
		if v.Name == area.Name {
			return &Response{
				ID:   v.ID,
				Name: v.Name,
			}, nil
		}
	}
	return nil, nil
}

// GetLoadTypes retrieves available load types
func (c *Client) GetLoadTypes(ctx context.Context) ([]Response, error) {
	var resp []Response
	err := c.get(ctx, pathLoadTypes, &resp)
	if err != nil {
		return nil, fmt.Errorf("get load types failed: %w", err)
	}
	return resp, nil
}

// GetPaymentTypes retrieves available payment types
func (c *Client) GetPaymentTypes(ctx context.Context) ([]Response, error) {
	var resp []Response
	err := c.get(ctx, pathTypesPayment, &resp)
	if err != nil {
		return nil, fmt.Errorf("get payment types failed: %w", err)
	}
	return resp, nil
}

// GetPackageTypes retrieves available package types
func (c *Client) GetPackageTypes(ctx context.Context) ([]Response, error) {
	var resp []Response
	err := c.get(ctx, pathPackage, &resp)
	if err != nil {
		return nil, fmt.Errorf("get package types failed: %w", err)
	}
	return resp, nil
}

// GetBodyTypes retrieves available body types
func (c *Client) GetBodyTypes(ctx context.Context, body Request) (*Response, error) {
	var resp []Response
	err := c.get(ctx, pathTypes, &resp)
	if err != nil {
		return nil, fmt.Errorf("get body types failed: %w", err)
	}
	for _, v := range resp {
		if body.Name == v.Name {
			return &Response{
				ID:   v.ID,
				Name: v.Name,
			}, nil
		}
	}

	return nil, nil
}

// GetPaymentMoments retrieves available payment moments
func (c *Client) GetPaymentMoments(ctx context.Context) ([]Response, error) {
	var resp []Response
	err := c.get(ctx, pathMoments, &resp)
	if err != nil {
		return nil, fmt.Errorf("get payment moments failed: %w", err)
	}
	return resp, nil
}

// GetCurrencies retrieves available currencies
func (c *Client) GetCurrencies(ctx context.Context, currency Request) (*Response, error) {
	var resp []Response
	err := c.get(ctx, pathCurrencies, &resp)
	if err != nil {
		return nil, fmt.Errorf("get currencies failed: %w", err)
	}
	for _, v := range resp {
		if v.Name == currency.Name {
			return &Response{
				ID:   v.ID,
				Name: v.Name,
			}, nil
		}
	}

	return nil, nil
}

// GetUnits retrieves available units
func (c *Client) GetUnits(ctx context.Context) ([]Response, error) {
	var resp []Response
	err := c.get(ctx, pathUnits, &resp)
	if err != nil {
		return nil, fmt.Errorf("get units failed: %w", err)
	}
	return resp, nil
}

// post performs a POST request
func (c *Client) post(ctx context.Context, path string, body interface{}, result interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.BaseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	return c.doRequest(req, result)
}

// get performs a GET request
func (c *Client) get(ctx context.Context, path string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.config.BaseURL+path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	return c.doRequest(req, result)
}

// doRequest performs the HTTP request and handles the response
func (c *Client) doRequest(req *http.Request, result interface{}) error {
	req.Header.Set("Authorization", c.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("language", c.config.Language)
	req.URL.RawQuery = q.Encode()

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		return &apiErr
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
