package healthchecksio

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	checkNamePrefix   = "go-healthchecksio-check-"
	channelNamePrefix = "go-healthchecksio-channel-"
)

var (
	defaultChannel            = channelNamePrefix + "default"
	seededRand     *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func newCheck() Healthcheck {
	return Healthcheck{
		Name:        randomName(),
		Description: "Basic check",
		Tags:        "test healthchecksio client",
		Timeout:     99,
		Grace:       100,
		Unique: []string{
			"name",
		},
	}
}

func randomName() string {
	return checkNamePrefix + "rand-" + strconv.Itoa(seededRand.Int())
}

func configureClient() *Client {
	envKey := "HEALTHCHECKSIO_API_KEY"
	apiKey := os.Getenv(envKey)
	if apiKey == "" {
		log.Printf("API Key must be set (env: %s)\n", envKey)
		os.Exit(1)
	}
	return NewClient(apiKey)
}

func cleanupTestData(tb testing.TB, client *Client) {
	checks, err := client.GetAll()
	if err != nil {
		tb.Error("Test cleanup failed, fetching all healthchecks.", err)
		return
	}

	for _, check := range checks {
		if !strings.HasPrefix(check.Name, checkNamePrefix) {
			continue
		}

		_, err := client.Delete(check.ID())
		if err != nil {
			tb.Error("Test cleanup failed, deleting healthcheck.", err)
			return
		}
	}
}

func getChannel(client *Client, name string) (*HealthcheckChannelResponse, error) {
	channels, err := client.GetAllChannels()
	if err != nil {
		return nil, err
	}
	for _, channel := range channels {
		if channel.Name == name {
			return channel, nil
		}
	}
	return nil, fmt.Errorf("channel not found: %s", name)
}

// You can use testing.T, if you want to test the code without benchmarking
func setupSuite(tb testing.TB) (teardown func(tb testing.TB), client *Client) {
	log.Println("setup suite")

	client = configureClient()
	if client.APIKey == "" {
		tb.Error("API key must be set.")
	}
	cleanupTestData(tb, client)

	// Return a function to teardown the test
	return func(tb testing.TB) {
		log.Println("teardown suite")
	}, client
}

// Almost the same as the above, but this one is for single test instead of collection of tests
func setupTest(tb testing.TB, client *Client) func(tb testing.TB) {
	log.Println("setup test")

	return func(tb testing.TB) {
		log.Println("teardown test")
		cleanupTestData(tb, client)
	}
}

func TestClient(t *testing.T) {
	teardownSuite, client := setupSuite(t)
	defer teardownSuite(t)

	createCheck := func(t *testing.T) *HealthcheckResponse {
		check := newCheck()
		c, err := client.Create(check)
		if err != nil {
			t.Error("Test setup failed.", err)
			return nil
		}
		return c
	}

	t.Run("client", func(t *testing.T) {
		teardownTest := setupTest(t, client)
		defer teardownTest(t)

		t.Run("create", func(t *testing.T) {
			t.Run("creates check", func(t *testing.T) {
				check := newCheck()
				created, err := client.Create(check)

				assert.NoError(t, err)
				assert.NotNil(t, created, "expected check, %v", created)
				assert.Equal(t, check.Name, created.Name)
				assert.Equal(t, check.Description, created.Description)
				assert.Equal(t, check.Tags, created.Tags)
				assert.Equal(t, check.Grace, created.Grace)
				assert.Equal(t, check.Timeout, created.Timeout)
				assert.Equal(t, "new", created.Status)
			})
		})

		t.Run("update", func(t *testing.T) {
			created := createCheck(t)

			t.Run("updates check", func(t *testing.T) {
				check := newCheck()
				check.Description = "Basic check updated"
				check.Tags = "test devops updated"
				check.Grace = 100
				check.Timeout = 100

				updated, err := client.Update(created.ID(), check)

				assert.NoError(t, err)
				assert.NotNil(t, updated, "expected check, %v", updated)
				assert.Equal(t, check.Name, updated.Name)
				assert.Equal(t, check.Description, updated.Description)
				assert.Equal(t, check.Tags, updated.Tags)
				assert.Equal(t, check.Grace, updated.Grace)
				assert.Equal(t, check.Timeout, updated.Timeout)
				assert.Equal(t, "new", updated.Status)
			})

			t.Run("updates check with channel", func(t *testing.T) {
				dc, err := getChannel(client, defaultChannel)
				assert.NoError(t, err)

				// Add channel
				update1, err := client.Update(created.ID(), Healthcheck{
					Channels: dc.Name,
				})
				assert.NoError(t, err)
				assert.Equal(t, dc.ID, update1.Channels)

				// Remove channel
				update2, err := client.Update(created.ID(), Healthcheck{
					Channels: "",
				})
				assert.NoError(t, err)
				assert.Equal(t, "", update2.Channels)
			})

			t.Run("updates check with methods", func(t *testing.T) {
				// Set methods
				update1, err := client.Update(created.ID(), Healthcheck{
					Methods: "POST",
				})
				assert.NoError(t, err)
				assert.Equal(t, "POST", update1.Methods)

				// Remove methods
				update2, err := client.Update(created.ID(), Healthcheck{
					Methods: "",
				})
				assert.NoError(t, err)
				assert.Equal(t, "", update2.Methods)
			})
		})

		t.Run("pause", func(t *testing.T) {
			created := createCheck(t)

			t.Run("pauses check", func(t *testing.T) {
				paused, err := client.Pause(created.ID())
				assert.NoError(t, err)
				assert.Equal(t, "paused", paused.Status)
			})
		})
	})

	t.Run("delete", func(t *testing.T) {
		created := createCheck(t)

		t.Run("deletes check", func(t *testing.T) {
			_, err := client.Delete(created.ID())
			assert.NoError(t, err)

			all, err := client.GetAll()
			assert.NoError(t, err)
			for _, c := range all {
				assert.NotEqual(t, created.ID(), c.ID())
			}
		})
	})

	t.Run("get all", func(t *testing.T) {
		created := createCheck(t)

		t.Run("gets checks", func(t *testing.T) {
			checks, err := client.GetAll()
			assert.NoError(t, err)
			assert.NotNil(t, checks, "expected checks, %v", checks)
			assert.GreaterOrEqual(t, len(checks), 1, "expected at least 1 check, %v", checks)

			for _, check := range checks {
				if check.ID() == created.ID() {
					assert.Equal(t, created.Name, check.Name)
					return
				}
			}

			t.Error(fmt.Errorf("expected check not found, %s", created.ID()))
		})
	})

	t.Run("get channels", func(t *testing.T) {
		t.Run("gets channels", func(t *testing.T) {
			channels, err := client.GetAllChannels()
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(channels), 1, "expected at least 1 channel, %v", channels)

			for _, channel := range channels {
				if channel.Name == defaultChannel {
					return
				}
			}

			t.Error(fmt.Errorf("expected default channel not found, %s", defaultChannel))
		})
	})
}
