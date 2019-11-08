# go-healthchecksio

A [go](https://golang.org) client for working with the [healthchecks.io](https://healthchecks.io) api.

## Installation

    go get github.com/kristofferahl/go-healthchecksio

## Usage

```go
package main

import (
	"log"
	"os"

	"github.com/kristofferahl/go-healthchecksio"
)

func main() {
	apiKey := os.Getenv("HEALTHCHECKSIO_API_KEY") // See https://healthchecks.io/docs/api/

	client := healthchecksio.NewClient(apiKey)
	check := healthchecksio.Healthcheck{
		Name: "My first test",
		Tags: "go ftw",
	}

	healthcheck, err := client.Create(check)
	if err != nil {
		log.Printf("[ERROR] error creating healthcheck: %s", err)
		os.Exit(1)
	}

	log.Printf("[DEBUG] created healthcheck: %s", healthcheck)
}
```

## Documentation

Docs can be found at [godoc.org](https://godoc.org/github.com/kristofferahl/go-healthchecksio).

## A note on logging

By default this package uses the NoOpLogger, essentially turning off all logging. For basic logging you may configure the client to use the included logger StandardLogger. You may also choose to supply your own logger by implementing the Logger interface.

```go
client := healthchecksio.NewClient(apiKey)
client.Log = &healthchecksio.StandardLogger{}
```


## Developing

Running the tests requires a valid healthchecks.io API key (See https://healthchecks.io/docs/api/). Make sure the following environment variables are set.

    export HEALTHCHECKSIO_API_KEY='{your api key}'

## Roadmap

- Required fields validation

## Contributing

...
