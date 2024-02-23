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

	. "github.com/franela/goblin"
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

func cleanupTestData(t *testing.T, client *Client) {
	checks, err := client.GetAll()
	if err != nil {
		t.Error("Test cleanup failed, fetching all healthchecks.", err)
		return
	}

	for _, check := range checks {
		if !strings.HasPrefix(check.Name, checkNamePrefix) {
			continue
		}

		_, err := client.Delete(check.ID())
		if err != nil {
			t.Error("Test cleanup failed, deleting healthcheck.", err)
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

func TestClient(t *testing.T) {
	g := Goblin(t)

	var client *Client

	g.Describe("client", func() {
		g.Before(func() {
			client = configureClient()
			if client.APIKey == "" {
				t.Error("API key must be set.")
			}
			cleanupTestData(t, client)
		})

		g.After(func() {
			cleanupTestData(t, client)
		})

		g.Describe("create", func() {
			g.It("creates check", func() {
				check := newCheck()
				created, err := client.Create(check)

				g.Assert(err).Equal(nil)
				g.Assert(created != nil).IsTrue("expected check, %v", created)
				g.Assert(created.Name).Equal(check.Name)
				g.Assert(created.Description).Equal(check.Description)
				g.Assert(created.Tags).Equal(check.Tags)
				g.Assert(created.Grace).Equal(check.Grace)
				g.Assert(created.Timeout).Equal(check.Timeout)
				g.Assert(created.Status).Equal("new")
			})
		})

		g.Describe("update", func() {
			var created *HealthcheckResponse
			g.BeforeEach(func() {
				check := newCheck()
				c, err := client.Create(check)
				if err != nil {
					t.Error("Test setup failed.", err)
					return
				}
				created = c
			})

			g.It("updates check", func() {
				check := newCheck()
				check.Description = "Basic check updated"
				check.Tags = "test devops updated"
				check.Grace = 100
				check.Timeout = 100

				updated, err := client.Update(created.ID(), check)

				g.Assert(err).Equal(nil)
				g.Assert(updated != nil).IsTrue("expected check, %v", updated)
				g.Assert(updated.Name).Equal(check.Name)
				g.Assert(updated.Description).Equal(check.Description)
				g.Assert(updated.Tags).Equal(check.Tags)
				g.Assert(updated.Grace).Equal(check.Grace)
				g.Assert(updated.Timeout).Equal(check.Timeout)
				g.Assert(updated.Status).Equal("new")
			})
		})

		g.Describe("pause", func() {
			var created *HealthcheckResponse
			g.BeforeEach(func() {
				check := newCheck()
				c, err := client.Create(check)
				if err != nil {
					t.Error("Test setup failed.", err)
					return
				}
				created = c
			})

			g.It("pauses check", func() {
				paused, err := client.Pause(created.ID())
				g.Assert(err).Equal(nil)
				g.Assert(paused.Status).Equal("paused")
			})
		})

		g.Describe("delete", func() {
			var created *HealthcheckResponse
			g.BeforeEach(func() {
				check := newCheck()
				c, err := client.Create(check)
				if err != nil {
					t.Error("Test setup failed.", err)
					return
				}
				created = c
			})

			g.It("deletes check", func() {
				_, err := client.Delete(created.ID())
				g.Assert(err).Equal(nil)
			})
		})

		g.Describe("get all", func() {
			var created *HealthcheckResponse
			g.BeforeEach(func() {
				check := newCheck()
				c, err := client.Create(check)
				if err != nil {
					t.Error("Test setup failed.", err)
					return
				}
				created = c
			})

			g.It("gets checks", func() {
				checks, err := client.GetAll()
				g.Assert(checks != nil).IsTrue("expected checks, %v", checks)
				g.Assert(err).Equal(nil)
				g.Assert(len(checks) > 0).IsTrue("expected at least 1 check, %v", checks)

				for _, check := range checks {
					if check.ID() == created.ID() {
						g.Assert(check.Name).Equal(created.Name)
						return
					}
				}

				g.Fail(fmt.Errorf("expected check not found, %s", created.ID()))
			})
		})

		g.Describe("get channels", func() {
			g.It("gets channels", func() {
				channels, err := client.GetAllChannels()
				g.Assert(err).Equal(nil)
				g.Assert(len(channels) > 0).IsTrue("expected at least 1 channel, %v", channels)

				for _, channel := range channels {
					if channel.Name == defaultChannel {
						return
					}
				}

				g.Fail(fmt.Errorf("expected default channel not found, %s", defaultChannel))
			})
		})
	})
}
