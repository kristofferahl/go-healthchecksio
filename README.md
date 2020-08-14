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

## Developing

Running the tests requires a valid healthchecks.io API key (See https://healthchecks.io/docs/api/). Make sure the following environment variables are set.

    export HEALTHCHECKSIO_API_KEY='{your api key}'

## Roadmap

- Required fields validation

## Contributing

...
