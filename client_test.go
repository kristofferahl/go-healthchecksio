package healthchecksio

import (
	"os"
	"testing"
)

func assertString(t *testing.T, expected string, actual string) {
	if actual != expected {
		t.Errorf("Expected string of %s but got %s", expected, actual)
	}
}

func configureClient() *Client {
	l := &StandardLogger{}
	envKey := "HEALTHCHECKSIO_API_KEY"
	apiKey := os.Getenv(envKey)
	if apiKey == "" {
		l.Errorf("API Key must be set (env: %s)", envKey)
		os.Exit(1)
	}
	c := NewClient(apiKey)
	c.Log = l
	return c
}

func TestClient(t *testing.T) {
	// Client
	// ----------------------------------------
	client := configureClient()
	if client.APIKey == "" {
		t.Error("API key must be set.")
	}

	// GetAll
	// ----------------------------------------
	checks, err := client.GetAll()
	if err != nil {
		t.Error("Fetching healthchecks failed.", err)
		return
	}

	client.Log.Debugf("Fetched %s", checks)

	// Create
	// ----------------------------------------
	created, err := client.Create(Healthcheck{
		Name: "testcheck",
		Tags: "test devops",
		Unique: []string{
			"name",
		},
	})
	if err != nil {
		t.Error("Create healthcheck failed.", err)
		return
	}

	client.Log.Debugf("[DEBUG] Created %s", created)

	// Update
	// ----------------------------------------
	updated, err := client.Update(created.ID(), Healthcheck{
		Name: "testcheck",
		Tags: "test devops go",
		Unique: []string{
			"name",
		},
	})
	if err != nil {
		t.Error("Updating healthcheck failed.", err)
		return
	}

	client.Log.Debugf("[DEBUG] Updated %s", updated)

	// Pause
	// ----------------------------------------
	paused, err := client.Pause(updated.ID())
	if err != nil {
		t.Error("Pausing healthcheck failed.", err)
		return
	}

	client.Log.Debugf("[DEBUG] Paused %s", paused)

	// Delete
	// ----------------------------------------
	deleted, err := client.Delete(paused.ID())
	if err != nil {
		t.Error("Deleting healthcheck failed.", err)
		return
	}

	client.Log.Debugf("[DEBUG] Deleted %s", deleted)

	// GetAllChannels
	// ----------------------------------------
	channels, err := client.GetAllChannels()
	if err != nil {
		t.Error("Fetching channels failed.", err)
		return
	}

	client.Log.Debugf("[DEBUG] Fetched %s", channels)
}
