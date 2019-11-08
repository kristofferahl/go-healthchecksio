package healthchecksio

import (
	"encoding/json"
	"strings"
)

// APIError represents an error that occured when talkinging to the Healthchecks.io api
type APIError struct {
	err        string
	method     string
	url        string
	status     string
	statusCode int
}

func (e *APIError) Error() string {
	return e.err
}

// Method returns the HTTP request method
func (e *APIError) Method() string {
	return e.method
}

// URL returns the HTTP request URL
func (e *APIError) URL() string {
	return e.url
}

// Status returns the HTTP response status
func (e *APIError) Status() string {
	return e.status
}

// StatusCode returns the HTTP response status code
func (e *APIError) StatusCode() int {
	return e.statusCode
}

// Healthcheck represents a healthcheck
type Healthcheck struct {
	Channels string   `json:"channels,omitempty"`
	Grace    int      `json:"grace,omitempty"`
	Name     string   `json:"name,omitempty"`
	Schedule string   `json:"schedule,omitempty"`
	Tags     string   `json:"tags,omitempty"`
	Timeout  int      `json:"timeout,omitempty"`
	Timezone string   `json:"tz,omitempty"`
	Unique   []string `json:"unique,omitempty"`
}

// HealthcheckResponse represents a healthcheck api response
type HealthcheckResponse struct {
	Channels  string `json:"channels,omitempty"`
	Grace     int    `json:"grace,omitempty"`
	LastPing  string `json:"last_ping,omitempty"`
	Name      string `json:"name,omitempty"`
	NextPing  string `json:"next_ping,omitempty"`
	PauseURL  string `json:"pause_url,omitempty"`
	Pings     int    `json:"n_pings,omitempty"`
	PingURL   string `json:"ping_url,omitempty"`
	Schedule  string `json:"schedule,omitempty"`
	Status    string `json:"status,omitempty"`
	Tags      string `json:"tags,omitempty"`
	Timeout   int    `json:"timeout,omitempty"`
	Timezone  string `json:"tz,omitempty"`
	UpdateURL string `json:"update_url,omitempty"`
}

// HealthcheckChannelResponse represents a channel response of healthcheck api
type HealthcheckChannelResponse struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Kind string `json:"kind,omitempty"`
}

// ToJSON returns a json representation of a healthcheck data
func (hc *Healthcheck) ToJSON() (string, error) {
	b, err := json.Marshal(hc)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ID returns the identifier of a healthcheck
func (hc *HealthcheckResponse) ID() string {
	a := strings.Split(hc.UpdateURL, "/")
	return a[len(a)-1]
}

// ToJSON returns a json representation of a healthcheck
func (hc *HealthcheckResponse) ToJSON() (string, error) {
	b, err := json.MarshalIndent(hc, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (hc *HealthcheckResponse) String() string {
	json, err := hc.ToJSON()
	if err != nil {
		return err.Error()
	}
	return json
}

// ToJSON returns a json representation of a healthcheck channel
func (hc *HealthcheckChannelResponse) ToJSON() (string, error) {
	b, err := json.MarshalIndent(hc, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (hc *HealthcheckChannelResponse) String() string {
	json, err := hc.ToJSON()
	if err != nil {
		return err.Error()
	}
	return json
}
