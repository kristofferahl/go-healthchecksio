package healthchecksio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const baseURL string = "https://healthchecks.io/api/v1"

// Client provides access to create, read, update and delete healthchecks.io resources
type Client struct {
	APIKey      string
	BaseURL     string
	ContentType string
	HTTPClient  *http.Client
	Log         Logger
}

type apiResponse HealthcheckResponse

type apiListResponse struct {
	Data []*HealthcheckResponse `json:"checks"`
}

type apiListChannelsResponse struct {
	Data []*HealthcheckChannelResponse `json:"channels"`
}

type apiErrorResponse struct {
	Message string `json:"error"`
}

// NewClient creates a new client
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:      apiKey,
		BaseURL:     baseURL,
		ContentType: "application/json",
		HTTPClient:  &http.Client{},
		Log:         &NoOpLogger{},
	}
}

func wrapError(err error, req *http.Request, res *http.Response) error {
	if err != nil {
		status := ""
		statusCode := 0

		if res != nil {
			status = res.Status
			statusCode = res.StatusCode
		}

		return &APIError{
			err:        err.Error(),
			method:     req.Method,
			url:        req.URL.String(),
			status:     status,
			statusCode: statusCode,
		}
	}
	return nil
}

func (c *Client) request(method string, path string, reader io.Reader) ([]byte, error) {
	url := baseURL + path
	req, _ := http.NewRequest(method, url, reader)
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("X-Api-Key", c.APIKey)

	c.Log.Debugf("HTTP %s %s", method, url)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, wrapError(err, req, nil)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode >= 300 {
		errorRes := new(apiErrorResponse)
		if err = json.Unmarshal(body, &errorRes); err != nil {
			return nil, wrapError(err, req, res)
		}

		return nil, wrapError(fmt.Errorf("response error, %s", errorRes.Message), req, res)
	}

	return body, wrapError(err, req, res)
}

func (c *Client) get(path string) ([]byte, error) {
	return c.request("GET", path, nil)
}

func (c *Client) post(path string, body io.Reader) ([]byte, error) {
	return c.request("POST", path, body)
}

func (c *Client) delete(path string) ([]byte, error) {
	return c.request("DELETE", path, nil)
}

// GetAll returns all healthchecks
func (c *Client) GetAll() ([]*HealthcheckResponse, error) {
	body, err := c.get("/checks/")
	if err != nil {
		return nil, err
	}

	r, err := toAPIListResponse(body)
	if err != nil {
		return nil, err
	}

	resp := (*r).Data
	return resp, nil
}

// Create creates a new healthcheck
func (c *Client) Create(check Healthcheck) (*HealthcheckResponse, error) {
	json, err := check.ToJSON()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte(json))
	body, err := c.post("/checks/", buf)
	if err != nil {
		return nil, err
	}

	r, err := toAPIResponse(body)
	if err != nil {
		return nil, err
	}

	resp := HealthcheckResponse(*r)
	return &resp, nil
}

// Update updates an existing healthcheck
func (c *Client) Update(id string, check Healthcheck) (*HealthcheckResponse, error) {
	json, err := check.ToJSON()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte(json))
	body, err := c.post(fmt.Sprintf("/checks/%s", id), buf)
	if err != nil {
		return nil, err
	}

	r, err := toAPIResponse(body)
	if err != nil {
		return nil, err
	}

	resp := HealthcheckResponse(*r)
	return &resp, nil
}

// Pause pauses monitoring on existing healthcheck
func (c *Client) Pause(id string) (*HealthcheckResponse, error) {
	body, err := c.post(fmt.Sprintf("/checks/%s/pause", id), nil)
	if err != nil {
		return nil, err
	}

	r, err := toAPIResponse(body)
	if err != nil {
		return nil, err
	}

	resp := HealthcheckResponse(*r)
	return &resp, nil
}

// Delete deletes an existing healthcheck
func (c *Client) Delete(id string) (*HealthcheckResponse, error) {
	body, err := c.delete(fmt.Sprintf("/checks/%s", id))
	if err != nil {
		return nil, err
	}

	r, err := toAPIResponse(body)
	if err != nil {
		return nil, err
	}

	resp := HealthcheckResponse(*r)
	return &resp, nil
}

// GetAllChannels returns all channels
func (c *Client) GetAllChannels() ([]*HealthcheckChannelResponse, error) {
	body, err := c.get("/channels/")
	if err != nil {
		return nil, err
	}

	r, err := toAPIListChannelsResponse(body)
	if err != nil {
		return nil, err
	}

	resp := (*r).Data
	return resp, nil
}

func toAPIListResponse(body []byte) (*apiListResponse, error) {
	var s = new(apiListResponse)
	err := json.Unmarshal(body, &s)
	return s, err
}

func toAPIResponse(body []byte) (*apiResponse, error) {
	var s = new(apiResponse)
	err := json.Unmarshal(body, &s)
	return s, err
}

func toAPIListChannelsResponse(body []byte) (*apiListChannelsResponse, error) {
	var s = new(apiListChannelsResponse)
	err := json.Unmarshal(body, &s)
	return s, err
}
